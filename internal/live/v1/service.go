package v1

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"strings"
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

func (s *service) GetAnswersResponseForHost(ctx context.Context, qType string, answers []any, answerCounts map[string]int) (any, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var res []any

	switch qType {
	case util.Choice, util.TrueFalse:
		cAns := make([]ChoiceAnswer, len(answers))
		for i, a := range answers {
			v, ok := a.(map[string]any)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			id, ok := v["id"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			content, ok := v["content"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			color, ok := v["color"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			m, ok := v["mark"].(float64)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			mark := int(m)
			isCorrect, ok := v["is_correct"].(bool)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			count := answerCounts[id]

			cAns[i] = ChoiceAnswer{
				ID:      id,
				Content: content,
				Color:   color,
				Mark:    mark,
				Correct: isCorrect,
				Count:   count,
			}
		}
		res = append(res, cAns)
	case util.FillBlank, util.Paragraph:
		tAns := make([]TextAnswer, len(answers))
		for i, a := range answers {
			v, ok := a.(map[string]any)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			id, ok := v["id"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			content, ok := v["content"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			caseSensitive, ok := v["case_sensitive"].(bool)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			m, ok := v["mark"].(float64)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			mark := int(m)

			tAns[i] = TextAnswer{
				ID:            id,
				Content:       content,
				CaseSensitive: caseSensitive,
				Mark:          mark,
			}
		}
		res = append(res, tAns)
	case util.Matching:
		mAns := make([]MatchingAnswer, len(answers))
		for i, a := range answers {
			v, ok := a.(map[string]any)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			prompt, ok := v["prompt_id"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			option, ok := v["option_id"].(string)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			m, ok := v["mark"].(float64)
			if !ok {
				return nil, errors.New("invalid type assertion")
			}
			mark := int(m)

			mAns[i] = MatchingAnswer{
				PromptID: prompt,
				OptionID: option,
				Mark:     mark,
			}
		}
		res = append(res, mAns)
	case util.Pool:
		var pAns []any
		for _, an := range answers {
			ans, ok := an.([]any)
			if !ok {
				return nil, errors.New("invalid type assertion1")
			}

			sqType := util.Paragraph
			for _, answ := range ans {
				sqType, ok = answ.(map[string]any)["type"].(string)
				if !ok {
					return nil, errors.New("invalid type assertion2")
				}
			}

			switch sqType {
			case util.Choice, util.TrueFalse:
				cAns := make([]ChoiceAnswer, len(ans))
				for i, a := range ans {
					v, ok := a.(map[string]any)
					if !ok {
						return nil, errors.New("invalid type assertion3")
					}
					id, ok := v["id"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion4")
					}
					content, ok := v["content"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion5")
					}
					color, ok := v["color"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion6")
					}
					m, ok := v["mark"].(float64)
					if !ok {
						return nil, errors.New("invalid type assertion7")
					}
					mark := int(m)
					isCorrect, ok := v["is_correct"].(bool)
					if !ok {
						return nil, errors.New("invalid type assertion8")
					}
					count := answerCounts[id]

					cAns[i] = ChoiceAnswer{
						ID:      id,
						Content: content,
						Color:   color,
						Mark:    mark,
						Correct: isCorrect,
						Count:   count,
					}
				}
				pAns = append(pAns, cAns)
			case util.FillBlank, util.Paragraph:
				tAns := make([]TextAnswer, len(ans))
				for i, a := range ans {
					v, ok := a.(map[string]any)
					if !ok {
						return nil, errors.New("invalid type assertion9")
					}
					id, ok := v["id"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion10")
					}
					content, ok := v["content"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion11")
					}
					caseSensitive, ok := v["case_sensitive"].(bool)
					if !ok {
						return nil, errors.New("invalid type assertion12")
					}
					m, ok := v["mark"].(float64)
					if !ok {
						return nil, errors.New("invalid type assertion13")
					}
					mark := int(m)

					tAns[i] = TextAnswer{
						ID:            id,
						Content:       content,
						CaseSensitive: caseSensitive,
						Mark:          mark,
					}
				}
				pAns = append(pAns, tAns)
			case util.Matching:
				mAns := make([]MatchingAnswer, len(ans))
				for i, a := range ans {
					v, ok := a.(map[string]any)
					if !ok {
						return nil, errors.New("invalid type assertion14")
					}
					prompt, ok := v["prompt_id"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion15")
					}
					option, ok := v["option_id"].(string)
					if !ok {
						return nil, errors.New("invalid type assertion16")
					}
					m, ok := v["mark"].(float64)
					if !ok {
						return nil, errors.New("invalid type assertion17")
					}
					mark := int(m)

					mAns[i] = MatchingAnswer{
						PromptID: prompt,
						OptionID: option,
						Mark:     mark,
					}
				}
				pAns = append(pAns, mAns)
			}

			res = pAns
		}
	}

	return res, nil
}

func (s *service) CalculateChoice(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (ChoiceAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]ChoiceAnswer, 0)

	for _, o := range options {
		oID, ok := o.(map[string]any)["id"].(string)
		if !ok {
			return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
		}
		for _, a := range answers {
			aID, ok := a.(map[string]any)["id"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
			}
			aContent, ok := a.(map[string]any)["content"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
			}
			aColor, ok := a.(map[string]any)["color"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
			}
			aIsCorrect, ok := a.(map[string]any)["is_correct"].(bool)
			if !ok {
				return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return ChoiceAnswerResponse{}, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			if oID == aID {
				marks += mark
				if status == util.Answering {
					mark += 1000
					res = append(res, ChoiceAnswer{
						ID:      aID,
						Content: aContent,
						Color:   aColor,
					})
				}
				if status == util.RevealingAnswer {
					res = append(res, ChoiceAnswer{
						ID:      aID,
						Content: aContent,
						Color:   aColor,
						Mark:    mark,
						Correct: aIsCorrect,
					})
				}
			}
		}
	}

	var mRes *int
	if status == util.Answering {
		mRes = nil
	}
	if status == util.RevealingAnswer {
		mRes = &marks
	}

	return ChoiceAnswerResponse{
		Answers: res,
		Marks:   mRes,
		Time:    int(time),
	}, nil
}

func (s *service) CalculateAndSaveChoiceResponse(ctx context.Context, options []any, answers []any, answerCounts map[string]int, time float64, timeLimit float64, timeFactor float64, response *Response) (ChoiceAnswerResponse, map[string]int, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]ChoiceAnswer, 0)
	stringifyOptions := make([]string, len(options))

	for i, o := range options {
		oID, ok := o.(map[string]any)["id"].(string)
		if !ok {
			return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
		}
		answerCounts[oID] += 1
		stringifyOptions[i] = string(oID)
		for _, a := range answers {
			aID, ok := a.(map[string]any)["id"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
			}
			aContent, ok := a.(map[string]any)["content"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
			}
			aColor, ok := a.(map[string]any)["color"].(string)
			if !ok {
				return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
			}
			aIsCorrect, ok := a.(map[string]any)["is_correct"].(bool)
			if !ok {
				return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return ChoiceAnswerResponse{}, nil, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			if oID == aID {
				marks += mark
				res = append(res, ChoiceAnswer{
					ID:      aID,
					Content: aContent,
					Color:   aColor,
					Mark:    mark,
					Correct: aIsCorrect,
				})
			}
		}
	}

	p, err := s.GetParticipantByID(context.Background(), response.ParticipantID)
	if err != nil {
		return ChoiceAnswerResponse{}, nil, err
	}
	p.Marks += marks
	_, err = s.UpdateParticipant(context.Background(), p)
	if err != nil {
		return ChoiceAnswerResponse{}, nil, err
	}

	stringifyAnswer := strings.Join(stringifyOptions, util.AnswerSplitter)
	if _, err := s.SaveResponse(context.Background(), &Response{
		ID:                response.ID,
		LiveQuizSessionID: response.LiveQuizSessionID,
		QuestionID:        response.QuestionID,
		ParticipantID:     response.ParticipantID,
		Type:              response.Type,
		TimeTaken:         int(time),
		Answer:            stringifyAnswer,
	}); err != nil {
		return ChoiceAnswerResponse{}, nil, err
	}

	return ChoiceAnswerResponse{
		Answers: res,
		Marks:   &marks,
		Time:    int(time),
	}, answerCounts, nil
}

func (s *service) CalculateFillBlank(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (TextAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]TextAnswer, 0)

	for _, o := range options {
		oID, ok := o.(map[string]any)["id"].(string)
		if !ok {
			return TextAnswerResponse{}, errors.New("invalid type assertion")
		}
		oContent, ok := o.(map[string]any)["content"].(string)
		if !ok {
			return TextAnswerResponse{}, errors.New("invalid type assertion")
		}
		for _, a := range answers {
			aID, ok := a.(map[string]any)["id"].(string)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aContent, ok := a.(map[string]any)["content"].(string)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aCaseSensitive, ok := a.(map[string]any)["case_sensitive"].(bool)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			isCorrect := (aCaseSensitive && aContent == oContent) || (!aCaseSensitive && strings.EqualFold(aContent, oContent))
			if oID == aID {
				m := 0
				if isCorrect {
					marks += mark
					m = mark
				}
				if status == util.Answering {
					res = append(res, TextAnswer{
						ID:            aID,
						CaseSensitive: aCaseSensitive,
						Content:       oContent,
					})
				}
				if status == util.RevealingAnswer {
					res = append(res, TextAnswer{
						ID:            aID,
						CaseSensitive: aCaseSensitive,
						Content:       oContent,
						Answer:        aContent,
						Correct:       isCorrect,
						Mark:          m,
					})
				}
			}
		}
	}

	var mRes *int
	if status == util.Answering {
		mRes = nil
	}
	if status == util.RevealingAnswer {
		mRes = &marks
	}

	return TextAnswerResponse{
		Answers: res,
		Marks:   mRes,
		Time:    int(time),
	}, nil
}

func (s *service) CalculateAndSaveFillBlankResponse(ctx context.Context, options []any, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (TextAnswerResponse, error) {
	_, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]TextAnswer, 0)
	stringifyOptions := make([]string, len(options))

	for i, o := range options {
		oID, ok := o.(map[string]any)["id"].(string)
		if !ok {
			return TextAnswerResponse{}, errors.New("invalid type assertion")
		}
		oContent, ok := o.(map[string]any)["content"].(string)
		if !ok {
			return TextAnswerResponse{}, errors.New("invalid type assertion")
		}
		stringifyOptions[i] = string(oContent)
		for _, a := range answers {
			aID, ok := a.(map[string]any)["id"].(string)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aContent, ok := a.(map[string]any)["content"].(string)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aCaseSensitive, ok := a.(map[string]any)["case_sensitive"].(bool)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return TextAnswerResponse{}, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			isCorrect := (aCaseSensitive && aContent == oContent) || (!aCaseSensitive && strings.EqualFold(aContent, oContent))
			if oID == aID {
				m := 0
				if isCorrect {
					marks += mark
					m = mark
				}
				res = append(res, TextAnswer{
					ID:            aID,
					CaseSensitive: aCaseSensitive,
					Content:       oContent,
					Answer:        aContent,
					Correct:       isCorrect,
					Mark:          m,
				})
			}
		}
	}

	p, err := s.GetParticipantByID(context.Background(), response.ParticipantID)
	if err != nil {
		return TextAnswerResponse{}, err
	}
	p.Marks += marks
	_, err = s.UpdateParticipant(context.Background(), p)
	if err != nil {
		return TextAnswerResponse{}, err
	}

	stringifyAnswer := strings.Join(stringifyOptions, util.AnswerSplitter)
	if _, err := s.SaveResponse(context.Background(), &Response{
		ID:                response.ID,
		LiveQuizSessionID: response.LiveQuizSessionID,
		QuestionID:        response.QuestionID,
		ParticipantID:     response.ParticipantID,
		Type:              response.Type,
		TimeTaken:         int(time),
		Answer:            stringifyAnswer,
	}); err != nil {
		return TextAnswerResponse{}, err
	}

	return TextAnswerResponse{
		Answers: res,
		Marks:   &marks,
		Time:    int(time),
	}, nil
}

func (s *service) CalculateParagraph(ctx context.Context, status string, content string, answers []any, time float64, timeLimit float64, timeFactor float64) (any, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]TextAnswer, 0)

	var r any
	r = content
	l := len(answers)
	if l > 0 {
		answer, ok := answers[0].(map[string]any)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aID, ok := answer["id"].(string)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aContent, ok := answer["content"].(string)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aCaseSensitive, ok := answer["case_sensitive"].(bool)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aM, ok := (answer["mark"].(float64))
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aMark := int(aM)

		var timeBonus float64
		if aMark > 0 {
			timeBonus = timeLimit - factoredTime
		}
		mark := (int(aM + timeBonus))

		isCorrect := (aCaseSensitive && aContent == content) || (!aCaseSensitive && strings.EqualFold(aContent, content))
		if isCorrect {
			marks += mark
		}
		if status == util.RevealingAnswer {
			res = append(res, TextAnswer{
				ID:            aID,
				CaseSensitive: aCaseSensitive,
				Content:       content,
				Answer:        aContent,
				Correct:       isCorrect,
				Mark:          mark,
			})
			r = TextAnswerResponse{
				Answers: res,
				Marks:   &marks,
				Time:    int(time),
			}
		}
	}

	var finalRes any
	if status == util.Answering {
		finalRes = ParagraphAnswerResponse{
			Answer: content,
			Marks:  nil,
			Time:   int(time),
		}
	}
	if status == util.RevealingAnswer {
		finalRes = r
	}

	return finalRes, nil
}

func (s *service) CalculateAndSaveParagraphResponse(ctx context.Context, content string, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (any, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]TextAnswer, 0)

	var r any
	r = content
	l := len(answers)
	if l > 0 {
		answer, ok := answers[0].(map[string]any)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aID, ok := answer["id"].(string)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aContent, ok := answer["content"].(string)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aCaseSensitive, ok := answer["case_sensitive"].(bool)
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aM, ok := (answer["mark"].(float64))
		if !ok {
			return nil, errors.New("invalid type assertion")
		}
		aMark := int(aM)

		var timeBonus float64
		if aMark > 0 {
			timeBonus = timeLimit - factoredTime
		}
		mark := (int(aM + timeBonus))

		isCorrect := (aCaseSensitive && aContent == content) || (!aCaseSensitive && strings.EqualFold(aContent, content))
		if isCorrect {
			marks += mark
		}
		res = append(res, TextAnswer{
			ID:            aID,
			CaseSensitive: aCaseSensitive,
			Content:       content,
			Answer:        aContent,
			Correct:       isCorrect,
			Mark:          mark,
		})

		r = TextAnswerResponse{
			Answers: res,
			Marks:   &marks,
			Time:    int(time),
		}

		p, err := s.GetParticipantByID(context.Background(), response.ParticipantID)
		if err != nil {
			return TextAnswerResponse{}, err
		}
		p.Marks += marks
		_, err = s.UpdateParticipant(context.Background(), p)
		if err != nil {
			return TextAnswerResponse{}, err
		}
	}

	if _, err := s.SaveResponse(context.Background(), &Response{
		ID:                response.ID,
		LiveQuizSessionID: response.LiveQuizSessionID,
		QuestionID:        response.QuestionID,
		ParticipantID:     response.ParticipantID,
		Type:              response.Type,
		TimeTaken:         int(time),
		Answer:            content,
	}); err != nil {
		return TextAnswerResponse{}, err
	}

	return r, nil
}

