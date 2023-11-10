package v1

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateChoiceResponse(ctx context.Context, rs *ChoiceResponse) (*ChoiceResponse, error) {
	res := r.db.WithContext(ctx).Create(rs)
	if res.Error != nil {
		return &ChoiceResponse{}, res.Error
	}
	return rs, nil
}

func (r *repository) GetChoiceResponsesByQuizID(ctx context.Context, quizID string) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}
	return responses, nil
}

func (r *repository) GetChoiceResponsesByQuizIDAndQuestionID(ctx context.Context, quizID string, questionID string) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("quiz_id = ? AND question_id = ?", quizID, questionID).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}
	return responses, nil
}

func (r *repository) GetChoiceResponsesByQuizIDAndUserID(ctx context.Context, quizID string, uid string) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("quiz_id = ? AND user_id = ?", quizID, uid).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}
	return responses, nil
}

func (r *repository) GetChoiceResponsesByQuizIDAndQuestionIDAndUserID(ctx context.Context, quizID string, questionID string, uid string) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("quiz_id = ? AND question_id = ? AND user_id = ?", quizID, questionID, uid).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}
	return responses, nil
}

func (r *repository) UpdateChoiceResponse(ctx context.Context, rs *ChoiceResponse) (*ChoiceResponse, error) {
	res := r.db.WithContext(ctx).Save(rs)
	if res.Error != nil {
		return &ChoiceResponse{}, res.Error
	}
	return rs, nil
}

func (r *repository) DeleteChoiceResponse(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&ChoiceResponse{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repository) RestoreChoiceResponse(ctx context.Context, id uuid.UUID) error {
	var choiceResponse ChoiceResponse
	res := r.db.WithContext(ctx).Unscoped().First(&choiceResponse, id)
	if res.Error != nil {
		return res.Error
	}
	res = r.db.WithContext(ctx).Unscoped().Model(&choiceResponse).Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
