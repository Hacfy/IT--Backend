package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Hacfy/IT_INVENTORY/pkg/templates"
	"gopkg.in/gomail.v2"
)

func SendLoginCredentials(email string, password string) error {

	smtpHost := os.Getenv("SMTP_HOST")

	if smtpHost == "" {
		return errors.New("SMTP_HOST env was missing")
	}

	smtpPort := os.Getenv("SMTP_PORT")

	if smtpPort == "" {
		return errors.New("SMTP_PORT env was missing")
	}

	hostEmail := os.Getenv("HOST_EMAIL")

	if hostEmail == "" {
		return errors.New("HOST_EMAIL env was missing")
	}

	appPassword := os.Getenv("APP_PASSWORD")

	if appPassword == "" {
		return errors.New("APP_PASSWORD env was missing")
	}

	to := email

	smtpPortInt, err := strconv.Atoi(smtpPort)

	if err != nil {
		return err
	}

	client := gomail.NewDialer(smtpHost, smtpPortInt, hostEmail, appPassword)

	htmlTemplate := templates.GetVerifyEmailOtpTemplate(password, email)

	message := gomail.NewMessage()

	message.SetHeader("From", hostEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Your Secure Login Credentials is Ready")
	message.SetBody("text/html", htmlTemplate)

	if err = client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func SendForgotPasswordEmail(email, otp string) error {

	smtpHost := os.Getenv("SMTP_HOST")

	if smtpHost == "" {
		return errors.New("SMTP_HOST env was missing")
	}

	smtpPort := os.Getenv("SMTP_PORT")

	if smtpPort == "" {
		return errors.New("SMTP_PORT env was missing")
	}

	hostEmail := os.Getenv("HOST_EMAIL")

	if hostEmail == "" {
		return errors.New("HOST_EMAIL env was missing")
	}

	appPassword := os.Getenv("APP_PASSWORD")

	if appPassword == "" {
		return errors.New("APP_PASSWORD env was missing")
	}

	to := email

	smtpPortInt, err := strconv.Atoi(smtpPort)

	if err != nil {
		return err
	}

	client := gomail.NewDialer(smtpHost, smtpPortInt, hostEmail, appPassword)

	htmlTemplate := templates.GetForgotPasswordTemplate(email)

	message := gomail.NewMessage()

	message.SetHeader("From", hostEmail)
	message.SetHeader("To", to)
	message.SetHeader("Subject", "Reset Your Password")
	message.SetBody("text/html", htmlTemplate)

	if err = client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
