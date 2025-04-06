package pkg

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

func GenerateToken() string {
	return uuid.NewString()
}

func SendVerificationEmail(to, token string) error {
	mail := gomail.NewMessage()

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")

	mail.SetHeader("From", "LIDM AAMIIN <yazid.al2418@gmail.com>")
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Email Verification")

	baseUrl := os.Getenv("BASE_URL")
	verifyLink := fmt.Sprintf("%s/auth/verify/%s", baseUrl, token)
	mail.SetBody("text/html", fmt.Sprintf("Click here to verify: <a href='%s'>%s</a>", verifyLink, verifyLink))

	dailer := gomail.NewDialer(smtpHost, 587, smtpUser, smtpPass)

	err := dailer.DialAndSend(mail)

	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	log.Print("Email sent to ", to)

	return err
}
