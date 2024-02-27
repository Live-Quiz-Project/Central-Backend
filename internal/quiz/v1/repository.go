package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) BeginTransaction() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *repository) CommitTransaction(tx *gorm.DB) error {
	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// ---------- Quiz related repository methods ---------- //
func (r *repository) CreateQuiz(ctx context.Context, tx *gorm.DB, quiz *Quiz) (*Quiz, error) {
	res := tx.WithContext(ctx).Create(quiz)
	if res.Error != nil {
		tx.Rollback()
		return &Quiz{}, res.Error
	}

	return quiz, nil
}

func (r *repository) GetQuizzesByUserID(ctx context.Context, uid uuid.UUID) ([]Quiz, error) {
	var quizzes []Quiz
	res := r.db.WithContext(ctx).Where("creator_id = ?", uid).Find(&quizzes)
	if res.Error != nil {
		return []Quiz{}, res.Error
	}

	return quizzes, nil
}

func (r *repository) GetQuizByID(ctx context.Context, id uuid.UUID) (*Quiz, error) {
	var quiz Quiz
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&quiz)
	if res.Error != nil {
		return &Quiz{}, res.Error
	}

	return &quiz, nil
}

func (r *repository) GetDeleteQuizByID(ctx context.Context, id uuid.UUID) (*Quiz, error) {
	var quiz Quiz
	res := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&quiz)
	if res.Error != nil {
		return &Quiz{}, res.Error
	}
	return &quiz, nil
}

func (r *repository) UpdateQuiz(ctx context.Context, tx *gorm.DB, quiz *Quiz) (*Quiz, error) {
	res := tx.WithContext(ctx).Save(quiz)
	if res.Error != nil {
		tx.Rollback()
		return &Quiz{}, res.Error
	}

	return quiz, nil
}

func (r *repository) DeleteQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&Quiz{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

func (r *repository) RestoreQuiz(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*Quiz, error) {
	var quiz Quiz
	res := r.db.WithContext(ctx).Unscoped().First(&quiz, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&quiz).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	return &quiz, nil
}

func (r *repository) CreateQuizHistory(ctx context.Context, tx *gorm.DB, quizHistory *QuizHistory) (*QuizHistory, error) {
	res := tx.WithContext(ctx).Create(quizHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuizHistory{}, res.Error
	}

	return quizHistory, nil
}

func (r *repository) GetQuizHistories(ctx context.Context) ([]QuizHistory, error) {
	var quizHistories []QuizHistory
	res := r.db.WithContext(ctx).Find(&quizHistories)
	if res.Error != nil {
		return []QuizHistory{}, res.Error
	}

	return quizHistories, nil
}

func (r *repository) GetQuizHistoryByID(ctx context.Context, id uuid.UUID) (*QuizHistory, error) {
	var quizHistory QuizHistory
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&quizHistory)
	if res.Error != nil {
		return &QuizHistory{}, res.Error
	}

	return &quizHistory, nil
}

func (r *repository) GetQuizHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuizHistory, error) {
	var quizHistories []QuizHistory
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&quizHistories)
	if res.Error != nil {
		return []QuizHistory{}, res.Error
	}

	return quizHistories, nil
}

func (r *repository) GetQuizHistoriesByUserID(ctx context.Context, uid uuid.UUID) ([]QuizHistory, error) {
	var quizHistories []QuizHistory
	res := r.db.WithContext(ctx).Where("creator_id = ?", uid).Find(&quizHistories)
	if res.Error != nil {
		return []QuizHistory{}, res.Error
	}

	return quizHistories, nil
}

func (r *repository) GetQuizHistoryByQuizIDAndCreatedDate(ctx context.Context, quizID uuid.UUID, createdDate time.Time) (*QuizHistory, error) {
	var quizHistory QuizHistory
	res := r.db.WithContext(ctx).Where("quiz_id = ? AND created_date = ?", quizID, createdDate).First(&quizHistory)
	if res.Error != nil {
		return &QuizHistory{}, res.Error
	}

	return &quizHistory, nil
}

