package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/go-gomail/gomail"
)

func SendConfirmationCode(toEmail, confirmationCode string) error {
	d := gomail.NewDialer("smtp.gmail.com", 587, "lq.platform59@gmail.com", "qapy brza gciy cezt")

	m := gomail.NewMessage()
	m.SetHeader("From", "lq.platform59@gmail.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Confirmation Code")
	m.SetBody("text/plain", fmt.Sprintf("Your confirmation code is: %s", confirmationCode))

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func GenerateConfirmationCode() (string, error) {
	randomBytes := make([]byte, 6)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	confirmationCode := base64.URLEncoding.EncodeToString(randomBytes)

	return confirmationCode, nil
}
