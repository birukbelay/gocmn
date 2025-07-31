package smtp

import (
	"github.com/birukbelay/gocmn/src/provider/email"
	"github.com/birukbelay/gocmn/src/provider/email/templates"
)

type SmtpSender struct {
	Host     string
	Port     string
	Password string
	From     string
}

func NewSmtp(host, port, pwd, from string) *SmtpSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

func NewSmtpSender(host, port, pwd, from string) email.EmailSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

//=====================    Verification Code Related

func NewVerificationCodeSender(host, port, pwd, from string) email.VerificationSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

// VerificationCodeData holds the data for the email template
type VerificationCodeData struct {
	Name string
	Code string
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendVerificationCode(to string, code string) error {

	f, err := templates.Embedded.Open(email.VerificationCodeTemplate.S())
	if err != nil {
		panic(err)
	}
	verificationData := VerificationCodeData{
		Code: code,
		Name: "Mr",
	}
	body, err := email.ParseEmbededTemplate(f, verificationData)
	if err != nil {
		return err
	}

	emailFields := email.EmailFields{
		To:      to,
		From:    h.From,
		Subject: "Verify Your Email",
	}

	return h.SendEmailT(emailFields, body)
}
