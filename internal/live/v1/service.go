package v1

import (
	"context"
	"encoding/json"
	"log"
	"sort"
	"time"

	u "github.com/Live-Quiz-Project/Backend/internal/user/v1"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
)

type service struct {
	Repository
	timeout  time.Duration
	userRepo u.Repository
}

func NewService(r Repository, uRepo u.Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(3) * time.Second,
		userRepo:   uRepo,
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
			ID:                  lqs.ID,
			HostID:              lqs.HostID,
			QuizID:              lqs.QuizID,
			Status:              lqs.Status,
			ExemptedQuestionIDs: lqs.ExemptedQuestionIDs,
			CreatedAt:           lqs.CreatedAt,
			UpdatedAt:           lqs.UpdatedAt,
			DeletedAt:           lqs.DeletedAt,
		},
	}, nil
}

func (s *service) GetLiveQuizSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]SessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sessions, err := s.Repository.GetLiveQuizSessionsByUserID(c, userID)
	if err != nil {
		return nil, err
	}

	var res []SessionResponse

	for _, lqs := range sessions {
		res = append(res, SessionResponse{
			Session: Session{
				ID:                  lqs.ID,
				HostID:              lqs.HostID,
				QuizID:              lqs.QuizID,
				Status:              lqs.Status,
				ExemptedQuestionIDs: lqs.ExemptedQuestionIDs,
				CreatedAt:           lqs.CreatedAt,
				UpdatedAt:           lqs.UpdatedAt,
				DeletedAt:           lqs.DeletedAt,
			},
		})
	}

	return res, nil
}

// ---------- Live Quiz Session related service methods ---------- //
func (s *service) CreateLiveQuizSession(ctx context.Context, quizID uuid.UUID, id uuid.UUID, code string, hostID uuid.UUID) (*CreateLiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess := &Session{
		ID:                  id,
		HostID:              hostID,
		QuizID:              quizID,
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

func (s *service) GetLiveQuizSessions(ctx context.Context, hub *Hub) ([]LiveQuizSessionResponse, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	lqses, err := s.Repository.GetLiveQuizSessions(c)
	if err != nil {
		return []LiveQuizSessionResponse{}, err
	}

	var lqsesRes []LiveQuizSessionResponse
	for _, lqs := range lqses {
		for _, lq := range hub.LiveQuizSessions {
			if lqs.ID == lq.ID {
				lqsesRes = append(lqsesRes, LiveQuizSessionResponse{
					ID:     lq.ID,
					HostID: lq.HostID,
					QuizID: lq.QuizID,
					Status: lq.Status,
					Code:   lq.Code,
				})
			}
		}
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

	err := s.Repository.CreateCache(c, code, cache)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	cache, err := s.Repository.GetCache(c, code)
	if err != nil {
		return &Cache{}, err
	}

	var mod *Cache
	err = json.Unmarshal([]byte(cache), &mod)
	if err != nil {
		if err.Error() == "redis: nil" {
			return &Cache{}, nil
		}
		return &Cache{}, err
	}

	return mod, nil
}

func (s *service) UpdateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.UpdateCache(c, code, cache)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) FlushLiveQuizSessionCache(ctx context.Context, code string) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.Repository.FlushCache(c, code)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) DoesLiveQuizSessionCacheExist(ctx context.Context, code string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	exists, err := s.Repository.DoesCacheExist(c, code)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *service) FlushAllLiveQuizSessionRelatedCache(ctx context.Context, code string) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	keys, err := s.Repository.ScanCache(c, code+"*")
	if err != nil {
		return err
	}

	for _, k := range keys {
		err := s.Repository.FlushCache(c, k)
		if err != nil {
			return err
		}
	}

	return nil
}

// ---------- Participant related service methods ---------- //
func (s *service) CreateParticipant(ctx context.Context, p *Participant) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.CreateParticipant(c, p)
	if err != nil {
		return &Participant{}, err
	}

	return p, nil
}

func (s *service) GetParticipantByID(ctx context.Context, id uuid.UUID) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.GetParticipantByID(c, id)
	if err != nil {
		return &Participant{}, err
	}

	return p, nil
}

func (s *service) GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	participants, err := s.Repository.GetParticipantsByLiveQuizSessionID(c, lqsID)
	if err != nil {
		return []Participant{}, err
	}

	return participants, nil
}

func (s *service) DoesParticipantExist(ctx context.Context, id uuid.UUID) (bool, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	exists, err := s.Repository.DoesParticipantExist(c, id)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *service) UpdateParticipant(ctx context.Context, p *Participant) (*Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.UpdateParticipant(c, p)
	if err != nil {
		return &Participant{}, err
	}

	return p, nil
}

// ---------- Response related service methods ---------- //
func (s *service) CreateResponse(ctx context.Context, code string, qid string, pid string, response any) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.CreateCache(c, code+qid+pid, response); err != nil {
		return err
	}

	return nil
}

func (s *service) GetResponses(ctx context.Context, code string, qid string) ([]any, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	keys, err := s.Repository.ScanCache(c, code+qid+"*")
	if err != nil {
		return nil, err
	}

	res := make([]any, 0)
	for _, k := range keys {
		response, err := s.Repository.GetCache(c, k)
		if err != nil {
			return nil, err
		}
		if response == "" {
			return nil, nil
		}
		var r any
		err = json.Unmarshal([]byte(response), &r)
		if err != nil {
			return nil, err
		}
		res = append(res, r)
	}

	return res, nil
}

func (s *service) GetResponse(ctx context.Context, code string, qid string, pid string) (any, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	response, err := s.Repository.GetCache(c, code+qid+pid)
	if err != nil {
		return nil, err
	}
	if response == "" {
		return nil, nil
	}
	var res any
	err = json.Unmarshal([]byte(response), &res)
	if err != nil {
		return nil, err
	}
	log.Println(res)

	return res, nil
}

func (s *service) UpdateResponse(ctx context.Context, code string, qid string, pid string, response any) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.UpdateCache(c, code+qid+pid, response); err != nil {
		return err
	}

	return nil
}

func (s *service) FlushResponse(ctx context.Context, code string, qid string, pid string) error {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if err := s.Repository.FlushCache(c, code+qid+pid); err != nil {
		return err
	}

	return nil
}

func (s *service) DoesResponseExist(ctx context.Context, code string, qid string, pid string) (bool, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	exists, err := s.Repository.DoesCacheExist(c, code+qid+pid)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *service) CountResponses(ctx context.Context, code string, qid string) (int, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	count, err := s.Repository.ScanCache(c, code+qid+"*")
	if err != nil {
		return 0, err
	}
	log.Println("Count: ", count)

	return len(count), nil
}

func (s *service) SaveResponse(ctx context.Context, response *Response) (*Response, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	response, err := s.Repository.CreateResponse(c, response)
	if err != nil {
		return &Response{}, err
	}

	return response, nil
}

// ---------- Leaderboard related service methods ---------- //
func (s *service) GetLeaderboard(ctx context.Context, lqsID uuid.UUID) ([]Participant, error) {
	c, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	p, err := s.Repository.GetParticipantsByLiveQuizSessionID(c, lqsID)
	if err != nil {
		return []Participant{}, err
	}

	sort.Sort(ByMarks(p))

	return p, nil
}
