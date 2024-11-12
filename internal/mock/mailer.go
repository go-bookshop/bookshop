package mock

import (
	"bookshop/internal/mailer"
	"strings"
)

func NewMailTrap() mailer.Mailer {
	return &MailTrap{}
}

type MailTrap struct {
}

func (m MailTrap) Send(recipient, templateFile string, data any) error {
	if localPart, _, found := strings.Cut(recipient, "@"); found {
		switch strings.ToLower(localPart) {
		case "panic":
			panic("recover me")
		}
	}
	return nil
}
