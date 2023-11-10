package v1

import (
	"context"
	"encoding/json"
	"strconv"
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

func (r *repository) CreateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error {
	val, err := json.Marshal(&Cache{
		ID:     lqs.ID,
		QuizID: lqs.QuizID,
		Status: lqs.Status,
	})
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, lqs.Code, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) CreateLiveQuizSessionResponseCache(ctx context.Context, code string, order int, response any) error {
	val, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code+"RESP"+strconv.Itoa(order), val, time.Duration(60*60*5)*time.Second)
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

func (r *repository) UpdateLiveQuizSessionCache(ctx context.Context, lqs *LiveQuizSession) error {
	val, err := json.Marshal(&Cache{
		ID:     lqs.ID,
		QuizID: lqs.QuizID,
		Status: lqs.Status,
	})
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, lqs.Code, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) UpdateLiveQuizSessionResponseCache(ctx context.Context, code string, order int, response any) error {
	val, err := json.Marshal(&response)
	if err != nil {
		return err
	}

	status := r.cache.Set(ctx, code+"RESP"+strconv.Itoa(order), val, time.Duration(60*60*5)*time.Second)
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
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ?", lqsID).Find(&participants)
	if res.Error != nil {
		return nil, res.Error
	}
	return participants, nil
}

func (r *repository) DoesParticipantExists(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int64
	res := r.db.WithContext(ctx).Model(&Participant{}).Where("user_id = ?", userID).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (r *repository) UpdateParticipantStatus(ctx context.Context, userID uuid.UUID, quizID uuid.UUID, status string) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("user_id = ? AND live_quiz_session_id = ?", userID, quizID).First(&participant)
	if res.Error != nil {
		return nil, res.Error
	}
	return &participant, nil
}
