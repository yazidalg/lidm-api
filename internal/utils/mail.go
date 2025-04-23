package utils

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

func SendPasswordResetEmail(to, otp string) error {
	mail := gomail.NewMessage()

	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")

	mail.SetHeader("From", "LIDM AAMIIN <yazid.al2418@gmail.com>")
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", "Password Reset Request")

	messageBody := fmt.Sprintf(`
	<html>
		<body>
			<p>Hello,</p>
			<p>We received a request to reset your password. Here is your OTP code:</p>
			<h2 style="font-size: 24px; padding: 10px; background-color: #f0f0f0; border-radius: 5px; text-align: center;">%s</h2>
			<p>This code will expire in 10 minutes.</p>
			<p>If you did not request a password reset, please ignore this email.</p>
			<p>Best regards,<br>LIDM Team</p>
		</body>
	</html>
	`, otp)

	mail.SetBody("text/html", messageBody)

	dialer := gomail.NewDialer(smtpHost, 587, smtpUser, smtpPass)

	err := dialer.DialAndSend(mail)

	if err != nil {
		fmt.Println("Error sending password reset email:", err)
		return err
	}

	log.Print("Password reset email sent to ", to)

	return nil
}
