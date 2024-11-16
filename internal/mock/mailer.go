package mock

import (
	"bookshop/internal/mailer"
	"errors"
)

func NewMailTrap() mailer.Mailer {
	return &MailTrap{}
}

type MailTrap struct {
}

func (m MailTrap) Send(recipient, templateFile string, data any) error {
	switch recipient {
	case PanicEmail:
		panic("recover me")
	case FailedToSendEmail:
		return errors.New("empty")
	}
	return nil
}
