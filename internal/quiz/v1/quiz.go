package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ---------- Quiz related models ---------- //
type Quiz struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	CreatorID      uuid.UUID      `json:"creator_id" gorm:"column:creator_id;type:uuid;not null;references:user(id)"`
	Title          string         `json:"title" gorm:"column:title;type:text;default:Untitled"`
	Description    string         `json:"description" gorm:"column:description;type:text"`
	CoverImage     string         `json:"cover_image" gorm:"column:cover_image;type:text"`
	Visibility     string         `json:"visibility" gorm:"column:visibility;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	Mark           int            `json:"mark" gorm:"column:mark;type:int"`
	SelectMin      int            `json:"select_min" gorm:"column:select_min;type:int"`
	SelectMax      int            `json:"select_max" gorm:"column:select_max;type:int"`
	CaseSensitive  bool           `json:"case_sensitive" gorm:"column:case_sensitive;type:boolean"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (Quiz) TableName() string {
	return "quiz"
}

type QuizHistory struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz(id)"`
	CreatorID      uuid.UUID      `json:"creator_id" gorm:"column:creator_id;type:uuid;not null;references:user(id)"`
	Title          string         `json:"title" gorm:"column:title;type:text;default:Untitled"`
	Description    string         `json:"description" gorm:"column:description;type:text"`
	CoverImage     string         `json:"cover_image" gorm:"column:cover_image;type:text"`
	Visibility     string         `json:"visibility" gorm:"column:visibility;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	Mark           int            `json:"mark" gorm:"column:mark;type:int"`
	SelectMin      int            `json:"select_min" gorm:"column:select_min;type:int"`
	SelectMax      int            `json:"select_max" gorm:"column:select_max;type:int"`
	CaseSensitive  bool           `json:"case_sensitive" gorm:"column:case_sensitive;type:boolean"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (QuizHistory) TableName() string {
	return "quiz_history"
}

// ---------- Question Pool related models ---------- //
type QuestionPool struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz(id)"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	PoolOrder      int            `json:"pool_order" gorm:"column:pool_order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	MediaType      string         `json:"media_type" gorm:"column:media_type;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (QuestionPool) TableName() string {
	return "question_pool"
}

type QuestionPoolHistory struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionPoolID uuid.UUID      `json:"question_pool_id" gorm:"column:question_pool_id;type:uuid;not null;references:question_pool(id)"`
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz_history(id)"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	PoolOrder      int            `json:"pool_order" gorm:"column:pool_order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	MediaType      string         `json:"media_type" gorm:"column:media_type;type:text"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (QuestionPoolHistory) TableName() string {
	return "question_pool_history"
}