func (r *repository) UpdateQuizHistory(ctx context.Context, tx *gorm.DB, quizHistory *QuizHistory) (*QuizHistory, error) {
	res := tx.WithContext(ctx).Save(quizHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuizHistory{}, res.Error
	}

	return quizHistory, nil
}

func (r *repository) DeleteQuizHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&QuizHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

// ---------- Question Pool related repository methods ---------- //
func (r *repository) CreateQuestionPool(ctx context.Context, tx *gorm.DB, questionPool *QuestionPool) (*QuestionPool, error) {
	res := tx.WithContext(ctx).Create(questionPool)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionPool{}, res.Error
	}

	return questionPool, nil
}

func (r *repository) GetQuestionPoolByID(ctx context.Context, questionPoolID uuid.UUID) (*QuestionPool, error) {
	var questionPool QuestionPool
	res := r.db.WithContext(ctx).Where("id = ?", questionPoolID).Find(&questionPool)
	if res.Error != nil {
		return &QuestionPool{}, res.Error
	}

	return &questionPool, nil
}

func (r *repository) GetQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPool, error) {
	var questionPools []QuestionPool
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&questionPools)
	if res.Error != nil {
		return []QuestionPool{}, res.Error
	}

	return questionPools, nil
}

func (r *repository) GetDeleteQuestionPoolsByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPool, error) {
	var questionPools []QuestionPool
	res := r.db.WithContext(ctx).Unscoped().Where("quiz_id = ?", quizID).Find(&questionPools)
	if res.Error != nil {
		return []QuestionPool{}, res.Error
	}

	return questionPools, nil
}

func (r *repository) UpdateQuestionPool(ctx context.Context, tx *gorm.DB, questionPool *QuestionPool) (*QuestionPool, error) {
	res := tx.WithContext(ctx).Save(questionPool)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionPool{}, res.Error
	}

	return questionPool, nil
}

func (r *repository) DeleteQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&QuestionPool{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

func (r *repository) RestoreQuestionPool(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*QuestionPool, error) {
	var questionPool QuestionPool
	res := r.db.WithContext(ctx).Unscoped().First(&questionPool, id)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&questionPool).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	return &questionPool, nil
}

func (r *repository) CreateQuestionPoolHistory(ctx context.Context, tx *gorm.DB, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error) {
	res := tx.WithContext(ctx).Create(questionPoolHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionPoolHistory{}, res.Error
	}

	return questionPoolHistory, nil
}

func (r *repository) GetQuestionPoolHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionPoolHistory, error) {
	var questionPoolHistories []QuestionPoolHistory
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&questionPoolHistories)
	if res.Error != nil {
		return []QuestionPoolHistory{}, res.Error
	}

	return questionPoolHistories, nil
}

func (r *repository) UpdateQuestionPoolHistory(ctx context.Context, tx *gorm.DB, questionPoolHistory *QuestionPoolHistory) (*QuestionPoolHistory, error) {
	res := tx.WithContext(ctx).Save(questionPoolHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionPoolHistory{}, res.Error
	}

	return questionPoolHistory, nil
}

func (r *repository) DeleteQuestionPoolHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&QuestionPoolHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

// ---------- Question related repository methods ---------- //
func (r *repository) CreateQuestion(ctx context.Context, tx *gorm.DB, question *Question) (*Question, error) {
	res := tx.WithContext(ctx).Create(question)
	if res.Error != nil {
		tx.Rollback()
		return &Question{}, res.Error
	}

	return question, nil
}

func (r *repository) GetQuestions(ctx context.Context) ([]Question, error) {
	var questions []Question
	res := r.db.WithContext(ctx).Find(&questions)
	if res.Error != nil {
		return []Question{}, res.Error
	}

	return questions, nil
}

func (r *repository) GetQuestionByID(ctx context.Context, id uuid.UUID) (*Question, error) {
	var question Question
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&question)
	if res.Error != nil {
		return &Question{}, res.Error
	}
	return &question, nil
}

func (r *repository) GetQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error) {
	var questions []Question
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&questions)
	if res.Error != nil {
		return []Question{}, res.Error
	}

	return questions, nil
}

