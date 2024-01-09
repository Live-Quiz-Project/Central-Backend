package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Participant struct {
	ID                uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID            *uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid"`
	LiveQuizSessionID uuid.UUID  `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null"`
	Status            string     `json:"status" gorm:"column:status;type:text;not null"`
	Name              string     `json:"name" gorm:"column:name;type:text"`
	Marks             int        `json:"marks" gorm:"column:marks;type:int"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`

}

func (Participant) TableName() string {
	return "participant"
}

type AnswerResponse struct {
	ID                uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	LiveQuizSessionID uuid.UUID      `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null;references:live_quiz_session(id)"`
	ParticipantID     uuid.UUID      `json:"participant_id" gorm:"column:participant_id;type:uuid;not null;references:participant(id)"`
	Type              string         `json:"type" gorm:"column:type;type:text"`
	QuestionID        uuid.UUID      `json:"question_id" gorm:"column:question_id:type:uuid"`
	Answer            string         `json:"answer" gorm:"column:answer;type:text"`
	UseTime           int            `json:"use_time" gorm:"column:use_time;type:int"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (AnswerResponse) TableName() string {
	return "answer_response"
}

type LiveAnswerRequest struct {
	AnswerResponse
}

type LiveAnswerResponse struct {
	AnswerResponse
}

type ParticipantRequest struct {
	Participant
}

type ParticipantResponse struct {
	Participant
	Order 				int `json:"order"`
}

// -------------------- REPOSITORY START --------------------
type Repository interface {
	// Transaction
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error

	// GET
	GetParticipantsByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]Participant, error)
	GetParticipantsByLiveQuizSessionIDCustom(ctx context.Context, liveQuizSessionID uuid.UUID, order string, limit int) ([]Participant, error)

	GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveSessionID uuid.UUID) ([]AnswerResponse, error)
	GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]AnswerResponse, error)
	GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]AnswerResponse, error)
}

// --------------------- REPOSITORY END ---------------------

// #################### SERVICE START ####################
type Service interface {
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	CommitTransaction(ctx context.Context,tx *gorm.DB) (error)

	GetParticipantsByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]ParticipantResponse, error)
	GetParticipantsByLiveQuizSessionIDCustom(ctx context.Context, liveQuizSessionID uuid.UUID, orderBy string, limit int) ([]ParticipantResponse, error)

	GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]LiveAnswerResponse, error)
	GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]LiveAnswerResponse, error)
	GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]LiveAnswerResponse, error)
}

// ##################### SERVICE END #####################
