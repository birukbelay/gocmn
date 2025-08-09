package smtp

import (
	"embed"

	"github.com/birukbelay/gocmn/src/provider/email"
	"github.com/birukbelay/gocmn/src/provider/email/templates"
)

type SmtpSender struct {
	Host     string
	Port     string
	Password string
	From     string
	fs       embed.FS
}

func NewSmtp(host, port, pwd, from string, fs embed.FS) *SmtpSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
		fs:       fs,
	}
}

func NewSmtpSender(host, port, pwd, from string, fs embed.FS) email.EmailSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
	}
}

// SendEmail accepths template path & parses it and sends it as email
func (h *SmtpSender) SendEmail(to, subject string, templatePath email.EmailTemplates, templateStruct any) error {

	f, err := h.fs.Open(templatePath.S())
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

//=====================    Verification Code Related. ================

func NewVerificationCodeSender(host, port, pwd, from string, fs embed.FS) email.VerificationSender {
	return &SmtpSender{
		Host:     host,
		Port:     port,
		Password: pwd,
		From:     from,
		fs:       fs,
	}
}

// VerificationCodeData holds the data for the email template
type VerificationCodeData struct {
	Name string
	Code string
}

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendVerificationCode(to string, code string) error {

	f, err := h.fs.Open(email.VerificationCodeTemplate.S())
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
