package config

type SMTPConfig struct {
	SmtpHost     string `envconfig:"SMTP_HOST"`
	SmtpPort     string `envconfig:"SMTP_PORT"`
	SmtpUsername string `envconfig:"SMTP_USERNAME"`
	SmtpPassword string `envconfig:"SMTP_PASSWORD"`
	EmailFrom    string `envconfig:"EMAIL_FROM"`
}
