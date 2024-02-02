package v1

import (
	"context"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(3) * time.Second,
	}
}

// ---------- Session related service methods ---------- //
func (s *service) GetLiveQuizSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqs, err := s.Repository.GetLiveQuizSessionBySessionID(c, sessionID)
	if err != nil {
		return &SessionResponse{}, err
	}

	return &SessionResponse{
		Session: Session{
			ID: lqs.ID,
			HostID : lqs.HostID,
			QuizID: lqs.QuizID,
			Status : lqs.Status,
			ExemptedQuestionIDs : lqs.ExemptedQuestionIDs,
			CreatedAt :lqs.CreatedAt,
			UpdatedAt : lqs.UpdatedAt,
			DeletedAt: lqs.DeletedAt,
		},
	}, nil
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

func (s *service) CreateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.CreateLiveQuizSessionCache(c, code, cache)
	if err != nil {
		return err
	}

	return nil
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

func (s *service) UpdateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.UpdateLiveQuizSessionCache(c, code, cache)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FlushLiveQuizSessionCache(ctx context.Context, code string) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.FlushLiveQuizSessionCache(c, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) CreateLiveQuizSessionResponseCache(ctx context.Context, code string, response any) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.CreateLiveQuizSessionResponseCache(c, code, response)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetLiveQuizSessionResponseCache(ctx context.Context, code string) (any, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.Repository.GetLiveQuizSessionResponseCache(c, code)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateLiveQuizSessionResponseCache(ctx context.Context, code string, response any) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.UpdateLiveQuizSessionResponseCache(c, code, response)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FlushLiveQuizSessionResponseCache(ctx context.Context, code string) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.FlushLiveQuizSessionResponseCache(c, code)
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

func (s *service) DoesParticipantExists(ctx context.Context, uid uuid.UUID, lqsID uuid.UUID) (bool, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	exists, err := s.Repository.DoesParticipantExists(c, uid, lqsID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *service) UpdateParticipantStatus(ctx context.Context, uid uuid.UUID, lqsID uuid.UUID, status string) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.UpdateParticipantStatus(c, uid, lqsID, status)
	if err != nil {
		return &Participant{}, err
	}

	return p, nil
}

func (s *service) UnregisterParticipants(ctx context.Context, lqsID uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.UnregisterParticipants(c, lqsID)
	if err != nil {
		return err
	}

	return nil
}

// ---------- Response related service methods ---------- //
// Choice response related methods
func (s *service) CreateChoiceResponse(ctx context.Context, req *CreateChoiceResponseRequest, uid uuid.UUID) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cr := &ChoiceResponse{
		ID:             uuid.New(),
		ParticipantID:  uid,
		OptionChoiceID: req.OptionChoiceID,
	}

	_, err := s.Repository.CreateChoiceResponse(c, cr)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetChoiceResponsesByParticipantID(ctx context.Context, participantID uuid.UUID) ([]ChoiceResponseResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	choiceResponses, err := s.Repository.GetChoiceResponsesByParticipantID(c, participantID)
	if err != nil {
		return []ChoiceResponseResponse{}, err
	}

	var choiceResponsesRes []ChoiceResponseResponse
	for _, cr := range choiceResponses {
		choiceResponsesRes = append(choiceResponsesRes, ChoiceResponseResponse{ChoiceResponse: cr})
	}

	return choiceResponsesRes, nil
}

func (s *service) GetChoiceResponsesByQuizID(ctx context.Context, quizID uuid.UUID) ([]ChoiceResponseResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	choiceResponses, err := s.Repository.GetChoiceResponsesByQuizID(c, quizID)
	if err != nil {
		return []ChoiceResponseResponse{}, err
	}

	var choiceResponsesRes []ChoiceResponseResponse
	for _, cr := range choiceResponses {
		choiceResponsesRes = append(choiceResponsesRes, ChoiceResponseResponse{ChoiceResponse: cr})
	}

	return choiceResponsesRes, nil
}

func (s *service) GetChoiceResponsesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceResponseResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	choiceResponses, err := s.Repository.GetChoiceResponsesByQuestionID(c, questionID)
	if err != nil {
		return []ChoiceResponseResponse{}, err
	}

	var choiceResponsesRes []ChoiceResponseResponse
	for _, cr := range choiceResponses {
		choiceResponsesRes = append(choiceResponsesRes, ChoiceResponseResponse{ChoiceResponse: cr})
	}

	return choiceResponsesRes, nil
}

func (s *service) GetChoiceResponseByParticipantIDAndQuestionID(ctx context.Context, participantID uuid.UUID, questionID uuid.UUID) (*ChoiceResponseResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cr, err := s.Repository.GetChoiceResponseByParticipantIDAndQuestionID(c, participantID, questionID)
	if err != nil {
		return &ChoiceResponseResponse{}, err
	}

	return &ChoiceResponseResponse{ChoiceResponse: *cr}, nil
}
