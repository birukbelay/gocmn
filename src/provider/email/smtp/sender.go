package smtp

import (
	"bytes"
	"net/smtp"

	"github.com/birukbelay/gocmn/src/provider/email"
)


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
