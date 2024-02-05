package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(repo Repository) Service {
	return &service{
		Repository: repo,
		timeout:    time.Duration(3) * time.Second,
	}
}

func (s *service) BeginTransaction(ctx context.Context) (*gorm.DB, error) {
	tx, err := s.Repository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *service) CommitTransaction(ctx context.Context,tx *gorm.DB) (error) {
	err := s.Repository.CommitTransaction(tx)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetParticipantsByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]ParticipantResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participants, err := s.Repository.GetParticipantsByLiveQuizSessionID(ctx, liveQuizSessionID)
	if err != nil {
		return nil, err
	}

	var res []ParticipantResponse
	for _, p := range participants {
		res = append(res, ParticipantResponse{
			Participant: Participant{
				ID: p.ID,
				UserID: p.UserID,
				LiveQuizSessionID: p.LiveQuizSessionID,
				Status: p.Status,
				Name: p.Name,
				Marks: p.Marks,
			},
		})
	}

	return res, nil
}

func (s *service) GetParticipantsByLiveQuizSessionIDCustom(ctx context.Context, liveQuizSessionID uuid.UUID, orderBy string, limit int) ([]ParticipantResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participants, err := s.Repository.GetParticipantsByLiveQuizSessionIDCustom(ctx, liveQuizSessionID, orderBy, limit)
	if err != nil {
		return nil, err
	}

	var res []ParticipantResponse
	for _, p := range participants {
		res = append(res, ParticipantResponse{
			Participant: Participant{
				ID: p.ID,
				UserID: p.UserID,
				LiveQuizSessionID: p.LiveQuizSessionID,
				Status: p.Status,
				Name: p.Name,
				Marks: p.Marks,
			},
		})
	}

	return res, nil
}

func (s *service) GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveQuizSessionID uuid.UUID) ([]LiveAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	liveAnswers, err := s.Repository.GetAnswerResponseByLiveQuizSessionID(ctx, liveQuizSessionID)
	if err != nil {
		return nil, err
	}

	var res []LiveAnswerResponse
	for _, liveAnswer := range liveAnswers {
		res = append(res, LiveAnswerResponse{
			AnswerResponse: AnswerResponse{
				ID:                liveAnswer.ID,
				LiveQuizSessionID: liveAnswer.LiveQuizSessionID,
				ParticipantID:     liveAnswer.ParticipantID,
				Type:              liveAnswer.Type,
				QuestionID:        liveAnswer.QuestionID,
				Answer:            liveAnswer.Answer,
				CreatedAt:         liveAnswer.CreatedAt,
				UpdatedAt:         liveAnswer.UpdatedAt,
				DeletedAt:         liveAnswer.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetAnswerResponseByQuestionID(ctx context.Context, questionID uuid.UUID) ([]LiveAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	liveAnswers, err := s.Repository.GetAnswerResponseByQuestionID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	var res []LiveAnswerResponse
	for _, liveAnswer := range liveAnswers {
		res = append(res, LiveAnswerResponse{
			AnswerResponse: AnswerResponse{
				ID:                liveAnswer.ID,
				LiveQuizSessionID: liveAnswer.LiveQuizSessionID,
				ParticipantID:     liveAnswer.ParticipantID,
				Type:              liveAnswer.Type,
				QuestionID:        liveAnswer.QuestionID,
				Answer:            liveAnswer.Answer,
				CreatedAt:         liveAnswer.CreatedAt,
				UpdatedAt:         liveAnswer.UpdatedAt,
				DeletedAt:         liveAnswer.DeletedAt,
			},
		})
	}

	return res, nil
}

func (s *service) GetAnswerResponseByParticipantID(ctx context.Context, participantID uuid.UUID) ([]LiveAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	liveAnswers, err := s.Repository.GetAnswerResponseByParticipantID(ctx, participantID)
	if err != nil {
		return nil, err
	}

	var res []LiveAnswerResponse
	for _, liveAnswer := range liveAnswers {
		res = append(res, LiveAnswerResponse{
			AnswerResponse: AnswerResponse{
				ID:                liveAnswer.ID,
				LiveQuizSessionID: liveAnswer.LiveQuizSessionID,
				ParticipantID:     liveAnswer.ParticipantID,
				Type:              liveAnswer.Type,
				QuestionID:        liveAnswer.QuestionID,
				Answer:            liveAnswer.Answer,
				CreatedAt:         liveAnswer.CreatedAt,
				UpdatedAt:         liveAnswer.UpdatedAt,
				DeletedAt:         liveAnswer.DeletedAt,
			},
		})
	}

	return res, nil
}
