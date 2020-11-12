package email

import (
	"bytes"
	"text/template"
	"vl_sa/logger"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// ClientInterface used for sending emails
type ClientInterface interface {
	SendMail(content interface{}, mailTemplate string, recipient string, subject string) error
}

//SendgridClient email client implementation
type SendgridClient struct {
	Settings Settings
}

// SendMail implementation
func (client *SendgridClient) SendMail(content interface{}, mailTemplate string, recipient string, subject string) (err error) {
	from := mail.NewEmail("vl_sa", client.Settings.FromAddress)
	to := mail.NewEmail(recipient, recipient)
	plainTextContent := "A"

	buf := new(bytes.Buffer)
	template, err := template.New("webpage").Parse(mailTemplate)
	if err := template.Execute(buf, content); err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, buf.String())
	sendgrid := sendgrid.NewSendClient(client.Settings.APIKey)
	response, err := sendgrid.Send(message)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 && response.StatusCode > 204 {
		logger.Error.Printf("failed sending order request: %d", response.StatusCode)
	}

	return
}
