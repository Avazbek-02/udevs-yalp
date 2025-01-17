package etc


import (
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

type Otp struct {
	Code string
}

func GenerateNotificationEmailBody(message string) (string, error) {
	templateString := `
<!DOCTYPE html>
<html>
<body>
    <p>Hello: you have sms: {{.Code}},</p>
</body>
</html>
`
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}
	otpData := Otp{message}

	var builder strings.Builder
	err = tmpl.Execute(&builder, otpData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return builder.String(), nil
}

// generateEmailBody generates the HTML email body for a student
func GenerateOtpEmailBody(otp string) (string, error) {
	templateString := `
<!DOCTYPE html>
<html>
<body>
    <p>Your Otp to verify your YELP account {{.Code}},</p>
</body>
</html>
`
	tmpl, err := template.New("email").Parse(templateString)
	if err != nil {
		return "", fmt.Errorf("failed to parse email template: %w", err)
	}
	otpData := Otp{otp}

	var builder strings.Builder
	err = tmpl.Execute(&builder, otpData)
	if err != nil {
		return "", fmt.Errorf("failed to execute email template: %w", err)
	}

	return builder.String(), nil
}

// sendEmail sends an email using SMTP
func SendEmail(smtpHost, smtpPort, from, password, to, body string) error {
	auth := smtp.PlainAuth("", from, password, smtpHost)

	msg := []byte(fmt.Sprintf("Subject: Otp code Mini twitter\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"From: %s\r\n"+
		"To: %s\r\n"+
		"\r\n%s", from, to, body))

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
