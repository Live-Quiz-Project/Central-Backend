package v1

import (
	"context"
	"time"

	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
)

type service struct {
	quiz q.Repository
	Repository
	timeout time.Duration
}

func NewService(qRepo q.Repository, lRepo Repository) Service {
	return &service{
		quiz:       qRepo,
		Repository: lRepo,
		timeout:    time.Duration(3) * time.Second,
	}
}

// ---------- Live Quiz Session related service methods ---------- //
func (s *service) CreateLiveQuizSession(ctx context.Context, req *CreateLiveQuizSessionRequest, id uuid.UUID, code string, hostID uuid.UUID) (*CreateLiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess := &Session{
		ID:                  id,
		HostID:              hostID,
		QuizID:              req.QuizID,
		Status:              util.Idle,
		ExemptedQuestionIDs: nil,
	}

	sess, err := s.Repository.CreateLiveQuizSession(c, sess)
	if err != nil {
		return &CreateLiveQuizSessionResponse{}, err
	}

	lqs := &LiveQuizSession{
		Session: *sess,
		Code:    code,
	}

	er := s.Repository.CreateLiveQuizSessionCache(c, lqs)
	if er != nil {
		return &CreateLiveQuizSessionResponse{}, er
	}

	return &CreateLiveQuizSessionResponse{
		ID:     lqs.ID,
		QuizID: lqs.QuizID,
		Code:   lqs.Code,
	}, nil
}

func (s *service) GetLiveQuizSessions(ctx context.Context) ([]LiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqses, err := s.Repository.GetLiveQuizSessions(c)
	if err != nil {
		return []LiveQuizSessionResponse{}, err
	}

	var lqsesRes []LiveQuizSessionResponse
	for _, lqs := range lqses {
		lqsesRes = append(lqsesRes, LiveQuizSessionResponse{
			ID:     lqs.ID,
			HostID: lqs.HostID,
			QuizID: lqs.QuizID,
			Status: lqs.Status,
			Code:   lqs.Code,
		})
	}

	return lqsesRes, nil
}

func (s *service) GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*LiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqs, err := s.Repository.GetLiveQuizSessionByID(c, id)
	if err != nil {
		return &LiveQuizSessionResponse{}, err
	}

	return &LiveQuizSessionResponse{
		ID:     lqs.ID,
		HostID: lqs.HostID,
		QuizID: lqs.QuizID,
		Status: lqs.Status,
		Code:   lqs.Code,
	}, nil
}

func (s *service) GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*LiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqs, err := s.Repository.GetLiveQuizSessionByQuizID(c, quizID)
	if err != nil {
		return &LiveQuizSessionResponse{}, err
	}

	return &LiveQuizSessionResponse{
		ID:     lqs.ID,
		HostID: lqs.HostID,
		QuizID: lqs.QuizID,
		Status: lqs.Status,
		Code:   lqs.Code,
	}, nil
}

func (s *service) GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqs, err := s.Repository.GetLiveQuizSessionCache(c, code)
	if err != nil {
		return &Cache{}, err
	}

	return lqs, nil
}

func (s *service) UpdateLiveQuizSession(ctx context.Context, req *UpdateLiveQuizSessionRequest, id uuid.UUID) (*LiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqs, err := s.Repository.GetLiveQuizSessionByID(c, id)
	if err != nil {
		return &LiveQuizSessionResponse{}, err
	}

	if req.Status != "" {
		lqs.Status = req.Status
	}
	if req.ExemptedQuesIDs != nil {
		lqs.ExemptedQuestionIDs = req.ExemptedQuesIDs
	}

	lqs, err = s.Repository.UpdateLiveQuizSession(c, lqs, id)
	if err != nil {
		return &LiveQuizSessionResponse{}, err
	}

	return &LiveQuizSessionResponse{
		ID:     lqs.ID,
		HostID: lqs.HostID,
		QuizID: lqs.QuizID,
		Status: lqs.Status,
		Code:   lqs.Code,
	}, nil
}

func (s *service) DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.Repository.GetLiveQuizSessionByID(c, id)
	if err != nil {
		return err
	}

	err = s.Repository.DeleteLiveQuizSession(c, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetQuestionByQuizIDAndOrder(ctx context.Context, quizID uuid.UUID, order int) (*q.Question, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	ques, err := s.quiz.GetQuestionByQuizID(c, quizID, order)
	if err != nil {
		return &q.Question{}, err
	}

	return ques, nil
}

func (s *service) CreateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.CreateLiveQuizSessionCache(c, lqs)
	if err != nil {
		return err
	}

	return nil
}

// ---------- Participant related service methods ---------- //
func (s *service) CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participant, err := s.Repository.CreateParticipant(c, participant)
	if err != nil {
		return &Participant{}, err
	}

	return participant, nil
}

func (s *service) GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) (*ParticipantsResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participants, err := s.Repository.GetParticipantsByLiveQuizSessionID(c, lqsID)
	if err != nil {
		return &ParticipantsResponse{}, err
	}

	return &ParticipantsResponse{
		Participants: participants,
	}, nil
}

func (s *service) DoesParticipantExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	exists, err := s.Repository.DoesParticipantExists(c, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *service) UpdateParticipantStatus(ctx context.Context, userID uuid.UUID, lqsID uuid.UUID, status string) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.UpdateParticipantStatus(c, userID, lqsID, status)
	if err != nil {
		return &Participant{}, err
	}

	return p, nil
}
