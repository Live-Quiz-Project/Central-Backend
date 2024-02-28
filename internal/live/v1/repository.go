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

// ---------- Session ------------- //
func (r *repository) GetLiveQuizSessionBySessionID(ctx context.Context, id uuid.UUID) (*Session, error) {
	var lqs Session
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return &lqs, nil
}

func (r *repository) GetLiveQuizSessionsByUserID(ctx context.Context, id uuid.UUID) ([]Session, error) {
	var lqs []Session
	res := r.db.WithContext(ctx).Where("host_id = ?", id).Order("created_at DESC").Find(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return lqs, nil
}

// ---------- Live Quiz Session related repository methods ---------- //
func (r *repository) CreateLiveQuizSession(ctx context.Context, lqs *Session) (*Session, error) {
	res := r.db.WithContext(ctx).Create(lqs)
	if res.Error != nil {
		return &Session{}, res.Error
	}

	return lqs, nil
}

func (r *repository) GetLiveQuizSessions(ctx context.Context) ([]Session, error) {
	var lqses []Session
	res := r.db.WithContext(ctx).Find(&lqses)
	if res.Error != nil {
		return nil, res.Error
	}
	return lqses, nil
}

func (r *repository) GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*Session, error) {
	var lqs Session
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return &lqs, nil
}

func (r *repository) GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*Session, error) {
	var lqs Session
	res := r.db.WithContext(ctx).Where("quiz_id = ?", quizID).First(&lqs)
	if res.Error != nil {
		return nil, res.Error
	}
	return &lqs, nil
}

func (r *repository) GetLiveQuizSessionByCode(ctx context.Context, code string) (*Session, error) {
	var lqs Session
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

func (r *repository) UpdateLiveQuizSession(ctx context.Context, lqs *Session, id uuid.UUID) (*Session, error) {
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

func (r *repository) CreateCache(ctx context.Context, key string, value any) error {
	val, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}

	status := r.cache.Set(ctx, key, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) GetCache(ctx context.Context, key string) (string, error) {
	val, err := r.cache.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			return "", nil
		}
		return "", err
	}

	return val, nil
}

func (r *repository) UpdateCache(ctx context.Context, key string, value any) error {
	val, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	if val == nil {
		return nil
	}

	status := r.cache.Set(ctx, key, val, time.Duration(60*60*5)*time.Second)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) FlushCache(ctx context.Context, key string) error {
	status := r.cache.Del(ctx, key)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *repository) DoesCacheExist(ctx context.Context, key string) (bool, error) {
	val, err := r.cache.Exists(ctx, key).Result()
	if err != nil {
		if err.Error() == redis.Nil.Error() {
			return false, nil
		}
		return false, err
	}

	return val == 1, nil
}

func (r *repository) ScanCache(ctx context.Context, pattern string) ([]string, error) {
	var cursor uint64
	var keys []string
	var err error

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			var scanResult []string
			scanResult, cursor, err = r.cache.Scan(ctx, cursor, pattern, 10).Result()
			if err != nil {
				return nil, err
			}
			keys = append(keys, scanResult...)
			if cursor == 0 {
				return keys, nil
			}
		}
	}
}

// ---------- Participant related repository methods ---------- //
func (r *repository) CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error) {
	res := r.db.WithContext(ctx).Create(participant)
	if res.Error != nil {
		return &Participant{}, res.Error
	}

	return participant, nil
}

func (r *repository) GetParticipantByID(ctx context.Context, id uuid.UUID) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("id = ?", id).First(&participant)
	if res.Error != nil {
		return nil, res.Error
	}

	return &participant, nil
}

func (r *repository) GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error) {
	p := make([]Participant, 0)
	res := r.db.WithContext(ctx).Where("live_quiz_session_id = ? AND status = ?", lqsID, util.Joined).Find(&p)
	if res.Error != nil {
		return nil, res.Error
	}

	return p, nil
}

func (r *repository) GetParticipantByUserIDAndLiveQuizSessionID(ctx context.Context, uid *uuid.UUID, lqsID uuid.UUID) (*Participant, error) {
	var participant Participant
	res := r.db.WithContext(ctx).Where("user_id = ? AND live_quiz_session_id = ?", uid, lqsID).First(&participant)
	if res.Error != nil {
		return nil, res.Error
	}

	return &participant, nil
}

func (r *repository) DoesParticipantExist(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	res := r.db.WithContext(ctx).Model(&Participant{}).Where("id = ?", id).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (r *repository) UpdateParticipant(ctx context.Context, p *Participant) (*Participant, error) {
	res := r.db.WithContext(ctx).Where("id = ?", p.ID).Updates(p)
	if res.Error != nil {
		return nil, res.Error
	}

	return p, nil
}

// ---------- Response related repository methods ---------- //
func (r *repository) CreateResponse(ctx context.Context, ansRes *Response) (*Response, error) {
	res := r.db.WithContext(ctx).Create(ansRes)
	if res.Error != nil {
		return &Response{}, res.Error
	}

	return ansRes, nil
}
