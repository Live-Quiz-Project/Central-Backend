package v1

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	return &service{
		Repository: repo,
		timeout:    time.Duration(3) * time.Second,
	}
}

// ---------- Auth related service methods ---------- //
func (s *service) LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(c, req.Email)
	if err != nil {
		return &LogInResponse{}, "", err
	}

	err = util.CheckPassword(u.Password, req.Password)
	if err != nil {
		return &LogInResponse{}, "", err
	}

	accessToken, err := util.GenerateToken(u.ID, time.Now().Add(24*time.Hour), os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return &LogInResponse{}, "", err
	}

	refreshToken, err := util.GenerateToken(u.ID, time.Now().Add(7*24*time.Hour), os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return &LogInResponse{}, "", err
	}

	return &LogInResponse{
		ID:          u.ID,
		Name:        u.Name,
		Image:       u.Image,
		AccessToken: accessToken,
	}, refreshToken, nil
}

func (s *service) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, string, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	u, err := s.Repository.GetUserByEmail(c, req.Email)
	if err != nil {
		return &CreateUserResponse{}, "", err
	}

	if u != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return &CreateUserResponse{}, "", err
	}

	formattedName := util.AbbreviateName(req.Name)

	r, err := s.Repository.CreateUser(c, &User{
		ID:            uuid.New(),
		Name:          req.Name,
		Email:         req.Email,
		Password:      hashedPassword,
		Image:         "default.png",
		DisplayName:   formattedName,
		DisplayEmoji:  util.SmileyFace,
		DisplayColor:  util.Gray,
		AccountStatus: util.Active,
	})
	if err != nil {
		return &CreateUserResponse{}, "", err
	}

	res, refreshToken, err := s.LogIn(ctx, &LogInRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return &CreateUserResponse{}, "", err
	}

	return &CreateUserResponse{
		UserResponse: UserResponse{
			ID:           r.ID,
			Name:         r.Name,
			Email:        r.Email,
			Password:     r.Password,
			Image:        r.Image,
			DisplayName:  r.DisplayName,
			DisplayEmoji: r.DisplayEmoji,
			DisplayColor: r.DisplayColor,
		},
		AccessToken: res.AccessToken,
	}, refreshToken, nil
}

func (s *service) GetUsers(ctx context.Context) ([]UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	users, err := s.Repository.GetUsers(c)
	if err != nil {
		return nil, err
	}

	var res []UserResponse
	for _, u := range users {
		res = append(res, UserResponse{
			ID:           u.ID,
			Name:         u.Name,
			Email:        u.Email,
			Password:     u.Password,
			Image:        u.Image,
			DisplayName:  u.DisplayName,
			DisplayEmoji: u.DisplayEmoji,
			DisplayColor: u.DisplayColor,
		})
	}

	return res, nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByID(c, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return &UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		Image:        user.Image,
		DisplayName:  user.DisplayName,
		DisplayEmoji: user.DisplayEmoji,
		DisplayColor: user.DisplayColor,
	}, nil
}

func (s *service) UpdateUser(ctx context.Context, req *UpdateUserRequest, uid uuid.UUID) (*UserResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByID(c, uid)
	if err != nil {
		return &UserResponse{}, err
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return &UserResponse{}, err
		}
		user.Password = hashedPassword
	}
	if req.Image != "" {
		user.Image = req.Image
	}
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.DisplayEmoji != "" {
		user.DisplayEmoji = req.DisplayEmoji
	}
	if req.DisplayColor != "" {
		user.DisplayColor = req.DisplayColor
	}

	r, err := s.Repository.UpdateUser(c, user)
	if err != nil {
		return &UserResponse{}, err
	}

	return &UserResponse{
		ID:           r.ID,
		Name:         r.Name,
		Email:        r.Email,
		Password:     r.Password,
		Image:        r.Image,
		DisplayName:  r.DisplayName,
		DisplayEmoji: r.DisplayEmoji,
		DisplayColor: r.DisplayColor,
	}, nil
}

func (s *service) DeleteUser(ctx context.Context, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.DeleteUser(c, uid)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GoogleSignIn(ctx context.Context, idToken string) (*LogInResponse, string, error) {
	// Verify the Google ID token and extract user info
	tokenInfo, err := util.VerifyGoogleIDToken(idToken)
	if err != nil {
		return nil, "", err
	}

	user, err := s.Repository.GetUserByGoogleID(ctx, tokenInfo.GoogleID)
	if err != nil {
		return nil, "", err
	}

	// If the user doesn't exist, create a new one
	if user == nil {
		formattedName := util.AbbreviateName(tokenInfo.Name)

		newUser := &User{
			ID:            uuid.New(),
			GoogleId:      tokenInfo.GoogleID,
			Name:          tokenInfo.Name,
			Email:         tokenInfo.Email,
			Image:         "default.png",
			DisplayName:   formattedName,
			DisplayEmoji:  util.SmileyFace,
			DisplayColor:  util.Gray,
			AccountStatus: util.Active,
		}

		// Check again before inserting to avoid duplicate key error
		existingUser, err := s.Repository.GetUserByGoogleID(ctx, newUser.GoogleId)
		if err != nil {
			return nil, "", err
		}

		if existingUser == nil {
			user, err = s.Repository.CreateUser(ctx, newUser)
			if err != nil {
				return nil, "", err
			}
		} else {
			// If the user was created by another concurrent request, use the existing user
			user = existingUser
		}
	}

	accessToken, err := util.GenerateToken(user.ID, time.Now().Add(24*time.Hour), os.Getenv("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return &LogInResponse{}, "", err
	}

	refreshToken, err := util.GenerateToken(user.ID, time.Now().Add(7*24*time.Hour), os.Getenv("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return &LogInResponse{}, "", err
	}

	// Generate JWT token for the user
	// jwtToken := ... your logic to generate JWT ...
	return &LogInResponse{
		ID:          user.ID,
		Name:        user.Name,
		Image:       user.Image,
		AccessToken: accessToken,
	}, refreshToken, nil
}

// ---------- Admin related service methods ---------- //
func (s *service) RestoreUser(ctx context.Context, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.RestoreUser(c, uid)
	if err != nil {
		return err
	}

	return nil
}
