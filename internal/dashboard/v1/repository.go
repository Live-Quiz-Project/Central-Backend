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

func (r *repository) GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(ctx context.Context, liveQuizSessionID uuid.UUID, questionID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ? AND question_id = ?", liveQuizSessionID, questionID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetAnswerResponsesByLiveQuizSessionIDAndParticipantID(ctx context.Context, liveQuizSessionID uuid.UUID, participantID uuid.UUID) ([]AnswerResponse, error) {
	var answerResponses []AnswerResponse
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ? AND participant_id = ?", liveQuizSessionID, participantID).Find(&answerResponses)
	if res.Error != nil {
		return []AnswerResponse{}, res.Error
	}
	return answerResponses, nil
}

func (r *repository) GetParticipantByID(ctx context.Context, participantID uuid.UUID) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("id = ?", participantID).Find(&participant)
	if res.Error != nil {
		return &Participant{}, res.Error
	}
	return &participant, nil
}

func(r *repository) GetOrderParticipantsByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]Participant, error) {
	var participant []Participant
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ?", liveQuizSessionID).Order("marks DESC,name ASC").Find(&participant)
	if res.Error != nil {
		return []Participant{}, res.Error
	}
	return participant, nil
}