// ---------- Question related models ---------- //
type Question struct {
	ID             uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz(id)"`
	QuestionPoolID *uuid.UUID     `json:"question_pool_id,omitempty" gorm:"column:question_pool_id;type:uuid;references:question_pool(id)"`
	PoolOrder      int            `json:"pool_order" gorm:"column:pool_order;type:int"`
	PoolRequired   bool           `json:"pool_required" gorm:"column:pool_required;type:bool"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	MediaType      string         `json:"media_type" gorm:"column:media_type;type:text"`
	UseTemplate    bool           `json:"use_template" gorm:"column:use_template;type:boolean"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectMin      int            `json:"select_min" gorm:"column:select_min;type:int"`
	SelectMax      int            `json:"select_max" gorm:"column:select_max;type:int"`
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
	QuizID         uuid.UUID      `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null;references:quiz_history(id)"`
	QuestionPoolID *uuid.UUID     `json:"question_pool_id,omitempty" gorm:"column:question_pool_id;type:uuid;references:question_pool_history(id)"`
	PoolOrder      int            `json:"pool_order" gorm:"column:pool_order;type:int"`
	PoolRequired   bool           `json:"pool_required" gorm:"column:pool_required;type:bool"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	MediaType      string         `json:"media_type" gorm:"column:media_type;type:text"`
	UseTemplate    bool           `json:"use_template" gorm:"column:use_template;type:boolean"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectMin      int            `json:"select_min" gorm:"column:select_min;type:int"`
	SelectMax      int            `json:"select_max" gorm:"column:select_max;type:int"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
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
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
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
	UpdatedAt     time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (TextOptionHistory) TableName() string {
	return "option_text_history"
}

// Matching related models
type MatchingOption struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	Type       string         `json:"type" gorm:"column:type;type:text"`
	Order      int            `json:"order" gorm:"column:order;type:int"`
	Content    string         `json:"content" gorm:"column:content;type:text"`
	Color      string         `json:"color" gorm:"column:color;type:text"`
	Eliminate  bool           `json:"eliminate" gorm:"column:eliminate;type:boolean"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (MatchingOption) TableName() string {
	return "option_matching"
}

type MatchingOptionHistory struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingID uuid.UUID      `json:"matching_option_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching(id)"`
	QuestionID       uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	Type             string         `json:"type" gorm:"column:type;type:text"`
	Order            int            `json:"order" gorm:"column:order;type:int"`
	Content          string         `json:"content" gorm:"column:content;type:text"`
	Color            string         `json:"color" gorm:"column:color;type:text"`
	Eliminate        bool           `json:"eliminate" gorm:"column:eliminate;type:boolean"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (MatchingOptionHistory) TableName() string {
	return "option_matching_history"
}

type MatchingAnswer struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	PromptID   uuid.UUID      `json:"prompt_id" gorm:"column:prompt_id;type:uuid"`
	OptionID   uuid.UUID      `json:"option_id" gorm:"column:option_id;type:uuid"`
	Mark       int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (MatchingAnswer) TableName() string {
	return "answer_matching"
}

type MatchingAnswerHistory struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	AnswerMatchingID uuid.UUID      `json:"matching_answer_id" gorm:"column:answer_matching_id;type:uuid;not null;references:answer_matching(id)"`
	QuestionID       uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	PromptID         uuid.UUID      `json:"prompt_id" gorm:"column:prompt_id;type:uuid"`
	OptionID         uuid.UUID      `json:"option_id" gorm:"column:option_id;type:uuid"`
	Mark             int            `json:"mark" gorm:"column:mark;type:int"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

func (MatchingAnswerHistory) TableName() string {
	return "answer_matching_history"
}

type Repository interface {
	// ---------- Transaction repository methods ---------- //
	BeginTransaction() (*gorm.DB, error)
	CommitTransaction(tx *gorm.DB) error

	// ---------- Quiz related repository methods ---------- //
	CreateQuiz(ctx context.Context, tx *gorm.DB, quiz *Quiz) (*Quiz, error)
	GetQuizzesByUserID(ctx context.Context, uid uuid.UUID) ([]Quiz, error)
	GetQuizByID(ctx context.Context, id uuid.UUID) (*Quiz, error)
	UpdateQuiz(ctx context.Context, tx *gorm.DB, quiz *Quiz) (*Quiz, error)
	DeleteQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*Quiz, error)
	CreateQuizHistory(ctx context.Context, tx *gorm.DB, quizHistory *QuizHistory) (*QuizHistory, error)
	GetQuizHistories(ctx context.Context) ([]QuizHistory, error)
	GetQuizHistoryByID(ctx context.Context, id uuid.UUID) (*QuizHistory, error)
	GetQuizHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuizHistory, error)
	GetQuizHistoriesByUserID(ctx context.Context, uid uuid.UUID) ([]QuizHistory, error)
	GetQuizHistoryByQuizIDAndCreatedDate(ctx context.Context, quizID uuid.UUID, createdDate time.Time) (*QuizHistory, error)
	GetDeleteQuizByID(ctx context.Context, id uuid.UUID) (*Quiz, error)
	UpdateQuizHistory(ctx context.Context, tx *gorm.DB, quizHistory *QuizHistory) (*QuizHistory, error)
	DeleteQuizHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// ---------- Question Pool related repository methods ---------- //
	CreateQuestionPool(ctx context.Context, tx *gorm.DB, questionPool *QuestionPool) (*QuestionPool, error)
	GetQuestionPoolByID(ctx context.Context, questionPoolID uuid.UUID) (*QuestionPool, error)
	GetQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPool, error)
	GetDeleteQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPool, error)
	UpdateQuestionPool(ctx context.Context, tx *gorm.DB, questionPool *QuestionPool) (*QuestionPool, error)
	DeleteQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*QuestionPool, error)
	CreateQuestionPoolHistory(ctx context.Context, tx *gorm.DB, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error)
	GetQuestionPoolHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolHistory, error)
	UpdateQuestionPoolHistory(ctx context.Context, tx *gorm.DB, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error)
	DeleteQuestionPoolHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// ---------- Question related repository methods ---------- //
	CreateQuestion(ctx context.Context, tx *gorm.DB, question *Question) (*Question, error)
	GetQuestions(ctx context.Context) ([]Question, error)
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error)
	GetDeleteQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error)
	GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error)
	GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
	UpdateQuestion(ctx context.Context, tx *gorm.DB, question *Question) (*Question, error)
	DeleteQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*Question, error)
	CreateQuestionHistory(ctx context.Context, tx *gorm.DB, questionHistory *QuestionHistory) (*QuestionHistory, error)
	GetQuestionHistories(ctx context.Context) ([]QuestionHistory, error)
	GetQuestionHistoryByID(ctx context.Context, id uuid.UUID) (*QuestionHistory, error)
	GetQuestionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]QuestionHistory, error)
	GetQuestionHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionHistory, error)
	GetQuestionHistoryByQuestionIDAndCreatedDate(ctx context.Context, questionID uuid.UUID, createdDate time.Time) (*QuestionHistory, error)
	UpdateQuestionHistory(ctx context.Context, tx *gorm.DB, questionHistory *QuestionHistory) (*QuestionHistory, error)
	DeleteQuestionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// ---------- Options related repository methods ---------- //
	// Choice related repository methods
	CreateChoiceOption(ctx context.Context, tx *gorm.DB, optionChoice *ChoiceOption) (*ChoiceOption, error)
	GetChoiceOptionByID(ctx context.Context, id uuid.UUID) (*ChoiceOption, error)
	GetChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error)
	GetDeleteChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error)
	GetChoiceAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error)
	UpdateChoiceOption(ctx context.Context, tx *gorm.DB, optionChoice *ChoiceOption) (*ChoiceOption, error)
	DeleteChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*ChoiceOption, error)
	CreateChoiceOptionHistory(ctx context.Context, tx *gorm.DB, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error)
	GetChoiceOptionHistoryByID(ctx context.Context, id uuid.UUID) (*ChoiceOptionHistory, error)
	GetChoiceOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionHistory, error)
	UpdateChoiceOptionHistory(ctx context.Context, tx *gorm.DB, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error)
	DeleteChoiceOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// Text related repository methods
	CreateTextOption(ctx context.Context, tx *gorm.DB, optionText *TextOption) (*TextOption, error)
	GetTextOptionByID(ctx context.Context, id uuid.UUID) (*TextOption, error)
	GetTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error)
	GetDeleteTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error)
	GetTextAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error)
	UpdateTextOption(ctx context.Context, tx *gorm.DB, optionText *TextOption) (*TextOption, error)
	DeleteTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*TextOption, error)
	CreateTextOptionHistory(ctx context.Context, tx *gorm.DB, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	GetTextOptionHistoryByID(ctx context.Context, id uuid.UUID) (*TextOptionHistory, error)
	GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistory, error)
	UpdateTextOptionHistory(ctx context.Context, tx *gorm.DB, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	DeleteTextOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// Option Matching related repository methods
	CreateMatchingOption(ctx context.Context, tx *gorm.DB, optionMatching *MatchingOption) (*MatchingOption, error)
	GetMatchingOptionByID(ctx context.Context, id uuid.UUID) (*MatchingOption, error)
	GetMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOption, error)
	GetDeleteMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOption, error)
	GetMatchingOptionByQuestionIDAndOrder(ctx context.Context, questionID uuid.UUID, order int) (*MatchingOption, error)
	UpdateMatchingOption(ctx context.Context, tx *gorm.DB, optionMatching *MatchingOption) (*MatchingOption, error)
	DeleteMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*MatchingOption, error)
	CreateMatchingOptionHistory(ctx context.Context, tx *gorm.DB, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error)
	GetMatchingOptionHistoryByID(ctx context.Context, id uuid.UUID) (*MatchingOptionHistory, error)
	GetOptionMatchingHistories(ctx context.Context) ([]MatchingOptionHistory, error)
	GetMatchingOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionHistory, error)
	UpdateMatchingOptionHistory(ctx context.Context, tx *gorm.DB, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error)
	DeleteMatchingOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	// Answer Matching related repository methods
	CreateMatchingAnswer(ctx context.Context, tx *gorm.DB, answerMatching *MatchingAnswer) (*MatchingAnswer, error)
	UpdateMatchingAnswer(ctx context.Context, tx *gorm.DB, answerMatching *MatchingAnswer) (*MatchingAnswer, error)
	DeleteMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	RestoreMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*MatchingAnswer, error)
	GetMatchingAnswerByID(ctx context.Context, id uuid.UUID) (*MatchingAnswer, error)
	GetMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswer, error)
	GetDeleteMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswer, error)
	CreateMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error)
	GetMatchingAnswerHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerHistory, error)
	UpdateMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error)
	DeleteMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
}

