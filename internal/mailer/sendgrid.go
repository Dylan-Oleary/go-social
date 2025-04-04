package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type SendGridMailer struct {
	apiKey    string
	client    *sendgrid.Client
	fromEmail string
	logger    *zap.SugaredLogger
}

func NewSendGrid(apiKey, fromEmail string, logger *zap.SugaredLogger) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGridMailer{
		apiKey:    apiKey,
		client:    client,
		fromEmail: fromEmail,
		logger:    logger,
	}
}

func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	from := mail.NewEmail(MailFromName, m.fromEmail)
	to := mail.NewEmail(username, email)

	// TODO: Template Parse + Build
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

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(
		&mail.MailSettings{
			SandboxMode: &mail.Setting{
				Enable: &isSandbox,
			},
		},
	)

	for i := 0; i < maxRetries; i++ {
		res, err := m.client.Send(message)
		if err != nil {
			m.logger.Errorw("Failed to send email", "recipient", email, "attempt", i+1, "maximum_attempts", maxRetries, "error", err.Error())

			// Exponential Backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		return res.StatusCode, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts", maxRetries)
}