func (r *repository) GetDeleteQuestionsByQuizID(ctx context.Context, quizID uuid.UUID) ([]Question, error) {
	var questions []Question
	res := r.db.WithContext(ctx).Unscoped().Where("quiz_id = ?", quizID).Find(&questions)
	if res.Error != nil {
		return []Question{}, res.Error
	}

	return questions, nil
}

func (r *repository) GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*Question, error) {
	var question Question
	res := r.db.WithContext(ctx).Where(`quiz_id = ? AND "order" = ?`, quizID, order).First(&question)
	if res.Error != nil {
		return &Question{}, res.Error
	}

	return &question, nil
}

func (r *repository) GetQuestionCountByQuizID(ctx context.Context, quizID uuid.UUID) (int, error) {
	var count int64
	res := r.db.WithContext(ctx).Model(&Question{}).Where("quiz_id = ?", quizID).Count(&count)
	if res.Error != nil {
		return 0, res.Error
	}

	return int(count), nil
}

func (r *repository) UpdateQuestion(ctx context.Context, tx *gorm.DB, question *Question) (*Question, error) {
	res := tx.WithContext(ctx).Save(question)
	if res.Error != nil {
		tx.Rollback()
		return &Question{}, res.Error
	}

	return question, nil
}

func (r *repository) DeleteQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&Question{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

func (r *repository) RestoreQuestion(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*Question, error) {
	var question Question
	res := r.db.WithContext(ctx).Unscoped().First(&question, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&question).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	return &question, nil
}

func (r *repository) CreateQuestionHistory(ctx context.Context, tx *gorm.DB, questionHistory *QuestionHistory) (*QuestionHistory, error) {
	res := tx.WithContext(ctx).Create(questionHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionHistory{}, res.Error
	}

	return questionHistory, nil
}

func (r *repository) GetQuestionHistories(ctx context.Context) ([]QuestionHistory, error) {
	var questionHistories []QuestionHistory
	res := r.db.WithContext(ctx).Find(&questionHistories)
	if res.Error != nil {
		return []QuestionHistory{}, res.Error
	}

	return questionHistories, nil
}

func (r *repository) GetQuestionHistoryByID(ctx context.Context, id uuid.UUID) (*QuestionHistory, error) {
	var questionHistory QuestionHistory
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&questionHistory)
	if res.Error != nil {
		return &QuestionHistory{}, res.Error
	}

	return &questionHistory, nil
}

func (r *repository) GetQuestionHistoriesByQuizID(ctx context.Context, quizID uuid.UUID) ([]QuestionHistory, error) {
	var questionHistories []QuestionHistory
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&questionHistories)
	if res.Error != nil {
		return []QuestionHistory{}, res.Error
	}

	return questionHistories, nil
}

func (r *repository) GetQuestionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]QuestionHistory, error) {
	var questionHistories []QuestionHistory
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&questionHistories)
	if res.Error != nil {
		return []QuestionHistory{}, res.Error
	}

	return questionHistories, nil
}

func (r *repository) GetQuestionHistoryByQuestionIDAndCreatedDate(ctx context.Context, questionID uuid.UUID, createdDate time.Time) (*QuestionHistory, error) {
	var questionHistory QuestionHistory
	res := r.db.WithContext(ctx).Where("question_id = ? AND created_at = ?", questionID, createdDate).First(&questionHistory)
	if res.Error != nil {
		return &QuestionHistory{}, res.Error
	}

	return &questionHistory, nil
}

func (r *repository) UpdateQuestionHistory(ctx context.Context, tx *gorm.DB, questionHistory *QuestionHistory) (*QuestionHistory, error) {
	res := tx.WithContext(ctx).Save(questionHistory)
	if res.Error != nil {
		tx.Rollback()
		return &QuestionHistory{}, res.Error
	}

	return questionHistory, nil
}

func (r *repository) DeleteQuestionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&QuestionHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

