package smtp

import (
	"os"
	"path/filepath"

	"github.com/birukbelay/gocmn/src/provider/email"
)

//==================.  unused, Might be use full for making this package as a library

// SendEmail sends an email using Gmail SMTP
func (h *SmtpSender) SendVerificationCodeEmbeded(to string, code string) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	//TODO: we might give this variable as a param
	fullPath := filepath.Join(wd, "src/provider/email/templates/verification_code.html")
	emailFields := email.EmailFields{
		To:       to,
		From:     h.From,
		Subject:  "Verify Your Email",
		HtmlPath: fullPath,
	}
	verificationData := VerificationCodeData{
		Code: code,
		Name: "Mr",
	}
	return h.SendEmailTmpl(emailFields, verificationData)
}

// SendEmail given htmlPath
func (h *SmtpSender) SendEmailTmpl(fields email.EmailFields, templateStruct any) error {
	// Parse the email template
	body, err := email.ParseEmailTemplate(fields.HtmlPath, templateStruct)
	if err != nil {
		return err
	}
	return h.SendEmailT(fields, body)
}
