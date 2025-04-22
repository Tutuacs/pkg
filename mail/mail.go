package mail

import (
	"fmt"
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

		// conf := config.GetMailer()

		auth := smtp.PlainAuth(
			"Sistema de Propostas Comerciais",
			"arthursilva.mailtest@gmail.com",
			"xwrf ywms jycc dqnn",
			"smtp.gmail.com",
		)

		mailerAuth = &MailerAuth{
			srv:   "smtp.gmail.com:587",
			auth:  auth,
			email: "arthursilva.mailtest@gmail.com",
		}
	}
	return mailerAuth
}

func (m *MailerAuth) SendEmail(to []string, subject string, body string) error {
	if len(subject) == 0 {
		subject = "Email from Gbase"
	}

	tmpl := `
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>` + subject + `</title>
		<style>
			* {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
			}
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				color: #333;
				padding: 20px;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 5px;
				box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
			}
			h1 {
				color: #333;
			}
			p {
				margin-bottom: 15px;
			}
			.footer {
				margin-top: 20px;
				font-size: 12px;
				color: #777;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>` + subject + `</h1>
			<p>` + body + `</p>
			<p class="footer">Este é um email automático, por favor não responda.</p>
		</div>
	</body>
	</html>
	`

	msg := "Subject: " + subject
	mime := ";\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(msg + mime + tmpl)
	return smtp.SendMail(m.srv, m.auth, m.email, to, message)
}

func (m *MailerAuth) SendForgotEmail(to string, token string) error {

	frontend := config.GetClient()

	tmpl := `
	<!DOCTYPE html>
	<html lang="pt-BR">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title> Forgot pass on gBase app </title>
		<style>
			* {
				margin: 0;
				padding: 0;
				box-sizing: border-box;
			}
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				color: #333;
				padding: 20px;
			}
			.container {
				max-width: 600px;
				margin: 0 auto;
				background-color: #fff;
				padding: 20px;
				border-radius: 5px;
				box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
			}
			h1 {
				color: #333;
			}
			p {
				margin-bottom: 15px;
			}
			.footer {
				margin-top: 20px;
				font-size: 12px;
				color: #777;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1> Forgot pass on gBase app </h1>
			<a href="` + fmt.Sprintf(`%s"`, frontend) + `/reset-password?token=` + fmt.Sprintf(`%s"`, token) + `">
				<p> Clique aqui para redefinir sua senha </p>
				<p> Se você não solicitou a redefinição de senha, ignore este email. </p>
				<p> Se você não conseguir clicar no link, copie e cole o seguinte URL em seu navegador: </p>
				<p> ` + fmt.Sprintf(`%s"`, frontend) + `/reset-password?token=` + fmt.Sprintf(`%s"`, token) + ` </p>
			</a>
			<p class="footer">Este é um email automático, por favor não responda.</p>
		</div>
	</body>
	</html>
	`

	subject := "Reset gBase Password"
	msg := "Subject: " + subject
	mime := ";\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	message := []byte(msg + mime + tmpl)

	return smtp.SendMail(m.srv, m.auth, m.email, []string{to}, message)
}