func (s *service) CalculateMatching(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (MatchingAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]MatchingAnswer, 0)

	for _, o := range options {
		oPrompt, ok := o.(map[string]any)["prompt"].(string)
		if !ok {
			return MatchingAnswerResponse{}, errors.New("invalid type assertion")
		}
		oOption, ok := o.(map[string]any)["option"].(string)
		if !ok {
			return MatchingAnswerResponse{}, errors.New("invalid type assertion")
		}
		for _, a := range answers {
			aPrompt, ok := a.(map[string]any)["prompt_id"].(string)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aOption, ok := a.(map[string]any)["option_id"].(string)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			isCorrect := oPrompt == aPrompt && oOption == aOption
			if oPrompt == aPrompt {
				m := 0
				if isCorrect {
					marks += mark
					m = mark
				}
				if status == util.Answering {
					res = append(res, MatchingAnswer{
						PromptID: aPrompt,
						OptionID: oOption,
					})
				}
				if status == util.RevealingAnswer {
					res = append(res, MatchingAnswer{
						PromptID: aPrompt,
						OptionID: oOption,
						Correct:  isCorrect,
						Mark:     m,
					})
				}
			}
		}
	}

	var mRes *int
	if status == util.Answering {
		mRes = nil
	}
	if status == util.RevealingAnswer {
		mRes = &marks
	}

	return MatchingAnswerResponse{
		Answers: res,
		Marks:   mRes,
		Time:    int(time),
	}, nil
}

