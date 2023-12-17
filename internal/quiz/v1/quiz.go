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
	SelectUpTo     int            `json:"select_up_to" gorm:"column:select_up_to;type:int"`
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
	SelectUpTo     int            `json:"select_up_to" gorm:"column:select_up_to;type:int"`
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
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
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
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
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
	QuestionPoolID *uuid.UUID      `json:"question_pool_id,omitempty" gorm:"column:question_pool_id;type:uuid;references:question_pool(id)"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	UseTemplate    bool           `json:"use_template" gorm:"column:use_template;type:boolean"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectUpTo     int            `json:"select_up_to" gorm:"column:select_up_to;type:int"`
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
	QuestionPoolID *uuid.UUID      `json:"question_pool_id,omitempty" gorm:"column:question_pool_id;type:uuid;references:question_pool_history(id)"`
	Type           string         `json:"type" gorm:"column:type;type:text"`
	Order          int            `json:"order" gorm:"column:order;type:int"`
	Content        string         `json:"content" gorm:"column:content;type:text"`
	Note           string         `json:"note" gorm:"column:note;type:text"`
	Media          string         `json:"media" gorm:"column:media;type:text"`
	UseTemplate    bool           `json:"use_template" gorm:"column:use_template;type:boolean"`
	TimeLimit      int            `json:"time_limit" gorm:"column:time_limit;type:int"`
	HaveTimeFactor bool           `json:"have_time_factor" gorm:"column:have_time_factor;type:boolean"`
	TimeFactor     int            `json:"time_factor" gorm:"column:time_factor;type:int"`
	FontSize       int            `json:"font_size" gorm:"column:font_size;type:int"`
	LayoutIdx      int            `json:"layout_idx" gorm:"column:layout_idx;type:int"`
	SelectUpTo     int            `json:"select_up_to" gorm:"column:select_up_to;type:int"`
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
	Content    string         `json:"content" gorm:"column:content;type:text"`
	Order      int            `json:"order" gorm:"column:order;type:int"`
	Type       string         `json:"type" gorm:"column:type;type:text"`
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
	OptionMatchingID uuid.UUID      `json:"option_matching_id" gorm:"column:option_matching_id;type:uuid;not null;references:option_matching(id)"`
	QuestionID       uuid.UUID      `json:"question_id" gorm:"column:question_id;type:uuid;not null;references:question_history(id)"`
	Content          string         `json:"content" gorm:"column:content;type:text"`
	Order            int            `json:"order" gorm:"column:order;type:int"`
	Type             string         `json:"type" gorm:"column:type;type:text"`
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
	AnswerMatchingID uuid.UUID      `json:"answer_matching_id" gorm:"column:answer_matching_id;type:uuid;not null;references:answer_matching(id)"`
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

	// ---------- Question Pool related repository methods ---------- //
	CreateQuestionPool(ctx context.Context, questionPool *QuestionPool) (*QuestionPool, error)
	UpdateQuestionPool(ctx context.Context, questionPool *QuestionPool) (*QuestionPool, error)
	DeleteQuestionPool(ctx context.Context, id uuid.UUID) error
	RestoreQuestionPool(ctx context.Context, id uuid.UUID) (*QuestionPool, error)
	CreateQuestionPoolHistory(ctx context.Context, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error)
	UpdateQuestionPoolHistory(ctx context.Context, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error)
	DeleteQuestionPoolHistory(ctx context.Context, id uuid.UUID) error

	// ---------- Question related repository methods ---------- //
	CreateQuestion(ctx context.Context, question *Question) (*Question, error)
	GetQuestions(ctx context.Context) ([]Question, error)
	GetQuestionByID(ctx context.Context, id uuid.UUID) (*Question, error)
	GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error)
	GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error)
	GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
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
	GetChoiceAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error)
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
	GetTextAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error)
	UpdateTextOption(ctx context.Context, optionText *TextOption) (*TextOption, error)
	DeleteTextOption(ctx context.Context, id uuid.UUID) error
	RestoreTextOption(ctx context.Context, id uuid.UUID) (*TextOption, error)
	CreateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	GetTextOptionHistoryByID(ctx context.Context, id uuid.UUID) (*TextOptionHistory, error)
	GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistory, error)
	UpdateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error)
	DeleteTextOptionHistory(ctx context.Context, id uuid.UUID) error

	// Option Matching related repository methods
	CreateMatchingOption(ctx context.Context, optionMatching *MatchingOption) (*MatchingOption, error)
	GetMachingOptionByID(ctx context.Context, id uuid.UUID) (*MatchingOption, error)
	GetMatchingOptionByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOption, error)
	UpdateMatchingOption(ctx context.Context, optionMatching *MatchingOption) (*MatchingOption, error)
	DeleteMatchingOption(ctx context.Context, id uuid.UUID) error
	RestoreMatchingOption(ctx context.Context, id uuid.UUID) (*MatchingOption, error)
	CreateMatchingOptionHistory(ctx context.Context, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error)
	GetMatchingOptionHistoryByID(ctx context.Context, id uuid.UUID) (*MatchingOptionHistory, error)
	GetOptionMatchingHistories(ctx context.Context) ([]MatchingOptionHistory, error)
	GetMatchingOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionHistory, error)
	UpdateMatchingOptionHistory(ctx context.Context, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error)
	DeleteMatchingOptionHistory(ctx context.Context, id uuid.UUID) error

	// Answer Matching related repository methods
	CreateMatchingAnswer(ctx context.Context, answerMatching *MatchingAnswer) (*MatchingAnswer, error)
	UpdateMatchingAnswer(ctx context.Context, answerMatching *MatchingAnswer) (*MatchingAnswer, error)
	DeleteMatchingAnswer(ctx context.Context, id uuid.UUID) error
	RestoreMatchingAnswer(ctx context.Context, id uuid.UUID) (*MatchingAnswer, error)
	CreateMatchingAnswerHistory(ctx context.Context, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error)
	UpdateMatchingAnswerHistory(ctx context.Context, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error)
	DeleteMatchingAnswerHistory(ctx context.Context, id uuid.UUID) error
}