// ---------- Options related repository methods ---------- //
// Choice related repository methods
func (r *repository) CreateChoiceOption(ctx context.Context, tx *gorm.DB, optionChoice *ChoiceOption) (*ChoiceOption, error) {
	res := tx.WithContext(ctx).Create(optionChoice)
	if res.Error != nil {
		tx.Rollback()
		return &ChoiceOption{}, res.Error
	}

	return optionChoice, nil
}

func (r *repository) GetChoiceOptionByID(ctx context.Context, id uuid.UUID) (*ChoiceOption, error) {
	var optionChoice ChoiceOption
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionChoice)
	if res.Error != nil {
		return &ChoiceOption{}, res.Error
	}

	return &optionChoice, nil
}

func (r *repository) GetChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error) {
	var optionChoices []ChoiceOption
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionChoices)
	if res.Error != nil {
		return []ChoiceOption{}, res.Error
	}

	return optionChoices, nil
}

func (r *repository) GetDeleteChoiceOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error) {
	var optionChoices []ChoiceOption
	res := r.db.WithContext(ctx).Unscoped().Where("question_id = ?", questionID).Find(&optionChoices)
	if res.Error != nil {
		return []ChoiceOption{}, res.Error
	}

	return optionChoices, nil
}

func (r *repository) GetChoiceAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error) {
	var optionChoices []ChoiceOption
	res := r.db.WithContext(ctx).Where("question_id = ? AND correct = ?", questionID, true).Find(&optionChoices)
	if res.Error != nil {
		return []ChoiceOption{}, res.Error
	}

	return optionChoices, nil
}

func (r *repository) UpdateChoiceOption(ctx context.Context, tx *gorm.DB, optionChoice *ChoiceOption) (*ChoiceOption, error) {
	res := tx.WithContext(ctx).Save(optionChoice)
	if res.Error != nil {
		tx.Rollback()
		return &ChoiceOption{}, res.Error
	}

	return optionChoice, nil
}

func (r *repository) DeleteChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&ChoiceOption{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

func (r *repository) RestoreChoiceOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*ChoiceOption, error) {
	var optionChoice ChoiceOption
	res := r.db.WithContext(ctx).Unscoped().First(&optionChoice, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&optionChoice).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	return &optionChoice, nil
}

func (r *repository) CreateChoiceOptionHistory(ctx context.Context, tx *gorm.DB, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error) {
	res := tx.WithContext(ctx).Create(optionChoiceHistory)
	if res.Error != nil {
		tx.Rollback()
		return &ChoiceOptionHistory{}, res.Error
	}

	return optionChoiceHistory, nil
}

func (r *repository) GetChoiceOptionHistoryByID(ctx context.Context, id uuid.UUID) (*ChoiceOptionHistory, error) {
	var optionChoiceHistory ChoiceOptionHistory
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionChoiceHistory)
	if res.Error != nil {
		return &ChoiceOptionHistory{}, res.Error
	}

	return &optionChoiceHistory, nil
}

func (r *repository) GetOptionChoiceHistories(ctx context.Context) ([]ChoiceOptionHistory, error) {
	var optionChoiceHistories []ChoiceOptionHistory
	res := r.db.WithContext(ctx).Find(&optionChoiceHistories)
	if res.Error != nil {
		return []ChoiceOptionHistory{}, res.Error
	}

	return optionChoiceHistories, nil
}

func (r *repository) GetChoiceOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOptionHistory, error) {
	var optionChoiceHistories []ChoiceOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionChoiceHistories)
	if res.Error != nil {
		return []ChoiceOptionHistory{}, res.Error
	}

	return optionChoiceHistories, nil
}

func (r *repository) GetChoiceOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*ChoiceOptionHistory, error) {
	var optionChoiceHistory ChoiceOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ? and content = ?", questionID, content).Find(&optionChoiceHistory)
	if res.Error != nil {
		return &ChoiceOptionHistory{}, res.Error
	}

	return &optionChoiceHistory, nil
}

func (r *repository) UpdateChoiceOptionHistory(ctx context.Context, tx *gorm.DB, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error) {
	res := tx.WithContext(ctx).Save(optionChoiceHistory)
	if res.Error != nil {
		tx.Rollback()
		return &ChoiceOptionHistory{}, res.Error
	}

	return optionChoiceHistory, nil
}

