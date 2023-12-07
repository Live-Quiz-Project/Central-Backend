package v1

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type repository struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewRepository(db *gorm.DB, cache *redis.Client) Repository {
	return &repository{
		db:    db,
		cache: cache,
	}
}

// ---------- Live Quiz Session related repository methods ---------- //
func (r *repository) CreateLiveQuizSession(ctx context.Context, lqs *Session) (*Session, error) {
	res := r.db.WithContext(ctx).Create(lqs)
	if res.Error != nil {
		return &Session{}, res.Error
	}

	return lqs, nil
}

func (r *repository) GetLiveQuizSessions(ctx context.Context) ([]LiveQuizSession, error) {
	var lqses []LiveQuizSession
	res := r.db.WithContext(ctx).Find(&lqses)
	if res.Error != nil {
		return nil, res.Error
	}
	return lqses, nil
}

func (r *repository) GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*LiveQuizSession, error) {
	var lqs LiveQuizSession
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return &lqs, nil
}

func (r *repository) GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*LiveQuizSession, error) {
	var lqs LiveQuizSession
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).First(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return &lqs, nil
}

func (r *repository) GetLiveQuizSessionByCode(ctx context.Context, code string) (*LiveQuizSession, error) {
	var lqs LiveQuizSession
	val, err := r.cache.Get(ctx, code).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(val), &lqs)
	if err != nil {
		return nil, err
	}
	return &lqs, nil
}

func (r *repository) UpdateLiveQuizSession(ctx context.Context, lqs *LiveQuizSession, id uuid.UUID) (*LiveQuizSession, error) {
	res := r.db.WithContext(ctx).Where("id = ?", id).Updates(lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return lqs, nil
}

func (r *repository) EndLiveQuizSession(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&LiveQuizSession{}).Where("id = ?", id).Update("status", util.Ended)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repository) DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&LiveQuizSession{}, id)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *repository) CreateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error {
	val, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error) {
	var lqs Cache
	val, err := r.cache.Get(ctx, code).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &lqs)
	if err != nil {
		return nil, err
	}

	return &lqs, nil
}

func (r *repository) UpdateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error {
	val, err := json.Marshal(&cache)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) FlushLiveQuizSessionCache(ctx context.Context, code string) error {
	status := r.cache.Del(ctx, code)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) CreateLiveQuizSessionResponseCache(ctx context.Context, code string, response any) error {
	val, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code+"RESP", val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) GetLiveQuizSessionResponseCache(ctx context.Context, code string) (any, error) {
	var response any
	val, err := r.cache.Get(ctx, code+"RESP").Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(val), &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *repository) UpdateLiveQuizSessionResponseCache(ctx context.Context, code string, response any) error {
	val, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code+"RESP", val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) FlushLiveQuizSessionResponseCache(ctx context.Context, code string) error {
	status := r.cache.Del(ctx, code+"RESP")
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

// ---------- Participant related repository methods ---------- //
func (r *repository) CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error) {
	res := r.db.WithContext(ctx).Create(participant)
	if res.Error != nil {
		return &Participant{}, res.Error
	}

	return participant, nil
}

func (r *repository) GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error) {
	var participants []Participant
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ? AND status = ?", lqsID, util.Joined).Find(&participants)
	if res.Error != nil {
		return nil, res.Error
	}
	return participants, nil
}

func (r *repository) GetParticipantByUserIDAndLiveQuizSessionID(ctx context.Context, uid uuid.UUID, lqsID uuid.UUID) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("user_id = ? AND live_quiz_session_id = ?", uid, lqsID).First(&participant)
	if res.Error != nil {
		return nil, res.Error
	}

	return &participant, nil
}

func (r *repository) DoesParticipantExists(ctx context.Context, uid uuid.UUID, lqsID uuid.UUID) (bool, error) {
	var count int64
	res := r.db.WithContext(ctx).Model(&Participant{}).Where("user_id = ? AND live_quiz_session_id = ?", uid, lqsID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (r *repository) UpdateParticipantStatus(ctx context.Context, uid uuid.UUID, quizID uuid.UUID, status string) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("user_id = ? AND live_quiz_session_id = ?", uid, quizID).First(&participant)
	if res.Error != nil {
		return nil, res.Error
	}

	participant.Status = status
	res = r.db.WithContext(ctx).Save(&participant)
	if res.Error != nil {
		return nil, res.Error
	}

	return &participant, nil
}

func (r *repository) UnregisterParticipants(ctx context.Context, lqsID uuid.UUID) error {
	res := r.db.WithContext(ctx).Model(&Participant{}).Where("live_quiz_session_id = ?", lqsID).Update("status", util.Left)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// ---------- Response related repository methods ---------- //
// Choice response related methods
func (r *repository) CreateChoiceResponse(ctx context.Context, cr *ChoiceResponse) (*ChoiceResponse, error) {
	res := r.db.WithContext(ctx).Create(cr)
	if res.Error != nil {
		return &ChoiceResponse{}, res.Error
	}

	return cr, nil
}

func (r *repository) GetChoiceResponsesByParticipantID(ctx context.Context, participantID uuid.UUID) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("participant_id = ?", participantID).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}

	return responses, nil
}

func (r *repository) GetChoiceResponsesByQuizID(ctx context.Context, quizID uuid.UUID) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}

	return responses, nil
}

func (r *repository) GetChoiceResponsesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceResponse, error) {
	var responses []ChoiceResponse
	res := r.db.WithContext(ctx).Where("question_id = ?", questionID).Find(&responses)
	if res.Error != nil {
		return nil, res.Error
	}

	return responses, nil
}

func (r *repository) GetChoiceResponseByParticipantIDAndQuestionID(ctx context.Context, participantID uuid.UUID, questionID uuid.UUID) (*ChoiceResponse, error) {
	var response ChoiceResponse
	res := r.db.WithContext(ctx).Where("participant_id = ? AND question_id = ?", participantID, questionID).First(&response)
	if res.Error != nil {
		return nil, res.Error
	}

	return &response, nil
}
