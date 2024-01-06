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
	err := tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *repository) GetParticipantsByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]Participant, error) {
	var participants []Participant
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ?", liveQuizSessionID).Order("marks desc").Find(&participants)
	if res.Error != nil {
		return []Participant{}, res.Error
	}
	return participants, nil
}

func (r *repository) GetParticipantsByLiveQuizSessionIDCustom(ctx context.Context, liveQuizSessionID uuid.UUID, order string, limit int) ([]Participant, error) {
	var participants []Participant
	res := r.db.WithContext(ctx).
		Where("live_quiz_session_id = ?", liveQuizSessionID).
		Order(order + " desc").
		Limit(limit).
		Find(&participants)

	if res.Error != nil {
		return []Participant{}, res.Error
	}
	return participants, nil
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
