package pkg

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func GenerateToken() string {
	return uuid.NewString()
}

func SendVerificationEmail(to, token string) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	mail := gomail.NewMessage()

	from := smtpUser

	mail.SetHeader("From", from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Email Verification")

	baseUrl := os.Getenv("BASE_URL")
	verifyLink := fmt.Sprintf("%s/verify?token=%s", baseUrl, token)
	mail.SetBody("text/html", fmt.Sprintf("Click here to verify: <a href='%s'>%s</a>", verifyLink, verifyLink))

	dailer := gomail.NewDialer(smtpHost, 587, smtpUser, smtpPass)

	err := dailer.DialAndSend(mail)

	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	return err
}
