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

func (s *service) GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(ctx context.Context, liveQuizSessionID uuid.UUID, questionID uuid.UUID) ([]LiveAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	liveAnswers, err := s.Repository.GetAnswerResponsesByLiveQuizSessionIDAndQuestionHistoryID(ctx, liveQuizSessionID, questionID)
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

func (s *service) GetParticipantByID(ctx context.Context, liveQuizSessionID uuid.UUID) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participant, err := s.Repository.GetParticipantByID(c, liveQuizSessionID)
	if err != nil {
		return nil, err
	}

	return &Participant{
		ID:                participant.ID,
		UserID:            participant.UserID,
		LiveQuizSessionID: participant.LiveQuizSessionID,
		Status:            participant.Status,
		Name:              participant.Name,
		Marks:             participant.Marks,
	}, nil
}
