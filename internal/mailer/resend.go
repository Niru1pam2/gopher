package mailer

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/resend/resend-go/v2"
)

type resendClient struct {
	fromEmail string
	client    *resend.Client
}

func NewResendClient(apiKey, fromEmail string) (resendClient, error) {
	if apiKey == "" {
		return resendClient{}, errors.New("api key is required")
	}

	client := resend.NewClient(apiKey)

	return resendClient{
		fromEmail: fromEmail,
		client:    client,
	}, nil
}

func (m resendClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// 1. Template parsing and building (Exactly the same as before!)
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

	// 2. Build the Resend API request
	params := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Subject: subject.String(),
		Html:    body.String(),
	}

	// 3. Send it via the Resend API
	_, err = m.client.Emails.Send(params)
	if err != nil {
		return -1, err
	}

	return 200, nil
}
