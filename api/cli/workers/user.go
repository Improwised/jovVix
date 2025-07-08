package workers

import (
	"log"

	helpers "github.com/Improwised/jovvix/api/helpers/smtp"
)

type WelcomeMail struct {
	FirstName string `json:"first_name" `
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
}

func (w WelcomeMail) Handle() error {
	log.Println("sending mail")

	smtp := helpers.NewSMTPHelper("localhost", "2525", "root", "pass")
	smtp.SetSubject("welocme")
	smtp.SetPlainBody([]byte("welcome to our org"))

	smtp.SetSender("support@improwised.com")
	smtp.SetReceivers([]string{w.Email})

	if err := smtp.SendMail(); err != nil {
		return err
	}

	log.Printf("mail send to %v", w.Email)
	return nil
}