// ---------- Quiz related structs ---------- //
type QuizResponse struct {
	Quiz
	CreatorName string             `json:"creator_name"`
	Questions   []QuestionResponse `json:"questions"`
}

type CreateQuizRequest struct {
	Quiz
	CreatorName string            `json:"creator_name"`
	Questions   []QuestionRequest `json:"questions"`
}

type CreateQuizResponse struct {
	QuizResponse
	QuizHistoryID uuid.UUID `json:"quiz_history_id"`
}

type UpdateQuizResponse struct {
	QuizResponse
	QuizHistoryID uuid.UUID `json:"quiz_history_id"`
}

type UpdateQuizRequest struct {
	Quiz
	Questions []QuestionRequest `json:"questions"`
}

type QuizHistoryResponse struct {
	QuizHistory
	CreatorName     string                    `json:"creator_name"`
	QuestionHistory []QuestionHistoryResponse `json:"questions"`
}

// ---------- Question Pool related structs ---------- //
type QuestionPoolResponse struct {
	QuestionPool
}

type CreateQuestionPoolResponse struct {
	QuestionPoolResponse
	QuestionPoolHistoryID uuid.UUID `json:"question_pool_history_id"`
}

type UpdateQuestionPoolResponse struct {
	QuestionPoolResponse
	QuestionPoolHistoryID uuid.UUID `json:"question_pool_history_id"`
}

