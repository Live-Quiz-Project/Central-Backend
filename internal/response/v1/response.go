package v1

import (
	"context"

	"github.com/google/uuid"
)

// ---------- Response related models ---------- //
type ChoiceResponse struct {
	ID             uuid.UUID `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	ParticipantID  uuid.UUID `json:"participant_id" gorm:"column:participant_id;type:uuid;not null"`
	OptionChoiceID uuid.UUID `json:"option_choice_id" gorm:"column:option_choice_id;type:uuid;not null"`
}

func (ChoiceResponse) TableName() string {
	return "response_choice"
}

type Repository interface {
	// ---------- Choice response related repository methods ---------- //
	CreateChoiceResponse(ctx context.Context, r *ChoiceResponse) (*ChoiceResponse, error)
	GetChoiceResponsesByQuizID(ctx context.Context, quizID string) ([]ChoiceResponse, error)
	GetChoiceResponsesByQuizIDAndQuestionID(ctx context.Context, quizID string, questionID string) ([]ChoiceResponse, error)
	GetChoiceResponsesByQuizIDAndUserID(ctx context.Context, quizID string, uid string) ([]ChoiceResponse, error)
	GetChoiceResponsesByQuizIDAndQuestionIDAndUserID(ctx context.Context, quizID string, questionID string, uid string) ([]ChoiceResponse, error)
	UpdateChoiceResponse(ctx context.Context, r *ChoiceResponse) (*ChoiceResponse, error)
	DeleteChoiceResponse(ctx context.Context, id uuid.UUID) error
	RestoreChoiceResponse(ctx context.Context, id uuid.UUID) error
}

// ---------- Response related structs ---------- //
// Choice response related structs
type ChoiceResponseResponse struct {
	ID             uuid.UUID `json:"id"`
	ParticipantID  uuid.UUID `json:"participant_id"`
	OptionChoiceID uuid.UUID `json:"option_choice_id"`
}
type CreateChoiceResponseRequest struct{}
type CreateChoiceResponseResponse struct{}
type UpdateChoiceResponseRequest struct{}

type Service interface {
	// // ---------- Choice response related service methods ---------- //
	// CreateChoiceResponse(ctx context.Context, req *CreateChoiceResponseRequest, uid uuid.UUID) (*CreateChoiceResponseResponse, error)
	// GetChoiceResponsesByQuizID(ctx context.Context, quizID string) ([]ChoiceResponseResponse, error)
	// GetChoiceResponsesByQuestionID(ctx context.Context, quizID string, questionID string) ([]ChoiceResponseResponse, error)
	// GetChoiceResponsesByUserID(ctx context.Context, quizID string, uid string) ([]ChoiceResponseResponse, error)
	// GetChoiceResponsesByQuestionIDAndUserID(ctx context.Context, quizID string, questionID string, uid string) ([]ChoiceResponseResponse, error)
	// UpdateChoiceResponse(ctx context.Context, req *UpdateChoiceResponseRequest, uid uuid.UUID) (*ChoiceResponseResponse, error)
	// DeleteChoiceResponse(ctx context.Context, uid uuid.UUID) error
	// RestoreChoiceResponse(ctx context.Context, uid uuid.UUID) error
}
