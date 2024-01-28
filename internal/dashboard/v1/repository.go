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

func (r *repository) BeginTransaction() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

func (r *repository) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *repository) GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveSessionID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ?", liveSessionID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("participant_id = ?", participantID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetAnswerResponsesByLiveQuizSessionIDAndQuestionID(ctx context.Context, liveQuizSessionID uuid.UUID, questionID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ? AND question_id = ?", liveQuizSessionID, questionID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetParticipantByID(ctx context.Context, participantID uuid.UUID) (*ParticipantResponse, error) {
	var participantResponse ParticipantResponse
	res := r.db.WithContext(ctx).Where("id = ?", participantID).Find(&participantResponse)
	if res.Error != nil {
		return &ParticipantResponse{}, res.Error
	}
	return &participantResponse, nil
}

// For Testing
func (r *repository) CreateAnswerResponse(ctx context.Context, tx *gorm.DB, answerResponse *AnswerResponse) (*AnswerResponse, error) {
	res := tx.WithContext(ctx).Create(answerResponse)
	if res.Error != nil {
		tx.Rollback()
		return nil, res.Error
	}
	return answerResponse, nil
}
