package main

import (
	"bookshop/internal/data"
	"bookshop/internal/httputil"
	"bookshop/internal/mailer"
	"bookshop/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	v := validator.New()
	if data.ValidatePassword(v, input.Password); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if data.ValidateUser(v, user); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	err = app.repositories.UserRepository.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateItem):
			v.AddError("email", "user already exists")
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.repositories.TokenRepository.New(user.ID, 24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	go func() {
		defer func() {
			if recoveredErr := recover(); recoveredErr != nil {
				app.logger.Error("%v", recoveredErr)
			}
		}()

		activationData := struct {
			FirstName string
			LastName  string
			URL       string
		}{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			URL:       fmt.Sprintf("www.shouldbesomefepageforregistration.com/activate?token=%s", token.Plaintext),
		}

		err := app.mailer.Send(user.Email, mailer.UserActivationTemplateFile, activationData)
		if err != nil {
			app.logger.Error("failed to send email", "err", err.Error())
			return
		}
	}()

	err = httputil.WriteJSON(w, http.StatusAccepted, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PlaintextToken string `json:"token"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateTokenPlaintext(v, input.PlaintextToken); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	user, err := app.repositories.UserRepository.GetByToken(data.ScopeActivation, input.PlaintextToken)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.repositories.UserRepository.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.repositories.TokenRepository.DeleteAllByUserID(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = httputil.WriteJSON(w, http.StatusOK, user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