func (s *service) CalculateAndSaveMatchingResponse(ctx context.Context, options []any, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (MatchingAnswerResponse, error) {
	_, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	factoredTime := 0.0
	factoredTime = time * (timeFactor / 10)
	marks := 0
	res := make([]MatchingAnswer, 0)
	stringifyOptions := make([]string, len(options))

	for i, o := range options {
		oPrompt, ok := o.(map[string]any)["prompt"].(string)
		if !ok {
			return MatchingAnswerResponse{}, errors.New("invalid type assertion")
		}
		oOption, ok := o.(map[string]any)["option"].(string)
		if !ok {
			return MatchingAnswerResponse{}, errors.New("invalid type assertion")
		}
		stringifyOptions[i] = string(oPrompt + ":" + oOption)
		for _, a := range answers {
			aPrompt, ok := a.(map[string]any)["prompt_id"].(string)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aOption, ok := a.(map[string]any)["option_id"].(string)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aM, ok := a.(map[string]any)["mark"].(float64)
			if !ok {
				return MatchingAnswerResponse{}, errors.New("invalid type assertion")
			}
			aMark := int(aM)

			var timeBonus float64
			if aMark > 0 {
				timeBonus = timeLimit - factoredTime
			}
			mark := (int(aM + timeBonus))

			isCorrect := oPrompt == aPrompt && oOption == aOption
			if oPrompt == aPrompt {
				m := 0
				if isCorrect {
					marks += mark
					m = mark
				}
				res = append(res, MatchingAnswer{
					PromptID: aPrompt,
					OptionID: oOption,
					Correct:  isCorrect,
					Mark:     m,
				})
			}
		}
	}

	p, err := s.GetParticipantByID(context.Background(), response.ParticipantID)
	if err != nil {
		return MatchingAnswerResponse{}, err
	}
	p.Marks += marks
	_, err = s.UpdateParticipant(context.Background(), p)
	if err != nil {
		return MatchingAnswerResponse{}, err
	}

	stringifyAnswer := strings.Join(stringifyOptions, util.AnswerSplitter)
	if _, err := s.SaveResponse(context.Background(), &Response{
		ID:                response.ID,
		LiveQuizSessionID: response.LiveQuizSessionID,
		QuestionID:        response.QuestionID,
		ParticipantID:     response.ParticipantID,
		Type:              response.Type,
		TimeTaken:         int(time),
		Answer:            stringifyAnswer,
	}); err != nil {
		return MatchingAnswerResponse{}, err
	}

	return MatchingAnswerResponse{
		Answers: res,
		Marks:   &marks,
		Time:    int(time),
	}, nil
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
