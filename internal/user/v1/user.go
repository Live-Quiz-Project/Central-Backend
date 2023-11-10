package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- User related models ---------- //
type User struct {
	ID            uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Name          string         `json:"name" gorm:"column:name;type:text"`
	Email         string         `json:"email" gorm:"column:email;type:text;unique"`
	Password      string         `json:"password" gorm:"column:password;type:text"`
	Image         string         `json:"image" gorm:"column:image;type:text"`
	DisplayName   string         `json:"display_name" gorm:"column:display_name;type:text"`
	DisplayEmoji  string         `json:"display_emoji" gorm:"column:display_emoji;type:text"`
	DisplayColor  string         `json:"display_color" gorm:"column:display_color;type:text"`
	AccountStatus string         `json:"account_status" gorm:"column:account_status;type:text"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (User) TableName() string {
	return "user"
}

// ---------- Admin related models ---------- //
type Admin struct {
	ID       uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Email    string    `json:"email" gorm:"column:email;type:text;unique"`
	Password string    `json:"password" gorm:"column:password;type:text"`
}

func (Admin) TableName() string {
	return "admin"
}

type Repository interface {
	// ---------- User related repository methods ---------- //
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// ---------- Admin related repository methods ---------- //
	RestoreUser(ctx context.Context, id uuid.UUID) error
}

// ---------- Auth related structs ---------- //
type LogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LogInResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	AccessToken string    `json:"accessToken"`
}
type DecodeTokenRequest struct {
	Token string `json:"token"`
}

// ---------- User related structs ---------- //
type UserResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Image        string    `json:"image"`
	DisplayName  string    `json:"display_name"`
	DisplayEmoji string    `json:"display_emoji"`
	DisplayColor string    `json:"display_color"`
}
type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type CreateUserResponse struct {
	UserResponse
	AccessToken string `json:"accessToken"`
}
type UpdateUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Image        string `json:"image"`
	DisplayName  string `json:"display_name"`
	DisplayEmoji string `json:"display_emoji"`
	DisplayColor string `json:"display_color"`
}

// ---------- Admin related structs ---------- //

type Service interface {
	// ---------- Auth related service methods ---------- //
	LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, string, error)

	// ---------- User related service methods ---------- //
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, string, error)
	GetUsers(ctx context.Context) ([]UserResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*UserResponse, error)
	UpdateUser(ctx context.Context, req *UpdateUserRequest, id uuid.UUID) (*UserResponse, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// ---------- Admin related service methods ---------- //
	RestoreUser(ctx context.Context, id uuid.UUID) error
}
