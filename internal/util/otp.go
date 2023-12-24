package util

import (
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func GenerateTOTPKey() (string, string, time.Time, error) {
	expirationTime := time.Now().Add(time.Duration(1) * time.Minute)

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

	ops := totp.ValidateOpts{
		Period:    uint(expirationTime.Second()),
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	}

	passcode, _ := totp.GenerateCodeCustom(key.Secret(), expirationTime, ops)

	return passcode, key.Secret(), expirationTime, nil
}

func VerifyOTP(userOTP, secret string) (bool, error) {
	expirationTime := time.Now().Add(time.Duration(1) * time.Minute)
	ops := totp.ValidateOpts{
		Period:    uint(expirationTime.Second()),
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA512,
	}

	return totp.ValidateCustom(userOTP, secret, expirationTime, ops)
}
