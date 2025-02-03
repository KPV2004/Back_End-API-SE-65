package adapters

import (
	"go-server/core"

	"gopkg.in/gomail.v2"
)

type EmailRepository struct {
	Mailer *gomail.Dialer
}

func NewEmailRepository(Mailer *gomail.Dialer) core.EmailRepository {
	return &EmailRepository{Mailer: Mailer}
}

func (m *EmailRepository) Send(message *gomail.Message) error {
	return m.Mailer.DialAndSend(message)
}