type QuestionPoolHistoryResponse struct {
	QuestionPoolHistory
}

// ---------- Question related structs ---------- //
type QuestionResponse struct {
	Question
	Options []any `json:"options,omitempty"`
}

type QuestionRequest struct {
	IsInPool bool `json:"is_in_pool"`
	Question
	Options []any `json:"options,omitempty"`
}

type CreateQuestionResponse struct {
	QuestionResponse
	QuestionHistoryID uuid.UUID `json:"question_history_id"`
}

type UpdateQuestionResponse struct {
	QuestionResponse
	QuestionHistoryID uuid.UUID `json:"question_history_id"`
}

type QuestionHistoryResponse struct {
	QuestionHistory
	Options []any `json:"options,omitempty"`
}

// ---------- Options related structs ---------- //
// Choice related structs
type ChoiceOptionResponse struct {
	ChoiceOption
}

type UpdateChoiceOptionResponse struct {
	ChoiceOption
}

type CreateChoiceOptionResponse struct {
	ChoiceOption
}

type ChoiceOptionRequest struct {
	ChoiceOption
}

type ChoiceOptionHistoryResponse struct {
	ChoiceOptionHistory
}

// Text related structs
type TextOptionResponse struct {
	TextOption
}

type TextOptionRequest struct {
	TextOption
}

type UpdateTextOptionResponse struct {
	TextOption
}

type CreateTextOptionResponse struct {
	TextOption
}

type TextOptionHistoryResponse struct {
	TextOptionHistory
}

// Matching related structs

