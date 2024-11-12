package mailer

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

//go:embed "templates"
var templateFS embed.FS

type Mailer interface {
	Send(recipient, templateFile string, data any) error
}

type MailTrap struct {
	url    string
	token  string
	sender string
}

func NewMailTrap(url, token, sender string) Mailer {
	return &MailTrap{
		url:    url,
		token:  token,
		sender: sender,
	}
}

type EmailPayload struct {
	From    email   `json:"from"`
	To      []email `json:"to"`
	Subject string  `json:"subject"`
	Text    string  `json:"text"`
	Html    string  `json:"html"`
}

type email struct {
	Email string `json:"email"`
}

func (m MailTrap) Send(recipient, templateFile string, data any) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	emailPayload := EmailPayload{
		From:    email{m.sender},
		To:      []email{{recipient}},
		Subject: subject.String(),
		Text:    plainBody.String(),
		Html:    htmlBody.String(),
	}

	payload, err := json.Marshal(emailPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, m.url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+m.token)
	req.Header.Add("Content-Type", "application/json")

	var res *http.Response
	client := &http.Client{}
	for range 5 {
		res, err = client.Do(req)
		if nil == err {
			res.Body.Close()
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return err
}
