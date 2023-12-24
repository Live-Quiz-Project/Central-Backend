package util

import (
	"os"
	"time"

	"github.com/pquerna/otp/totp"
)

func GenerateTOTPKey() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Getenv("OTP_ISSUER"),
		AccountName: os.Getenv("OTP_ACCOUNT_NAME"),
	})
	if err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(time.Duration(5) * time.Minute)

	code, err := totp.GenerateCode(key.Secret(), expirationTime)
	if err != nil {
		return "", err
	}

	return code, nil
}

func GetOTPSecret(issuer, accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

func ValidateTOTP(userOTP, secret string) bool {
	return totp.Validate(userOTP, secret)
}
