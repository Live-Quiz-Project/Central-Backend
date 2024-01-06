package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
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

func (s *service) CreateAnswerResponse(ctx context.Context, req *LiveAnswerRequest) (*LiveAnswerResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	tx, err := s.Repository.BeginTransaction()
	if err != nil {
		return nil, err
	}

	la := &AnswerResponse{
		ID:                uuid.New(),
		LiveQuizSessionID: req.LiveQuizSessionID,
		ParticipantID:     req.ParticipantID,
		Type:              req.Type,
		QuestionID:        req.QuestionID,
		Answer:            req.Answer,
	}

	liveAnswer, err := s.Repository.CreateAnswerResponse(c, tx, la)
	if err != nil {
		return &LiveAnswerResponse{}, err
	}

	s.Repository.CommitTransaction(tx)

	return &LiveAnswerResponse{
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
	}, nil
}

func (s *service) GetAnswerResponseByLiveQuizSessionID(ctx context.Context, liveSessionID uuid.UUID) ([]LiveAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	liveAnswers, err := s.Repository.GetAnswerResponseByLiveQuizSessionID(ctx, liveSessionID)
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