type MatchingOptionAndAnswerResponse struct {
	ID         uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	QuestionID uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question(id)"`
	Type       string         `json:"type,omitempty" gorm:"column:type;type:text"`
	Order      *int           `json:"order,omitempty" gorm:"column:order;type:int"`
	Content    *string        `json:"content,omitempty" gorm:"column:content;type:text"`
	Color      *string        `json:"color,omitempty" gorm:"column:color;type:text"`
	Eliminate  bool           `json:"eliminate" gorm:"column:eliminate;type:boolean"`
	PromptID   *uuid.UUID     `json:"prompt_id,omitempty" gorm:"column:prompt_id;type:uuid"`
	OptionID   *uuid.UUID     `json:"option_id,omitempty" gorm:"column:option_id;type:uuid"`
	Mark       *int           `json:"mark,omitemtpy" gorm:"column:mark;type:int"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

type MatchingOptionAndAnswerHistoryResponse struct {
	ID               uuid.UUID      `json:"id" gorm:"column:id;type:uuid;primaryKey;not null"`
	OptionMatchingID *uuid.UUID     `json:"option_matching_id,omitempty"`
	AmswerMatchingID *uuid.UUID     `json:"answer_matching_id,omitempty"`
	QuestionID       uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	Type             string         `json:"type,omitempty" gorm:"column:type;type:text"`
	Order            *int           `json:"order,omitempty" gorm:"column:order;type:int"`
	Content          *string        `json:"content,omitempty" gorm:"column:content;type:text"`
	Color            *string        `json:"color,omitempty" gorm:"column:color;type:text"`
	Eliminate        bool           `json:"eliminate" gorm:"column:eliminate;type:boolean"`
	PromptID         *uuid.UUID     `json:"prompt_id,omitempty" gorm:"column:prompt_id;type:uuid"`
	OptionID         *uuid.UUID     `json:"option_id,omitempty" gorm:"column:option_id;type:uuid"`
	Mark             *int           `json:"mark,omitemtpy" gorm:"column:mark;type:int"`
	CreatedAt        time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt        time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp;not null"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:timestamp"`
}

// ------ Matching Option ------
type MatchingOptionResponse struct {
	MatchingOption
}

type MatchingOptionRequest struct {
	MatchingOption
}

type UpdateMatchingOptionResponse struct {
	MatchingOption
}

type CreateMatchingOptionResponse struct {
	MatchingOption
}

type MatchingOptionHistoryResponse struct {
	MatchingOptionHistory
}

// ------ Matching Answer -------
type MatchingAnswerResponse struct {
	MatchingAnswer
}

type MatchingAnswerRequest struct {
	Prompt int `json:"prompt_id"`
	Option int `json:"option_id"`
	MatchingAnswer
}

type UpdateMatchingAnswerResponse struct {
	MatchingAnswer
}

type CreateMatchingAnswerResponse struct {
	MatchingAnswer
}

type MatchingAnswerHistoryResponse struct {
	MatchingAnswerHistory
}

