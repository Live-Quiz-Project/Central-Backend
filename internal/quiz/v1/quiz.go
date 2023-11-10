package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- Quiz related models ---------- //
type Quiz struct {
	ID          uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	CreatorID   uuid.UUID      `json:"creator_id" gorm:"column:creator_id;type:uuid;not null;references:user(id)"`
	Title       string         `json:"title" gorm:"column:title;type:text;default:Untitled"`
	Description string         `json:"description" gorm:"column:description;type:text"`
	CoverImage  string         `json:"cover_image" gorm:"column:cover_image;type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (Quiz) TableName() string {
	return "quiz"
}

type QuizHistory struct {
	ID          uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuizID      uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz(id)"`
	CreatorID   uuid.UUID      `json:"creator_id" gorm:"column:creator_id;type:uuid;not null;references:user(id)"`
	Title       string         `json:"title" gorm:"column:title;type:text;default:Untitled"`
	Description string         `json:"description" gorm:"column:description;type:text"`
	CoverImage  string         `json:"cover_image" gorm:"column:cover_image;type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted     bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy   uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (QuizHistory) TableName() string {
	return "quiz_history"
}

// ---------- Question related models ---------- //
type Question struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz(id)"`
	ParentID       *uuid.UUID     `json:"parent_id" gorm:"column:parent_id;type:uuid;references:question(id)"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectedUpTo   int            `json:"selected_up_to" gorm:"column:selected_up_to;type:int"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (Question) TableName() string {
	return "question"
}

type QuestionHistory struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID     uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	QuizHistoryID  uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz_history(id)"`
	ParentID       *uuid.UUID     `json:"parent_id" gorm:"column:parent_id;type:uuid;references:question_history(id)"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectedUpTo   int            `json:"selected_up_to" gorm:"column:selected_up_to;type:int"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted        bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy      uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (QuestionHistory) TableName() string {
	return "question_history"
}

// ---------- Options related models ---------- //
// Choice related models
type ChoiceOption struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	Order      int            `json:"order" gorm:"column:order;type:int"`
	Content    string         `json:"content" gorm:"column:content;type:text"`
	Mark       int            `json:"mark" gorm:"column:mark;type:int"`
	Color      string         `json:"color" gorm:"column:color;type:text"`
	Correct    bool           `json:"correct" gorm:"column:correct;type:boolean"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (ChoiceOption) TableName() string {
	return "option_choice"
}

type ChoiceOptionHistory struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	ChoiceOptionID uuid.UUID      `json:"option_choice_id" gorm:"column:option_choice_id;type:uuid;not null;references:option_choice(id)"`
	QuestionID     uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Mark           int            `json:"mark" gorm:"column:mark;type:int"`
	Color          string         `json:"color" gorm:"column:color;type:text"`
	Correct        bool           `json:"correct" gorm:"column:correct;type:boolean"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted        bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy      uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (ChoiceOptionHistory) TableName() string {
	return "option_choice_history"
}

// Text related models
type TextOption struct {
	ID            uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID    uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	Order         int            `json:"order" gorm:"column:order;type:int"`
	Content       string         `json:"content" gorm:"column:content;type:text"`
	Mark          int            `json:"mark" gorm:"column:mark;type:int"`
	CaseSensitive bool           `json:"case_sensitive" gorm:"column:case_sensitive;type:boolean"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (TextOption) TableName() string {
	return "option_text"
}

type TextOptionHistory struct {
	ID            uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionTextID  uuid.UUID      `json:"option_text_id" gorm:"column:option_text_id;type:uuid;not null;references:option_text(id)"`
	QuestionID    uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	Order         int            `json:"order" gorm:"column:order;type:int"`
	Content       string         `json:"content" gorm:"column:content;type:text"`
	Mark          int            `json:"mark" gorm:"column:mark;type:int"`
	CaseSensitive bool           `json:"case_sensitive" gorm:"column:case_sensitive;type:boolean"`
	CreatedAt     time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted       bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy     uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (TextOptionHistory) TableName() string {
	return "option_text_history"
}

// Matching related models
type MatchingOption struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	PromptID   uuid.UUID      `json:"prompt_id" gorm:"column:prompt_id;type:uuid"`
	OptionID   uuid.UUID      `json:"option_id" gorm:"column:option_id;type:uuid"`
	Mark       int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

type MatchingOptionPrompt struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingID uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching(id)"`
	Content          string         `json:"content" gorm:"column:content;type:text"`
	Order            int            `json:"order" gorm:"column:order;type:int"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

type MatchingOptionOption struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingID uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching(id)"`
	Content          string         `json:"content" gorm:"column:content;type:text"`
	Order            int            `json:"order" gorm:"column:order;type:int"`
	Eliminated       bool           `json:"eliminated" gorm:"column:eliminated;type:bool"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (MatchingOption) TableName() string {
	return "option_matching"
}
func (MatchingOptionPrompt) TableName() string {
	return "option_matching_prompt"
}
func (MatchingOptionOption) TableName() string {
	return "option_matching_option"
}

type MatchingOptionHistory struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingID uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching(id)"`
	QuestionID       uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	PromptID         uuid.UUID      `json:"prompt_id" gorm:"column:prompt_id;type:uuid"`
	OptionID         uuid.UUID      `json:"option_id" gorm:"column:option_id;type:uuid"`
	Mark             int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted          bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy        uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

type MatchingOptionPromptHistory struct {
	ID                     uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingPromptID uuid.UUID      `json:"option_matching_prompt_id" gorm:"column:option_matching_prompt_id;type:uuid;not null;references:option_matching_prompt(id)"`
	OptionMatchingID       uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching_history(id)"`
	Content                string         `json:"content" gorm:"column:content;type:text"`
	Order                  int            `json:"order" gorm:"column:order;type:int"`
	CreatedAt              time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted                bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt              gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy              uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

type MatchingOptionOptionHistory struct {
	ID                     uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingOptionID uuid.UUID      `json:"option_matching_option_id" gorm:"column:option_matching_option_id;type:uuid;not null;references:option_matching_option(id)"`
	OptionMatchingID       uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching_history(id)"`
	Content                string         `json:"content" gorm:"column:content;type:text"`
	Order                  int            `json:"order" gorm:"column:order;type:int"`
	Eliminated             bool           `json:"eliminated" gorm:"column:eliminated;type:bool"`
	CreatedAt              time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted                bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt              gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy              uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (MatchingOptionHistory) TableName() string {
	return "option_matching_history"
}
func (MatchingOptionPromptHistory) TableName() string {
	return "option_matching_prompt_history"
}
func (MatchingOptionOptionHistory) TableName() string {
	return "option_matching_option_history"
}

// Pin related models
type PinOption struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	XAxis      int            `json:"x_axis" gorm:"column:x_axis;type:int"`
	YAxis      int            `json:"y_axis" gorm:"column:y_axis;type:int"`
	Mark       int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (PinOption) TableName() string {
	return "option_pin"
}

type PinOptionHistory struct {
	ID          uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionPinID uuid.UUID      `json:"option_pin_id" gorm:"column:option_pin_id;type:uuid;not null;references:option_pin(id)"`
	QuestionID  uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	XAxis       int            `json:"x_axis" gorm:"column:x_axis;type:int"`
	YAxis       int            `json:"y_axis" gorm:"column:y_axis;type:int"`
	Mark        int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt   time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	Deleted     bool           `json:"deleted" gorm:"column:deleted;type:boolean;default:false"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
	UpdatedBy   uuid.UUID      `json:"updated_by" gorm:"column:updated_by;type:uuid;not null;references:user(id)"`
}

func (PinOptionHistory) TableName() string {
	return "option_pin_history"
}

type Repository interface {
	// ---------- Quiz related repository methods ---------- //
	CreateQuiz(ctx context.Context, quiz *Quiz) (*Quiz, error)
	GetQuizzesByUserID(ctx context.Context, uid uuid.UUID) ([]Quiz, error)
	GetQuizByID(ctx context.Context, id uuid.UUID) (*Quiz, error)
	UpdateQuiz(ctx context.Context, quiz *Quiz) (*Quiz, error)
	DeleteQuiz(ctx context.Context, id uuid.UUID) error
	RestoreQuiz(ctx context.Context, id uuid.UUID) (*Quiz, error)
	CreateQuizHistory(ctx context.Context, quizHistory *QuizHistory) (*QuizHistory, error)
	GetQuizHistories(ctx context.Context) ([]QuizHistory, error)
	GetQuizHistoryByID(ctx context.Context, id uuid.UUID) (*QuizHistory, error)
	GetQuizHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuizHistory, error)
	GetQuizHistoryByQuizIDAndCreatedDate(ctx context.Context, quizID uuid.UUID, createdDate time.Time) (*QuizHistory, error)
	UpdateQuizHistory(ctx context.Context, quizHistory *QuizHistory) (*QuizHistory, error)
	DeleteQuizHistory(ctx context.Context, id uuid.UUID) error

	// ---------- Question related repository methods ---------- //
	CreateQuestion(ctx context.Context, question *Question) (*Question, error)
	GetQuestions(ctx context.Context) ([]Question, error)
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error)
	GetQuestionByQuizID(ctx context.Context, quizID uuid.UUID, order int) (*Question, error)
	UpdateQuestion(ctx context.Context, question *Question) (*Question, error)
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
	RestoreQuestion(ctx context.Context, id uuid.UUID) (*Question, error)
	CreateQuestionHistory(ctx context.Context, questionHistory *QuestionHistory) (*QuestionHistory, error)
	GetQuestionHistories(ctx context.Context) ([]QuestionHistory, error)
	GetQuestionHistoryByID(ctx context.Context, id uuid.UUID) (*QuestionHistory, error)
	GetQuestionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]QuestionHistory, error)
	GetQuestionHistoryByQuestionIDAndCreatedDate(ctx context.Context, questionID uuid.UUID, createdDate time.Time) (*QuestionHistory, error)
	UpdateQuestionHistory(ctx context.Context, questionHistory *QuestionHistory) (*QuestionHistory, error)
	DeleteQuestionHistory(ctx context.Context, id uuid.UUID) error

	// ---------- Options related repository methods ---------- //
	// Choice related repository methods
	CreateChoiceOption(ctx context.Context, optionChoice *ChoiceOption) (*ChoiceOption, error)
	GetChoiceOptionByID(ctx context.Context, id uuid.UUID) (*ChoiceOption, error)
	GetChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error)
	UpdateChoiceOption(ctx context.Context, optionChoice *ChoiceOption) (*ChoiceOption, error)
	DeleteChoiceOption(ctx context.Context, id uuid.UUID) error
	RestoreChoiceOption(ctx context.Context, id uuid.UUID) (*ChoiceOption, error)
	CreateChoiceOptionHistory(ctx context.Context, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error)
	GetChoiceOptionHistoryByID(ctx context.Context, id uuid.UUID) (*ChoiceOptionHistory, error)
	GetChoiceOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionHistory, error)
	UpdateChoiceOptionHistory(ctx context.Context, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error)
	DeleteChoiceOptionHistory(ctx context.Context, id uuid.UUID) error

	// Text related repository methods
	CreateTextOption(ctx context.Context, optionText *TextOption) (*TextOption, error)
	GetTextOptionByID(ctx context.Context, id uuid.UUID) (*TextOption, error)
	GetTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error)
	UpdateTextOption(ctx context.Context, optionText *TextOption) (*TextOption, error)
	DeleteTextOption(ctx context.Context, id uuid.UUID) error
	RestoreTextOption(ctx context.Context, id uuid.UUID) (*TextOption, error)
	CreateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	GetTextOptionHistoryByID(ctx context.Context, id uuid.UUID) (*TextOptionHistory, error)
	GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistory, error)
	UpdateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	DeleteTextOptionHistory(ctx context.Context, id uuid.UUID) error

	// Matching related repository methods
	// Pin related repository methods
}

// ---------- Quiz related structs ---------- //
type QuizResponse struct {
	ID          uuid.UUID          `json:"id"`
	CreatorID   uuid.UUID          `json:"creator_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	CoverImage  string             `json:"cover_image"`
	Questions   []QuestionResponse `json:"questions"`
}
type CreateQuizRequest struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Questions   []CreateQuestionRequest `json:"questions"`
}
type CreateQuizResponse struct {
	QuizResponse
	QuizHistoryID uuid.UUID `json:"quiz_history_id"`
}
type UpdateQuizRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverImage  string `json:"cover_image"`
}

// ---------- Question related structs ---------- //
type QuestionResponse struct {
	ID             uuid.UUID          `json:"id"`
	QuizID         uuid.UUID          `json:"quiz_id"`
	ParentID       *uuid.UUID         `json:"parent_id"`
	Type           string             `json:"type"`
	Order          int                `json:"order"`
	Content        string             `json:"content"`
	Note           string             `json:"note"`
	Media          string             `json:"media"`
	TimeLimit      int                `json:"time_limit"`
	HaveTimeFactor bool               `json:"have_time_factor"`
	TimeFactor     int                `json:"time_factor"`
	FontSize       int                `json:"font_size"`
	LayoutIdx      int                `json:"layout_idx"`
	SelectedUpTo   int                `json:"selected_up_to"`
	SubQuestions   []QuestionResponse `json:"sub_questions"`
	Options        []any              `json:"options"`
}
type CreateQuestionRequest struct {
	Type           string `json:"type"`
	Order          int    `json:"order"`
	Content        string `json:"content"`
	Note           string `json:"note"`
	Media          string `json:"media"`
	TimeLimit      int    `json:"time_limit"`
	HaveTimeFactor bool   `json:"have_time_factor"`
	TimeFactor     int    `json:"time_factor"`
	FontSize       int    `json:"font_size"`
	LayoutIdx      int    `json:"layout_idx"`
	SelectedUpTo   int    `json:"selected_up_to"`
	Options        []any  `json:"options"`
}
type CreateQuestionResponse struct {
	QuestionResponse
	QuestionHistoryID uuid.UUID `json:"question_history_id"`
}
type UpdateQuestionRequest struct {
	Type           string `json:"type"`
	Order          int    `json:"order"`
	Content        string `json:"content"`
	Note           string `json:"note"`
	Media          string `json:"media"`
	TimeLimit      int    `json:"time_limit"`
	HaveTimeFactor bool   `json:"have_time_factor"`
	TimeFactor     int    `json:"time_factor"`
	FontSize       int    `json:"font_size"`
	LayoutIdx      int    `json:"layout_idx"`
	SelectedUpTo   int    `json:"selected_up_to"`
}

// ---------- Options related structs ---------- //
// Choice related structs
type ChoiceOptioneResponse struct {
	ID         uuid.UUID `json:"id"`
	QuestionID uuid.UUID `json:"question_id"`
	Order      int       `json:"order"`
	Content    string    `json:"content"`
	Mark       int       `json:"mark"`
	Color      string    `json:"color"`
	Correct    bool      `json:"correct"`
}
type CreateChoiceOptionRequest struct {
	Order   int    `json:"order"`
	Content string `json:"content"`
	Mark    int    `json:"mark"`
	Color   string `json:"color"`
	Correct bool   `json:"correct"`
}
type UpdateChoiceOptionRequest struct {
	Order   int    `json:"order"`
	Content string `json:"content"`
	Mark    int    `json:"mark"`
	Color   string `json:"color"`
	Correct bool   `json:"correct"`
}

// Text related structs
type TextOptionResponse struct {
	ID            uuid.UUID `json:"id"`
	QuestionID    uuid.UUID `json:"question_id"`
	Order         int       `json:"order"`
	Content       string    `json:"content"`
	Mark          int       `json:"mark"`
	CaseSensitive bool      `json:"case_sensitive"`
}
type CreateTextOptionRequest struct {
	Order         int    `json:"order"`
	Content       string `json:"content"`
	Mark          int    `json:"mark"`
	CaseSensitive bool   `json:"case_sensitive"`
}
type UpdateTextOptionRequest struct {
	Order         int    `json:"order"`
	Content       string `json:"content"`
	Mark          int    `json:"mark"`
	CaseSensitive bool   `json:"case_sensitive"`
}

type Service interface {
	// ---------- Quiz related service methods ---------- //
	CreateQuiz(ctx context.Context, req *CreateQuizRequest, uid uuid.UUID) (*CreateQuizResponse, error)
	GetQuizzes(ctx context.Context, uid uuid.UUID) ([]QuizResponse, error)
	GetQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	UpdateQuiz(ctx context.Context, req *UpdateQuizRequest, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	DeleteQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)

	// ---------- Question related service methods ---------- //
	CreateQuestion(ctx context.Context, req *CreateQuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error)
	GetQuestionsByQuizID(ctx context.Context, id uuid.UUID) ([]QuestionResponse, error)
	UpdateQuestion(ctx context.Context, req *UpdateQuestionRequest, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error)
	DeleteQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error)

	// ---------- Options related service methods ---------- //
	// Choice related service methods
	CreateChoiceOption(ctx context.Context, req *CreateChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)
	GetChoiceOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptioneResponse, error)
	UpdateChoiceOption(ctx context.Context, req *UpdateChoiceOptionRequest, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)
	DeleteChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)

	// Text related service methods
	CreateTextOption(ctx context.Context, req *CreateTextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)
	GetTextOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error)
	UpdateTextOption(ctx context.Context, req *UpdateTextOptionRequest, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)
	DeleteTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)

	// Matching related service methods
	// Pin related service methods
}