func (r *repository) DeleteChoiceOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&ChoiceOptionHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

// Text related repository methods
func (r *repository) CreateTextOption(ctx context.Context, tx *gorm.DB, optionText *TextOption) (*TextOption, error) {
	res := tx.WithContext(ctx).Create(optionText)
	if res.Error != nil {
		tx.Rollback()
		return &TextOption{}, res.Error
	}

	return optionText, nil
}

func (r *repository) GetTextOptionByID(ctx context.Context, id uuid.UUID) (*TextOption, error) {
	var optionText TextOption
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionText)
	if res.Error != nil {
		return &TextOption{}, res.Error
	}

	return &optionText, nil
}

func (r *repository) GetTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error) {
	var optionTexts []TextOption
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionTexts)
	if res.Error != nil {
		return []TextOption{}, res.Error
	}

	return optionTexts, nil
}

func (r *repository) GetDeleteTextOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error) {
	var optionTexts []TextOption
	res := r.db.WithContext(ctx).Unscoped().Where("question_id = ?", questionID).Find(&optionTexts)
	if res.Error != nil {
		return []TextOption{}, res.Error
	}

	return optionTexts, nil
}

func (r *repository) UpdateTextOption(ctx context.Context, tx *gorm.DB, optionText *TextOption) (*TextOption, error) {
	res := tx.WithContext(ctx).Save(optionText)
	if res.Error != nil {
		tx.Rollback()
		return &TextOption{}, res.Error
	}

	return optionText, nil
}

func (r *repository) DeleteTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&TextOption{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

func (r *repository) RestoreTextOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*TextOption, error) {
	var optionText TextOption
	res := r.db.WithContext(ctx).Unscoped().First(&optionText, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&optionText).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}

	return &optionText, nil
}

func (r *repository) CreateTextOptionHistory(ctx context.Context, tx *gorm.DB, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error) {
	res := tx.WithContext(ctx).Create(optionTextHistory)
	if res.Error != nil {
		tx.Rollback()
		return &TextOptionHistory{}, res.Error
	}

	return optionTextHistory, nil
}

func (r *repository) GetTextOptionHistoryByID(ctx context.Context, id uuid.UUID) (*TextOptionHistory, error) {
	var optionTextHistory TextOptionHistory
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionTextHistory)
	if res.Error != nil {
		return &TextOptionHistory{}, res.Error
	}

	return &optionTextHistory, nil
}

func (r *repository) GetTextOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOptionHistory, error) {
	var optionTextHistories []TextOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionTextHistories)
	if res.Error != nil {
		return []TextOptionHistory{}, res.Error
	}

	return optionTextHistories, nil
}

func (r *repository) GetTextOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*TextOptionHistory, error) {
	var optionTextHistory TextOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ? AND content = ?", questionID, content).Find(&optionTextHistory)
	if res.Error != nil {
		return &TextOptionHistory{}, res.Error
	}

	return &optionTextHistory, nil
}

func (r *repository) UpdateTextOptionHistory(ctx context.Context, tx *gorm.DB, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error) {
	res := tx.WithContext(ctx).Save(optionTextHistory)
	if res.Error != nil {
		tx.Rollback()
		return &TextOptionHistory{}, res.Error
	}

	return optionTextHistory, nil
}

func (r *repository) DeleteTextOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&TextOptionHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}

// Matching related repository methods
// Option Matching
func (r *repository) CreateMatchingOption(ctx context.Context, tx *gorm.DB, optionMatching *MatchingOption) (*MatchingOption, error) {
	res := tx.WithContext(ctx).Create(optionMatching)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingOption{}, res.Error
	}
	return optionMatching, nil
}

func (r *repository) GetMatchingOptionByID(ctx context.Context, id uuid.UUID) (*MatchingOption, error) {
	var optionMatching MatchingOption
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionMatching)
	if res.Error != nil {
		return &MatchingOption{}, res.Error
	}
	return &optionMatching, nil
}

