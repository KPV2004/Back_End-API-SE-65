package core

import "gopkg.in/gomail.v2"

type EmailRepository interface {
	Send(message *gomail.Message) error
}
