package services

import (
	"strconv"

	"github.com/Improwised/quizz-app/api/config"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	logger *zap.Logger
	config *config.SMTPConfig
}

func NewEmailService(logger *zap.Logger, config *config.SMTPConfig) *EmailService {

	portEnv := config.SmtpPort
	port := 1025

	if portEnv != "" {
		if parsedPort, err := strconv.Atoi(portEnv); err == nil {
			port = parsedPort
		}
	}

	dialer := gomail.NewDialer(config.SmtpHost, port, config.SmtpUsername, config.SmtpPassword)
	dialer.SSL = false
	dialer.TLSConfig = nil

	return &EmailService{
		dialer: dialer,
		logger: logger,
		config: config,
	}
}

// SendEmail sends an email with the specified subject and body to the recipient
func (es *EmailService) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", es.config.EmailFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := es.dialer.DialAndSend(m); err != nil {
		return err
	}
	es.logger.Info("Email sent successfully.", zap.Any("to", to))

	return nil
}