// ---------- Quiz related structs ---------- //
type QuizResponse struct {
	Quiz
	Questions []QuestionResponse `json:"questions"`
}

type CreateQuizRequest struct {
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	CoverImage     string                  `json:"cover_image"`
	Visibility     string                  `json:"visiblity"`
	TimeLimit      int                     `json:"time_limit"`
	HaveTimeFactor bool                    `json:"have_time_limit"`
	TimeFactor     int                     `json:"time_factor"`
	FontSize       int                     `json:"font_size"`
	Mark           int                     `json:"mark"`
	SelectUpTo     int                     `json:"select_up_to"`
	CaseSensitive  bool                    `json:"case_sensitive"`
	Questions      []CreateQuestionRequest `json:"questions"`
}

type CreateQuizResponse struct {
	QuizResponse
	QuizHistoryID uuid.UUID `json:"quiz_history_id"`
}

type UpdateQuizRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	CoverImage     string `json:"cover_image"`
	Visibility     string `json:"visiblity"`
	TimeLimit      int    `json:"time_limit"`
	HaveTimeFactor bool   `json:"have_time_limit"`
	TimeFactor     int    `json:"time_factor"`
	FontSize       int    `json:"font_size"`
	Mark           int    `json:"mark"`
	SelectUpTo     int    `json:"select_up_to"`
}

// ---------- Question Pool related structs ---------- //
type QuestionPoolResponse struct {
	QuestionPool
}

type CreateQuestionPoolResponse struct {
	QuestionPoolResponse
	QuestionPoolHistoryID uuid.UUID `json:"question_pool_history_id"`
}

// ---------- Question related structs ---------- //
type QuestionResponse struct {
	Question
	Options []any `json:"options"`
}

type CreateQuestionRequest struct {
	IsInPool			 bool       `json:"is_in_pool"`
	Type           string     `json:"type"`
	Order          int        `json:"order"`
	Content        string     `json:"content"`
	Note           string     `json:"note"`
	Media          string     `json:"media"`
	UseTemplate    bool       `json:"use_template"`
	TimeLimit      int        `json:"time_limit"`
	HaveTimeFactor bool       `json:"have_time_factor"`
	TimeFactor     int        `json:"time_factor"`
	FontSize       int        `json:"font_size"`
	LayoutIdx      int        `json:"layout_idx"`
	SelectUpTo     int        `json:"select_up_to"`
	Options        []any      `json:"options,omitempty"`
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
	UseTemplate    string `json:"use_template"`
	TimeLimit      int    `json:"time_limit"`
	HaveTimeFactor bool   `json:"have_time_factor"`
	TimeFactor     int    `json:"time_factor"`
	FontSize       int    `json:"font_size"`
	LayoutIdx      int    `json:"layout_idx"`
	SelectUpTo     int    `json:"select_up_to"`
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

// Matching related structs

type Service interface {
	// ---------- Quiz related service methods ---------- //
	CreateQuiz(ctx context.Context, req *CreateQuizRequest, uid uuid.UUID) (*CreateQuizResponse, error)
	GetQuizzes(ctx context.Context, uid uuid.UUID) ([]QuizResponse, error)
	GetQuizByID(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	UpdateQuiz(ctx context.Context, req *UpdateQuizRequest, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)
	DeleteQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreQuiz(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuizResponse, error)

	// ---------- Question Pool related service methods ---------- //
	CreateQuestionPool(ctx context.Context, req *CreateQuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID) (*CreateQuestionPoolResponse, error)

	// ---------- Question related service methods ---------- //
	CreateQuestion(ctx context.Context, req *CreateQuestionRequest, quizID uuid.UUID, quizHistoryID uuid.UUID, questionPoolID *uuid.UUID, QuestionPoolHistoryID *uuid.UUID, uid uuid.UUID) (*CreateQuestionResponse, error)
	GetQuestionsByQuizID(ctx context.Context, id uuid.UUID) ([]QuestionResponse, error)
	GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error)
	GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error)
	UpdateQuestion(ctx context.Context, req *UpdateQuestionRequest, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error)
	DeleteQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreQuestion(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*QuestionResponse, error)

	// ---------- Options related service methods ---------- //
	// Choice related service methods
	CreateChoiceOption(ctx context.Context, req *CreateChoiceOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)
	GetChoiceOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptioneResponse, error)
	GetChoiceAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]ChoiceOptioneResponse, error)
	UpdateChoiceOption(ctx context.Context, req *UpdateChoiceOptionRequest, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)
	DeleteChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreChoiceOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*ChoiceOptioneResponse, error)

	// Text related service methods
	CreateTextOption(ctx context.Context, req *CreateTextOptionRequest, questionID uuid.UUID, questionHistoryID uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)
	GetTextOptionsByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error)
	GetTextAnswersByQuestionID(ctx context.Context, id uuid.UUID) ([]TextOptionResponse, error)
	UpdateTextOption(ctx context.Context, req *UpdateTextOptionRequest, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)
	DeleteTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) error
	RestoreTextOption(ctx context.Context, id uuid.UUID, uid uuid.UUID) (*TextOptionResponse, error)

	// Matching related service methods
}
