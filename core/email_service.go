package core

import (
	"gopkg.in/gomail.v2"

	"log"
)

type EmailService interface {
	Message(message string, usermail string) error // Message sends an email
}

type EmailServiceImpl struct {
	repo EmailRepository
}

func NewEmailService(repo EmailRepository) EmailService {
	return &EmailServiceImpl{repo: repo}
}

func (e *EmailServiceImpl) Message(message string, usermail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "semailsender01@gmail.com")
	m.SetHeader("To", usermail)
	m.SetHeader("Subject", "OTP Verification TiawPao app")
	m.SetBody("text/html", "<h1>Your One Time Password (OTP) is : "+message+"</h1>")

	err := e.repo.Send(m)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}
	return nil
}
