package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

type CreateLiveAnswerRequest struct {
	LiveQuizSessionID uuid.UUID `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null;references:live_quiz_session(id)"`
	Answers           []AnswerResponse
}

type ParticipantResponse struct {
	ID                uuid.UUID  `json:"id" `
	UserID            *uuid.UUID `json:"user_id"`
	LiveQuizSessionID *uuid.UUID  `json:"live_quiz_session_id,omitempty"`
	Status            string     `json:"status,omitempty"`
	Name              string     `json:"name"`
	Marks             int        `json:"marks,omitempty"`
}

type LiveAnswerRequest struct {
	AnswerResponse
}

type LiveAnswerResponse struct {
	AnswerResponse
}

type QuestionViewQuizResponse struct {
	ID          uuid.UUID                      `json:"id"`
	CreatorID   uuid.UUID                      `json:"creator_id"`
	Title       string                         `json:"title"`
	Description string                         `json:"description"`
	CoverImage  string                         `json:"cover_image"`
	CreatedAt   time.Time                      `json:"created_at"`
	Questions   []QuestionViewQuestionResponse `json:"questions"`
}

type QuestionViewQuestionResponse struct {
	ID             uuid.UUID     `json:"id"`
	Order          int           `json:"order"`
	Content        string        `json:"content"`
	Type           string        `json:"type"`
	Note           string        `json:"note"`
	Media          string        `json:"media"`
	UseTemplate    bool          `json:"use_template"`
	TimeLimit      int           `json:"time_limit"`
	HaveTimeFactor bool          `json:"have_time_factor"`
	TimeFactor     int           `json:"time_factor"`
	FontSize       int           `json:"font_size"`
	SelectUpTo     int           `json:"select_up_to"`
	Options        []interface{} `json:"options"`
}

type QuestionViewOptionChoice struct {
	ID           uuid.UUID     `json:"id"`
	Order        int           `json:"order"`
	Content      string        `json:"content"`
	Mark         int           `json:"mark"`
	Correct      bool          `json:"correct"`
	Participants []ParticipantResponse `json:"participants"`
}

type QuestionViewOptionText struct {
	ID            uuid.UUID     `json:"id"`
	Order         int           `json:"order"`
	Content       string        `json:"content"`
	Mark          int           `json:"mark"`
	CaseSensitive bool          `json:"case_sensitive"`
	Participants  []ParticipantResponse `json:"participants"`
}

type QuestionViewMatching struct {
	ID       uuid.UUID `json:"id"`
	OptionID uuid.UUID `json:"option_id"`
	OptionContent string `json:"option_content"`
	PromptID uuid.UUID `json:"prompt_id"`
	PromptContent string `json:"prompt_content"`
	Mark     int       `json:"mark"`
	Participants []ParticipantResponse
}

// type QuestionViewOptionMatching struct {
// 	ID        uuid.UUID `json:"id"`
// 	Type      string    `json:"type"`
// 	Order     int       `json:"order"`
// 	Content   string    `json:"content"`
// 	Eliminate string    `json:"eliminate"`
// }

// type QuestionViewAnswerMatching struct {
	
// }

type QuestionViewParticipant struct {
	ID   uuid.UUID
	Name string
}

type Ranking struct {
	ID        uuid.UUID
	Order     int
	Name      string
	Marks     string
	Correct   int
	Incorrect int
	Timeup    int
	TimeUse   int
}

// -------------------- REPOSITORY START --------------------
type Repository interface {
	// Transaction
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error

	// CREATE
	CreateAnswerResponse(ctx context.Context, tx *gorm.DB, liveAnswer *AnswerResponse) (*AnswerResponse, error)

	// GET

	GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(ctx context.Context, liveQuizSessionID uuid.UUID, questionID uuid.UUID) ([]AnswerResponse, error)
	GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveSessionID uuid.UUID) ([]AnswerResponse, error)
	GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]AnswerResponse, error)
	GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]AnswerResponse, error)
	GetParticipantByID(ctx context.Context, participantID uuid.UUID) (*Participant, error)
}

// #################### SERVICE START ####################
type Service interface {
	// CREATE
	CreateAnswerResponse(ctx context.Context, req *LiveAnswerRequest) (*LiveAnswerResponse, error)

	// GET
	GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveSessionID uuid.UUID) ([]LiveAnswerResponse, error)
	GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]LiveAnswerResponse, error)
	GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]LiveAnswerResponse, error)
	GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(ctx context.Context, liveQuizSessionID uuid.UUID, questionID uuid.UUID) ([]LiveAnswerResponse, error)
	GetParticipantByID(ctx context.Context, liveQuizSessionID uuid.UUID) (*Participant, error)
}