func (r *repository) GetMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOption, error) {
	var optionMatchings []MatchingOption
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionMatchings)
	if res.Error != nil {
		return []MatchingOption{}, res.Error
	}
	return optionMatchings, nil
}

func (r *repository) GetDeleteMatchingOptionsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOption, error) {
	var optionMatchings []MatchingOption
	res := r.db.WithContext(ctx).Unscoped().Where("question_id = ?", questionID).Find(&optionMatchings)
	if res.Error != nil {
		return []MatchingOption{}, res.Error
	}
	return optionMatchings, nil
}

func (r *repository) GetMatchingOptionByQuestionIDAndOrder(ctx context.Context, questionID uuid.UUID, order int) (*MatchingOption, error) {
	var optionMatching MatchingOption
	res := r.db.WithContext(ctx).Where(`question_id = ? AND "order" = ?`, questionID, order).First(&optionMatching)
	if res.Error != nil {
		return &MatchingOption{}, res.Error
	}
	return &optionMatching, nil
}

func (r *repository) UpdateMatchingOption(ctx context.Context, tx *gorm.DB, optionMatching *MatchingOption) (*MatchingOption, error) {
	res := tx.WithContext(ctx).Save(optionMatching)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingOption{}, res.Error
	}
	return optionMatching, nil
}

func (r *repository) DeleteMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&MatchingOption{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}

func (r *repository) RestoreMatchingOption(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*MatchingOption, error) {
	var optionMatching MatchingOption
	res := r.db.WithContext(ctx).Unscoped().First(&optionMatching, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&optionMatching).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	return &optionMatching, nil
}

func (r *repository) CreateMatchingOptionHistory(ctx context.Context, tx *gorm.DB, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error) {
	res := tx.WithContext(ctx).Create(optionMatchingHistory)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingOptionHistory{}, res.Error
	}
	return optionMatchingHistory, nil
}

func (r *repository) GetMatchingOptionHistoryByOptionMatchingID(ctx context.Context, optionMatchingID uuid.UUID) (*MatchingOptionHistory, error) {
	var optionMatchingHistory MatchingOptionHistory
	res := r.db.WithContext(ctx).Where("option_matching_id = ?", optionMatchingID).First(&optionMatchingHistory)
	if res.Error != nil {
		return &MatchingOptionHistory{}, res.Error
	}
	return &optionMatchingHistory, nil
}

func (r *repository) GetMatchingOptionHistoryByID(ctx context.Context, id uuid.UUID) (*MatchingOptionHistory, error) {
	var optionMatchingHistory MatchingOptionHistory
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&optionMatchingHistory)
	if res.Error != nil {
		return &MatchingOptionHistory{}, res.Error
	}
	return &optionMatchingHistory, nil
}

func (r *repository) GetOptionMatchingHistories(ctx context.Context) ([]MatchingOptionHistory, error) {
	var optionMatchingHistories []MatchingOptionHistory
	res := r.db.WithContext(ctx).Find(&optionMatchingHistories)
	if res.Error != nil {
		return []MatchingOptionHistory{}, res.Error
	}
	return optionMatchingHistories, nil
}

func (r *repository) GetMatchingOptionHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingOptionHistory, error) {
	var optionMatchingHistories []MatchingOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&optionMatchingHistories)
	if res.Error != nil {
		return []MatchingOptionHistory{}, res.Error
	}
	return optionMatchingHistories, nil
}

func (r *repository) GetMatchingOptionHistoryByQuestionIDAndContent(ctx context.Context, questionID uuid.UUID, content string) (*MatchingOptionHistory, error) {
	var optionMatchingHistory MatchingOptionHistory
	res := r.db.WithContext(ctx).Where("question_id = ? AND content = ?", questionID, content).Find(&optionMatchingHistory)
	if res.Error != nil {
		return &MatchingOptionHistory{}, res.Error
	}
	return &optionMatchingHistory, nil
}

func (r *repository) UpdateMatchingOptionHistory(ctx context.Context, tx *gorm.DB, optionMatchingHistory *MatchingOptionHistory) (*MatchingOptionHistory, error) {
	res := tx.WithContext(ctx).Save(optionMatchingHistory)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingOptionHistory{}, res.Error
	}
	return optionMatchingHistory, nil
}

func (r *repository) DeleteMatchingOptionHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&MatchingOptionHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}

// Answer Matching
func (r *repository) CreateMatchingAnswer(ctx context.Context, tx *gorm.DB, answerMatching *MatchingAnswer) (*MatchingAnswer, error) {
	res := tx.WithContext(ctx).Create(answerMatching)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingAnswer{}, res.Error
	}
	return answerMatching, nil
}

func (r *repository) UpdateMatchingAnswer(ctx context.Context, tx *gorm.DB, answerMatching *MatchingAnswer) (*MatchingAnswer, error) {
	res := tx.WithContext(ctx).Save(answerMatching)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingAnswer{}, res.Error
	}
	return answerMatching, nil
}

func (r *repository) DeleteMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&MatchingAnswer{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}

func (r *repository) RestoreMatchingAnswer(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*MatchingAnswer, error) {
	var answerMatching MatchingAnswer
	res := r.db.WithContext(ctx).Unscoped().First(&answerMatching, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = tx.WithContext(ctx).Unscoped().Model(&answerMatching).Update("deleted_at", nil)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	return &answerMatching, nil
}

func (r *repository) GetMatchingAnswerByID(ctx context.Context, id uuid.UUID) (*MatchingAnswer, error) {
	var answerMatching MatchingAnswer
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&answerMatching)
	if res.Error != nil {
		return &MatchingAnswer{}, res.Error
	}
	return &answerMatching, nil
}

func (r *repository) GetMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswer, error) {
	var answerMatchings []MatchingAnswer
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&answerMatchings)
	if res.Error != nil {
		return []MatchingAnswer{}, res.Error
	}
	return answerMatchings, nil
}

func (r *repository) GetDeleteMatchingAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswer, error) {
	var answerMatchings []MatchingAnswer
	res := r.db.WithContext(ctx).Unscoped().Where("question_id = ?", questionID).Find(&answerMatchings)
	if res.Error != nil {
		return []MatchingAnswer{}, res.Error
	}
	return answerMatchings, nil
}

func (r *repository) CreateMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error) {
	res := tx.WithContext(ctx).Create(answerMatchingHistory)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingAnswerHistory{}, res.Error
	}
	return answerMatchingHistory, nil
}

func (r *repository) GetMatchingAnswerHistoriesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]MatchingAnswerHistory, error) {
	var answerMatchingHistories []MatchingAnswerHistory
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&answerMatchingHistories)
	if res.Error != nil {
		return []MatchingAnswerHistory{}, res.Error
	}
	return answerMatchingHistories, nil
}

func (r *repository) GetMatchingAnswerHistoryByPromptIDAndOptionID(ctx context.Context, promptID uuid.UUID, optionID uuid.UUID) (*MatchingAnswerHistory, error) {
	var answerMatchingHistory MatchingAnswerHistory
	res := r.db.WithContext(ctx).Where("prompt_id = ? AND option_id = ?", promptID, optionID).Find(&answerMatchingHistory)
	if res.Error != nil {
		return &MatchingAnswerHistory{}, res.Error
	}
	return &answerMatchingHistory, nil
}

func (r *repository) UpdateMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, answerMatchingHistory *MatchingAnswerHistory) (*MatchingAnswerHistory, error) {
	res := tx.WithContext(ctx).Save(answerMatchingHistory)
	if res.Error != nil {
		tx.Rollback()
		return &MatchingAnswerHistory{}, res.Error
	}
	return answerMatchingHistory, nil
}

func (r *repository) DeleteMatchingAnswerHistory(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	res := tx.WithContext(ctx).Delete(&MatchingAnswerHistory{}, id)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}
	return nil
}

func (r *repository) GetLatestQuizHistoryByQuizID(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	var qh QuizHistory
	res := r.db.WithContext(ctx).Order("created_at desc").Where("quiz_id = ?", id).First(&qh)
	if res.Error != nil {
		return nil, res.Error
	}
	return &qh.ID, nil
}
