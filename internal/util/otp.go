package util

import (
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPKey() (string, string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(3) * time.Minute)

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      os.Getenv("OTP_ISSUER"),
		AccountName: os.Getenv("OTP_ACCOUNT_NAME"),
		Algorithm:   otp.AlgorithmSHA512,
		Digits:      otp.DigitsSix,
		Period:      uint(expirationTime.Second()),
	})
	if err != nil {
		return "", "", time.Time{}, err
	}

	passcode, _ := totp.GenerateCodeCustom(key.Secret(), expirationTime, totp.ValidateOpts{
		Period:    uint(expirationTime.Second()),
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	})

	return passcode, key.Secret(), expirationTime, nil
}

func VerifyOTP(userOTP, secret string, expirationTime time.Time) (bool, error) {
	ops := totp.ValidateOpts{
		Period:    uint(expirationTime.Second()),
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	}

	return totp.ValidateCustom(userOTP, secret, expirationTime, ops)
}
