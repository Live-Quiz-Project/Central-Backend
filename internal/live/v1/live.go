package v1

import (
	"context"
	"time"

	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/google/uuid"
)

// ---------- Live Quiz Session related models ---------- //
type Session struct {
	ID                  uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	HostID              uuid.UUID  `json:"host_id" gorm:"column:host_id;type:uuid;not null"`
	QuizID              uuid.UUID  `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null"`
	Status              string     `json:"status" gorm:"column:status;not null"`
	ExemptedQuestionIDs *string    `json:"exempted_question_ids" gorm:"column:exempted_question_ids"`
	CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamptz;not null"`
	DeletedAt           *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamptz"`
}

func (Session) TableName() string {
	return "live_quiz_session"
}

type LiveQuizSession struct {
	Session
	Code    string                `json:"code"`
	Clients map[uuid.UUID]*Client `json:"clients"`
}

type Cache struct {
	ID              uuid.UUID      `json:"id"`
	QuizID          uuid.UUID      `json:"quiz_id"`
	QuestionCount   int            `json:"question_count"`
	CurrentQuestion int            `json:"current_question"`
	Status          string         `json:"status"`
	Config          Configurations `json:"config"`
}

type Configurations struct {
	ParticipantConfig ParticipantConfigurations `json:"participant"`
	LeaderboardConfig LeaderboardConfigurations `json:"leaderboard"`
	OptionConfig      OptionConfigurations      `json:"option"`
}

type ParticipantConfigurations struct {
	Reanswer bool `json:"reanswer"`
}

type LeaderboardConfigurations struct {
	DuringQuestions  bool `json:"during_questions"`
	BetweenQuestions bool `json:"between_questions"`
}

type OptionConfigurations struct {
	Colorless         bool `json:"colorless"`
	ShowCorrectAnswer bool `json:"show_correct_answer"`
}

// ---------- Participant related models ---------- //
type Participant struct {
	ID                uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID            *uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid"`
	LiveQuizSessionID uuid.UUID  `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null"`
	Status            string     `json:"status" gorm:"column:status;type:text;not null"`
	Name              string     `json:"name" gorm:"column:name;type:text"`
	Marks             int        `json:"marks" gorm:"column:marks;type:int"`
}

func (Participant) TableName() string {
	return "participant"
}

type Repository interface {
	// ---------- Live Quiz Session related repository methods ---------- //
	CreateLiveQuizSession(ctx context.Context, lqs *Session) (*Session, error)
	GetLiveQuizSessions(ctx context.Context) ([]LiveQuizSession, error)
	GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*LiveQuizSession, error)
	GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*LiveQuizSession, error)
	GetLiveQuizSessionByCode(ctx context.Context, code string) (*LiveQuizSession, error)
	UpdateLiveQuizSession(ctx context.Context, lqs *LiveQuizSession, id uuid.UUID) (*LiveQuizSession, error)
	EndLiveQuizSession(ctx context.Context, id uuid.UUID) error
	DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error

	CreateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error
	CreateLiveQuizSessionResponseCache(ctx context.Context, code string, order int, response any) error
	GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error)
	UpdateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error
	UpdateLiveQuizSessionResponseCache(ctx context.Context, code string, order int, response any) error

	// ---------- Participant related repository methods ---------- //
	CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error)
	GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error)
	DoesParticipantExists(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateParticipantStatus(ctx context.Context, userID uuid.UUID, lqsID uuid.UUID, status string) (*Participant, error)
}

type LiveQuizSessionResponse struct {
	ID     uuid.UUID `json:"id"`
	HostID uuid.UUID `json:"host_id"`
	QuizID uuid.UUID `json:"quiz_id"`
	Code   string    `json:"code"`
	Status string    `json:"status"`
}

type CreateLiveQuizSessionRequest struct {
	QuizID uuid.UUID `json:"quiz_id"`
}
type CreateLiveQuizSessionResponse struct {
	ID     uuid.UUID `json:"id"`
	QuizID uuid.UUID `json:"quiz_id"`
	Code   string    `json:"code"`
}

type ParticipantsResponse struct {
	Participants []Participant `json:"participants"`
}

type UpdateLiveQuizSessionRequest struct {
	Status          string  `json:"status"`
	ExemptedQuesIDs *string `json:"exempted_question_ids"`
}

type CheckLiveQuizSessionAvailabilityResponse struct {
	ID              uuid.UUID `json:"id"`
	QuizID          uuid.UUID `json:"quiz_id"`
	Code            string    `json:"code"`
	QuestionCount   int       `json:"question_count"`
	CurrentQuestion int       `json:"current_question"`
	Status          string    `json:"status"`
}

type CountDownPayload struct {
	LiveQuizSessionID uuid.UUID `json:"live_quiz_session_id"`
	TimeLeft          int       `json:"time_left"`
	QuestionCount     int       `json:"question_count"`
	CurrentQuestion   int       `json:"current_question"`
	Status            string    `json:"status"`
}

type Service interface {
	// ---------- Live Quiz Session related service methods ---------- //
	CreateLiveQuizSession(ctx context.Context, req *CreateLiveQuizSessionRequest, id uuid.UUID, code string, hostID uuid.UUID) (*CreateLiveQuizSessionResponse, error)
	GetLiveQuizSessions(ctx context.Context) ([]LiveQuizSessionResponse, error)
	GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*LiveQuizSessionResponse, error)
	UpdateLiveQuizSession(ctx context.Context, req *UpdateLiveQuizSessionRequest, id uuid.UUID) (*LiveQuizSessionResponse, error)
	DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error

	CreateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error
	GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error)

	GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*q.Question, error)

	// ---------- Participant related service methods ---------- //
	CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error)
	GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) (*ParticipantsResponse, error)
	DoesParticipantExists(ctx context.Context, userID uuid.UUID) (bool, error)
	UpdateParticipantStatus(ctx context.Context, userID uuid.UUID, lqsID uuid.UUID, status string) (*Participant, error)
}
