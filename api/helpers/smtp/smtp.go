package helpers

import (
	"fmt"
	"net/smtp"
)

type SMTPHelper struct {
	Host        string
	Port        string
	Username    string
	Password    string
	MailDetails MailDetails
}

type MailDetails struct {
	Receivers []string
	Subject   string
	Data      []byte
	Sender    string
	mimeType  string
}

func NewSMTPHelper(host, port, username, password string) SMTPHelper {
	return SMTPHelper{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

func (sh *SMTPHelper) SetSender(sender string) {
	sh.MailDetails.Sender = sender
}

func (sh *SMTPHelper) SetReceivers(receivers []string) {
	sh.MailDetails.Receivers = receivers
}

func (sh *SMTPHelper) SetSubject(subject string) {
	sh.MailDetails.Subject = fmt.Sprintf("Subject: %s%s", subject, "\n\n")
}

func (sh *SMTPHelper) SetPlainBody(body []byte) {
	sh.MailDetails.Data = body
	sh.MailDetails.mimeType = "Content-Type: text/plain; charset=\"UTF-8\";\n\n"
}

func (sh *SMTPHelper) SetHTMLBody(body []byte) {
	sh.MailDetails.Data = body
	sh.MailDetails.mimeType = "Content-Type: text/html; charset=\"UTF-8\";\n\n"
}

func (sh *SMTPHelper) SendMail() error {
	auth := smtp.PlainAuth("", sh.Username, sh.Password, sh.Host)

	sh.MailDetails.Data = append([]byte(sh.MailDetails.Subject), sh.MailDetails.Data...)

	if err := smtp.SendMail(fmt.Sprintf("%s:%s", sh.Host, sh.Port), auth, sh.MailDetails.Sender, sh.MailDetails.Receivers, sh.MailDetails.Data); err != nil {
		return err
	}
	return nil
}
