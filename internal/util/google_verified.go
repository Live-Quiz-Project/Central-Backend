package util

import (
	"context"
	"os"

	"google.golang.org/api/idtoken"
)

// GoogleUserInfo represents the basic user information extracted from the Google ID token
type GoogleUserInfo struct {
	Email         string
	VerifiedEmail bool
	Name          string
	PictureURL    string
	GoogleID      string
}

// VerifyGoogleIDToken verifies the Google ID token and returns the user's information
func VerifyGoogleIDToken(idToken string) (*GoogleUserInfo, error) {
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	payload, err := idtoken.Validate(context.Background(), idToken, client_id)
	if err != nil {
		return nil, err
	}

	userInfo := &GoogleUserInfo{
		Email:         payload.Claims["email"].(string),
		VerifiedEmail: payload.Claims["email_verified"].(bool),
		Name:          payload.Claims["name"].(string),
		PictureURL:    payload.Claims["picture"].(string),
		GoogleID:      payload.Subject,
	}

	return userInfo, nil
}