type Service interface {
	// ---------- Transaction related service methods ---------- //
	BeginTransaction(ctx context.Context) (*gorm.DB, error)
	CommitTransaction(ctx context.Context, tx *gorm.DB) error

	// ---------- Quiz related service methods ---------- //
	CreateQuiz(ctx context.Context, tx *gorm.DB, req *CreateQuizRequest, uid uuid.UUID) (*CreateQuizResponse, error)
	GetQuizzes(ctx context.Context, uid uuid.UUID) ([]QuizResponse, error)
	GetQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	GetDeleteQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	UpdateQuiz(ctx context.Context, tx *gorm.DB, req *UpdateQuizRequest, id uuid.UUID, uid uuid.UUID) (*UpdateQuizResponse, error)
	DeleteQuiz(ctx context.Context, tx *gorm.DB, quizID uuid.UUID) error
	RestoreQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetQuizHistories(ctx context.Context, uid uuid.UUID) ([]QuizHistoryResponse, error)
	GetQuizHistoryByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizHistoryResponse, error)

	// ---------- Question Pool related service methods ---------- //
	CreateQuestionPool(ctx context.Context, tx *gorm.DB, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID) (*CreateQuestionPoolResponse, error)
	GetQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolResponse, error)
	GetDeleteQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolResponse, error)
	UpdateQuestionPool(ctx context.Context, tx *gorm.DB, req *QuestionRequest, user_id uuid.UUID, questionPoolID uuid.UUID, quizHistoryID uuid.UUID) (*UpdateQuestionPoolResponse, error)
	DeleteQuestionPool(ctx context.Context, tx *gorm.DB, questionPoolID uuid.UUID) error
	RestoreQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetQuestionPoolHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolHistoryResponse, error)

	// ---------- Question related service methods ---------- //
	CreateQuestion(ctx context.Context, tx *gorm.DB, req *QuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, questionPoolID *uuid.UUID, QuestionPoolHistoryID *uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error)
	GetQuestionsByQuizID(ctx context.Context, id uuid.UUID) ([]QuestionResponse, error)
	GetDeleteQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionResponse, error)
	GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error)
	GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
	UpdateQuestion(ctx context.Context, tx *gorm.DB, req *QuestionRequest, user_id uuid.UUID, questionID uuid.UUID, quizHistoryID uuid.UUID, questionPoolHistoryID *uuid.UUID) (*UpdateQuestionResponse, error)
	DeleteQuestion(ctx context.Context, tx *gorm.DB, questionID uuid.UUID) error
	RestoreQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetQuestionHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionHistoryResponse, error)

	// ---------- Options related service methods ---------- //
	// Choice related service methods
	CreateChoiceOption(ctx context.Context, tx *gorm.DB, req *ChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateChoiceOptionResponse, error)
	GetChoiceOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptionResponse, error)
	GetDeleteChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionResponse, error)
	GetChoiceAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptionResponse, error)
	UpdateChoiceOption(ctx context.Context, tx *gorm.DB, req *ChoiceOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateChoiceOptionResponse, error)
	DeleteChoiceOption(ctx context.Context, tx *gorm.DB, choiceOptionID uuid.UUID) error
	RestoreChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetChoiceOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionHistoryResponse, error)

	// Text related service methods
	CreateTextOption(ctx context.Context, tx *gorm.DB, req *TextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateTextOptionResponse, error)
	GetTextOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error)
	GetDeleteTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionResponse, error)
	GetTextAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error)
	UpdateTextOption(ctx context.Context, tx *gorm.DB, req *TextOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateTextOptionResponse, error)
	DeleteTextOption(ctx context.Context, tx *gorm.DB, textOptionID uuid.UUID) error
	RestoreTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistoryResponse, error)

	// Matching related service methods
	// ----- Matching Option ------
	CreateMatchingOption(ctx context.Context, tx *gorm.DB, req *MatchingOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateMatchingOptionResponse, error)
	GetMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionResponse, error)
	GetDeleteMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionResponse, error)
	GetMatchingOptionByQuestionIDAndOrder(ctx context.Context, questionID uuid.UUID, order int) (*MatchingOptionResponse, error)
	UpdateMatchingOption(ctx context.Context, tx *gorm.DB, req *MatchingOptionRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateMatchingOptionResponse, error)
	DeleteMatchingOption(ctx context.Context, tx *gorm.DB, matchingOptionID uuid.UUID) error
	RestoreMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetMatchingOptionHistoryByID(ctx context.Context, id uuid.UUID) (*MatchingOptionHistoryResponse, error)
	GetMatchingOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionHistoryResponse, error)

	// ----- Matching Answer ------
	CreateMatchingAnswer(ctx context.Context, tx *gorm.DB, req *MatchingAnswerRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*CreateMatchingAnswerResponse, error)
	GetMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerResponse, error)
	GetDeleteMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerResponse, error)
	UpdateMatchingAnswer(ctx context.Context, tx *gorm.DB, req *MatchingAnswerRequest, userID uuid.UUID, optionID uuid.UUID, questionHistoryID uuid.UUID) (*UpdateMatchingAnswerResponse, error)
	DeleteMatchingAnswer(ctx context.Context, tx *gorm.DB, matchingAnswerID uuid.UUID) error
	RestoreMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) error

	GetMatchingAnswerHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerHistoryResponse, error)
}
