package mail

import (
	"net/smtp"

	"github.com/Tutuacs/pkg/config"
)

type MailerAuth struct {
	srv   string
	auth  smtp.Auth
	email string
}

var mailerAuth *MailerAuth

func init() {
	mailerAuth = nil
}

func UseMailer() *MailerAuth {
	if mailerAuth == nil {

		conf := config.GetMailer()

		auth := smtp.PlainAuth(
			"",
			conf.SMTP_MAIL,
			conf.SMTP_PASS,
			conf.SMTP_HOST,
		)

		mailerAuth = &MailerAuth{
			srv:   conf.SMTP_ADDR,
			auth:  auth,
			email: conf.SMTP_MAIL,
		}
	}
	return mailerAuth
}

func (m *MailerAuth) SendEmail(to []string, subject string, body string) error {
	if len(subject) == 0 {
		subject = "Email from Gbase"
	}
	msg := []byte("Subject: " + subject + "\n" + body)

	return smtp.SendMail(m.srv, m.auth, m.email, to, msg)
}
