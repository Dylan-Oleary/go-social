package mailer

import (
	"bytes"
	"errors"

	"text/template"

	gomail "gopkg.in/mail.v2"
)

type MailtrapMailer struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapClient(apiKey, fromEmail string) (*MailtrapMailer, error) {
	if apiKey == "" {
		return &MailtrapMailer{}, errors.New("api key is required")
	}

	return &MailtrapMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m MailtrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	return 200, nil
}
