package smtp

import (
	"bytes"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/birukbelay/gocmn/src/provider/email"
	"github.com/birukbelay/gocmn/src/provider/email/templates"
)

// SendEmail accepths template path & parses it and sends it as email
func (h *SmtpSender) SendEmail(to, subject string, templatePath email.EmailTemplates, templateStruct any) error {

	f, err := templates.Embedded.Open(templatePath.S())
	if err != nil {
		panic(err)
	}

	body, err := email.ParseEmbededTemplate(f, templateStruct)
	if err != nil {
		return err
	}

	emailFields := email.EmailFields{
		To:      to,
		From:    h.From,
		Subject: subject,
	}

	return h.SendEmailT(emailFields, body)
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendEmailT(fields email.EmailFields, body string) error {

	// Set up SMTP server configuration
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"
	auth := smtp.PlainAuth("", fields.From, h.Password, h.Host)

	// Create the email message
	headers := make(map[string]string)
	headers["From"] = fields.From
	headers["To"] = fields.To
	headers["Subject"] = fields.Subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	// Build the message
	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// Send the email
	err := smtp.SendMail(
		h.Host+":"+h.Port,
		auth,
		fields.From,
		[]string{fields.To},
		msg.Bytes(),
	)
	return err
}

