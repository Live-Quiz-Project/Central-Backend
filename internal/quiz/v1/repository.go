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

// ---------- Quiz related repository methods ---------- //
func (r *repository) CreateQuiz(ctx context.Context, quiz *Quiz) (*Quiz, error) {
	res := r.db.WithContext(ctx).Create(quiz)
	if res.Error != nil {
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

func (r *repository) UpdateQuiz(ctx context.Context, quiz *Quiz) (*Quiz, error) {
	res := r.db.WithContext(ctx).Save(quiz)
	if res.Error != nil {
		return &Quiz{}, res.Error
	}

	return quiz, nil
}

func (r *repository) DeleteQuiz(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&Quiz{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *repository) RestoreQuiz(ctx context.Context, id uuid.UUID) (*Quiz, error) {
	var quiz Quiz
	res := r.db.WithContext(ctx).Unscoped().First(&quiz, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Unscoped().Model(&quiz).Updates(Quiz{DeletedAt: gorm.DeletedAt{}})
	if res.Error != nil {
		return nil, res.Error
	}

	return &quiz, nil
}

func (r *repository) CreateQuizHistory(ctx context.Context, quizHistory *QuizHistory) (*QuizHistory, error) {
	res := r.db.WithContext(ctx).Create(quizHistory)
	if res.Error != nil {
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

func (r *repository) GetQuizHistoryByQuizIDAndCreatedDate(ctx context.Context, quizID uuid.UUID, createdDate time.Time) (*QuizHistory, error) {
	var quizHistory QuizHistory
	res := r.db.WithContext(ctx).Where("quiz_id = ? AND created_date = ?", quizID, createdDate).First(&quizHistory)
	if res.Error != nil {
		return &QuizHistory{}, res.Error
	}

	return &quizHistory, nil
}

func (r *repository) UpdateQuizHistory(ctx context.Context, quizHistory *QuizHistory) (*QuizHistory, error) {
	res := r.db.WithContext(ctx).Save(quizHistory)
	if res.Error != nil {
		return &QuizHistory{}, res.Error
	}

	return quizHistory, nil
}

func (r *repository) DeleteQuizHistory(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&QuizHistory{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// ---------- Question related repository methods ---------- //
func (r *repository) CreateQuestion(ctx context.Context, question *Question) (*Question, error) {
	res := r.db.WithContext(ctx).Create(question)
	if res.Error != nil {
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

func (r *repository) UpdateQuestion(ctx context.Context, question *Question) (*Question, error) {
	res := r.db.WithContext(ctx).Save(question)
	if res.Error != nil {
		return &Question{}, res.Error
	}

	return question, nil
}

func (r *repository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&Question{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *repository) RestoreQuestion(ctx context.Context, id uuid.UUID) (*Question, error) {
	var question Question
	res := r.db.WithContext(ctx).Unscoped().First(&question, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Unscoped().Model(&question).Updates(Question{DeletedAt: gorm.DeletedAt{}})
	if res.Error != nil {
		return nil, res.Error
	}

	return &question, nil
}

func (r *repository) CreateQuestionHistory(ctx context.Context, questionHistory *QuestionHistory) (*QuestionHistory, error) {
	res := r.db.WithContext(ctx).Create(questionHistory)
	if res.Error != nil {
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
	res := r.db.WithContext(ctx).Where("question_id = ? AND created_date = ?", questionID, createdDate).First(&questionHistory)
	if res.Error != nil {
		return &QuestionHistory{}, res.Error
	}

	return &questionHistory, nil
}

func (r *repository) UpdateQuestionHistory(ctx context.Context, questionHistory *QuestionHistory) (*QuestionHistory, error) {
	res := r.db.WithContext(ctx).Save(questionHistory)
	if res.Error != nil {
		return &QuestionHistory{}, res.Error
	}

	return questionHistory, nil
}

func (r *repository) DeleteQuestionHistory(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&QuestionHistory{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// ---------- Options related repository methods ---------- //
// Choice related repository methods
func (r *repository) CreateChoiceOption(ctx context.Context, optionChoice *ChoiceOption) (*ChoiceOption, error) {
	res := r.db.WithContext(ctx).Create(optionChoice)
	if res.Error != nil {
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

func (r *repository) GetChoiceAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceOption, error) {
	var optionChoices []ChoiceOption
	res := r.db.WithContext(ctx).Where("question_id = ? AND correct = ?", questionID, true).Find(&optionChoices)
	if res.Error != nil {
		return []ChoiceOption{}, res.Error
	}

	return optionChoices, nil
}

func (r *repository) UpdateChoiceOption(ctx context.Context, optionChoice *ChoiceOption) (*ChoiceOption, error) {
	res := r.db.WithContext(ctx).Save(optionChoice)
	if res.Error != nil {
		return &ChoiceOption{}, res.Error
	}

	return optionChoice, nil
}

func (r *repository) DeleteChoiceOption(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&ChoiceOption{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *repository) RestoreChoiceOption(ctx context.Context, id uuid.UUID) (*ChoiceOption, error) {
	var optionChoice ChoiceOption
	res := r.db.WithContext(ctx).Unscoped().First(&optionChoice, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Unscoped().Model(&optionChoice).Updates(ChoiceOption{DeletedAt: gorm.DeletedAt{}})
	if res.Error != nil {
		return nil, res.Error
	}

	return &optionChoice, nil
}

func (r *repository) CreateChoiceOptionHistory(ctx context.Context, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error) {
	res := r.db.WithContext(ctx).Create(optionChoiceHistory)
	if res.Error != nil {
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

func (r *repository) UpdateChoiceOptionHistory(ctx context.Context, optionChoiceHistory *ChoiceOptionHistory) (*ChoiceOptionHistory, error) {
	res := r.db.WithContext(ctx).Save(optionChoiceHistory)
	if res.Error != nil {
		return &ChoiceOptionHistory{}, res.Error
	}

	return optionChoiceHistory, nil
}

func (r *repository) DeleteChoiceOptionHistory(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&ChoiceOptionHistory{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// Text related repository methods
func (r *repository) CreateTextOption(ctx context.Context, optionText *TextOption) (*TextOption, error) {
	res := r.db.WithContext(ctx).Create(optionText)
	if res.Error != nil {
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

func (r *repository) GetTextAnswersByQuestionID(ctx context.Context, questionID uuid.UUID) ([]TextOption, error) {
	var optionTexts []TextOption
	res := r.db.WithContext(ctx).Where("question_id = ? AND is_correct = ?", questionID, true).Find(&optionTexts)
	if res.Error != nil {
		return []TextOption{}, res.Error
	}

	return optionTexts, nil
}

func (r *repository) UpdateTextOption(ctx context.Context, optionText *TextOption) (*TextOption, error) {
	res := r.db.WithContext(ctx).Save(optionText)
	if res.Error != nil {
		return &TextOption{}, res.Error
	}

	return optionText, nil
}

func (r *repository) DeleteTextOption(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&TextOption{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *repository) RestoreTextOption(ctx context.Context, id uuid.UUID) (*TextOption, error) {
	var optionText TextOption
	res := r.db.WithContext(ctx).Unscoped().First(&optionText, id)
	if res.Error != nil {
		return nil, res.Error
	}

	res = r.db.WithContext(ctx).Unscoped().Model(&optionText).Updates(TextOption{DeletedAt: gorm.DeletedAt{}})
	if res.Error != nil {
		return nil, res.Error
	}

	return &optionText, nil
}

func (r *repository) CreateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error) {
	res := r.db.WithContext(ctx).Create(optionTextHistory)
	if res.Error != nil {
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

func (r *repository) UpdateTextOptionHistory(ctx context.Context, optionTextHistory *TextOptionHistory) (*TextOptionHistory, error) {
	res := r.db.WithContext(ctx).Save(optionTextHistory)
	if res.Error != nil {
		return &TextOptionHistory{}, res.Error
	}

	return optionTextHistory, nil
}

func (r *repository) DeleteTextOptionHistory(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&TextOptionHistory{}, id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
