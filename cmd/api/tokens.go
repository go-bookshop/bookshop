package main

import (
	"bookshop/internal/httputil"
	"bookshop/internal/mailer"
	"bookshop/internal/models"
	"bookshop/internal/repository"
	"bookshop/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) resendActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}

	err := httputil.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if models.ValidateEmail(v, input.Email); !v.Valid() {
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	user, err := app.repositories.UserRepository.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.validationErrorResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Activated {
		v.AddError("email", "user has already been activated")
		app.validationErrorResponse(w, r, v.Errors)
		return
	}

	token, err := app.repositories.TokenRepository.New(user.ID, 24*time.Hour, models.ScopeActivation)
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

		err := app.mailer.Send(user.Email, mailer.ResendActivationTemplateFile, activationData)
		if err != nil {
			app.logger.Error("failed to send email", "err", err.Error())
			return
		}
	}()
	w.WriteHeader(http.StatusAccepted)
}
