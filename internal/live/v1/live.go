package v1

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ---------- Live Quiz Session related models ---------- //
type Session struct {
	ID                  uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	HostID              uuid.UUID  `json:"host_id" gorm:"column:host_id;type:uuid;not null"`
	QuizID              uuid.UUID  `json:"quiz_id" gorm:"column:quiz_id;type:uuid;not null"`
	Status              string     `json:"status" gorm:"column:status;not null"`
	ExemptedQuestionIDs *string    `json:"exempted_question_ids" gorm:"column:exempted_question_ids"`
	CreatedAt           time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamptz;not null"`
	DeletedAt           *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamptz"`
}

func (Session) TableName() string {
	return "live_quiz_session"
}

type LiveQuizSession struct {
	Session
	Code    string                `json:"code"`
	Clients map[uuid.UUID]*Client `json:"clients"`
}

type Cache struct {
	LiveQuizSessionID uuid.UUID                 `json:"live_quiz_session_id"`
	QuizID            uuid.UUID                 `json:"quiz_id"`
	HostID            uuid.UUID                 `json:"host_id"`
	QuizTitle         string                    `json:"quiz_title"`
	QuestionCount     int                       `json:"question_count"`
	CurrentQuestion   int                       `json:"current_question"`
	Questions         []any                     `json:"questions"`
	Answers           []any                     `json:"answers"`
	AnswerCounts      map[string]map[string]int `json:"answer_counts"`
	Status            string                    `json:"status"`
	Config            Configurations            `json:"config"`
	Locked            bool                      `json:"locked"`
	Interrupted       bool                      `json:"interrupted"`
	Orders            []int                     `json:"orders"`
	ResponseCount     int                       `json:"response_count"`
	ParticipantCount  int                       `json:"participant_count"`
}

type SessionResponse struct {
	Session
}

type Configurations struct {
	ShuffleConfig     ShuffleConfigurations     `json:"shuffle"`
	ParticipantConfig ParticipantConfigurations `json:"participant"`
	LeaderboardConfig LeaderboardConfigurations `json:"leaderboard"`
	OptionConfig      OptionConfigurations      `json:"option"`
}

type ShuffleConfigurations struct {
	Question bool `json:"question"`
	Option   bool `json:"option"`
}

type ParticipantConfigurations struct {
	Reanswer bool `json:"reanswer"`
}

type LeaderboardConfigurations struct {
	DuringQuestions bool `json:"during"`
	AfterQuestions  bool `json:"after"`
}

type OptionConfigurations struct {
	Colorless         bool `json:"colorless"`
	ShowCorrectAnswer bool `json:"show_correct_answer"`
}

// ---------- Participant related models ---------- //
type Participant struct {
	ID                uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	UserID            *uuid.UUID `json:"user_id" gorm:"column:user_id;type:uuid"`
	LiveQuizSessionID uuid.UUID  `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null"`
	Status            string     `json:"status" gorm:"column:status;type:text;not null"`
	Name              string     `json:"display_name" gorm:"column:name;type:text"`
	Emoji             string     `json:"display_emoji" gorm:"column:emoji;type:text"`
	Color             string     `json:"display_color" gorm:"column:color;type:text"`
	Marks             int        `json:"marks" gorm:"column:marks;type:int"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamptz;not null"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamptz"`
}

func (Participant) TableName() string {
	return "participant"
}

// ---------- Response related models ---------- //
type Response struct {
	ID                uuid.UUID  `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	LiveQuizSessionID uuid.UUID  `json:"live_quiz_session_id" gorm:"column:live_quiz_session_id;type:uuid;not null"`
	ParticipantID     uuid.UUID  `json:"participant_id" gorm:"column:participant_id;type:uuid;not null"`
	QuestionID        uuid.UUID  `json:"question_id" gorm:"column:question_id;type:uuid;not null"`
	Type              string     `json:"type" gorm:"column:type;type:text;not null"`
	Answer            any        `json:"answer" gorm:"column:answer;type:text;not null"`
	TimeTaken         int        `json:"time" gorm:"column:use_time;type:int;not null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at;type:timestamptz;not null"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:timestamptz"`
}

func (Response) TableName() string {
	return "answer_response"
}

type Repository interface {
	GetLiveQuizSessionBySessionID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetLiveQuizSessionsByUserID(ctx context.Context, id uuid.UUID) ([]Session, error)

	// ---------- Live Quiz Session related repository methods ---------- //
	CreateLiveQuizSession(ctx context.Context, lqs *Session) (*Session, error)
	GetLiveQuizSessions(ctx context.Context) ([]Session, error)
	GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*Session, error)
	GetLiveQuizSessionByCode(ctx context.Context, code string) (*Session, error)
	UpdateLiveQuizSession(ctx context.Context, lqs *Session, id uuid.UUID) (*Session, error)
	EndLiveQuizSession(ctx context.Context, id uuid.UUID) error
	DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error

	CreateCache(ctx context.Context, key string, value any) error
	GetCache(ctx context.Context, key string) (string, error)
	UpdateCache(ctx context.Context, key string, value any) error
	FlushCache(ctx context.Context, key string) error
	DoesCacheExist(ctx context.Context, key string) (bool, error)
	ScanCache(ctx context.Context, pattern string) ([]string, error)

	// ---------- Participant related repository methods ---------- //
	CreateParticipant(ctx context.Context, participant *Participant) (*Participant, error)
	GetParticipantByID(ctx context.Context, id uuid.UUID) (*Participant, error)
	GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error)
	GetParticipantByUserIDAndLiveQuizSessionID(ctx context.Context, uid *uuid.UUID, lqsID uuid.UUID) (*Participant, error)
	DoesParticipantExist(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateParticipant(ctx context.Context, participant *Participant) (*Participant, error)

	// ---------- Response related repository methods ---------- //
	CreateResponse(ctx context.Context, ansRes *Response) (*Response, error)
	// Choice response related repository methods
	// CreateChoiceResponse(ctx context.Context, r *ChoiceResponse) (*ChoiceResponse, error)
	// GetChoiceResponsesByParticipantID(ctx context.Context, participantID uuid.UUID) ([]ChoiceResponse, error)
	// GetChoiceResponsesByQuizID(ctx context.Context, quizID uuid.UUID) ([]ChoiceResponse, error)
	// GetChoiceResponsesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]ChoiceResponse, error)
	// GetChoiceResponseByParticipantIDAndQuestionID(ctx context.Context, participantID uuid.UUID, questionID uuid.UUID) (*ChoiceResponse, error)
}

// ---------- Live Quiz Session related structs ---------- //
type LiveQuizSessionResponse struct {
	ID     uuid.UUID `json:"id"`
	HostID uuid.UUID `json:"host_id"`
	QuizID uuid.UUID `json:"quiz_id"`
	Code   string    `json:"code"`
	Status string    `json:"status"`
}
type CreateLiveQuizSessionRequest struct {
	QuizID uuid.UUID      `json:"quiz_id"`
	Config Configurations `json:"config"`
}
type CreateLiveQuizSessionResponse struct {
	ID     uuid.UUID `json:"id"`
	QuizID uuid.UUID `json:"quiz_id"`
	Code   string    `json:"code"`
}
type UpdateLiveQuizSessionRequest struct {
	Status          string  `json:"status"`
	ExemptedQuesIDs *string `json:"exempted_question_ids"`
}

type JoinedMessage struct {
	Code    string    `json:"code"`
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Emoji   string    `json:"emoji"`
	Color   string    `json:"color"`
	IsHost  bool      `json:"is_host"`
	Answers any       `json:"answers"`
	Marks   int       `json:"marks"`
}

type CheckLiveQuizSessionAvailabilityResponse struct {
	ID              uuid.UUID `json:"id"`
	QuizID          uuid.UUID `json:"quiz_id"`
	QuizTitle       string    `json:"quiz_title"`
	Code            string    `json:"code"`
	QuestionCount   int       `json:"question_count"`
	CurrentQuestion int       `json:"current_question"`
	Status          string    `json:"status"`
}

type CountDownPayload struct {
	TimeLeft        float64 `json:"time_left"`
	CurrentQuestion int     `json:"current_question"`
	Status          string  `json:"status"`
}

type ByMarks []Participant

func (a ByMarks) Len() int      { return len(a) }
func (a ByMarks) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByMarks) Less(i, j int) bool {
	if a[i].Marks == a[j].Marks {
		return a[i].Name < a[j].Name
	}
	return a[i].Marks > a[j].Marks
}

type ChoiceAnswer struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Color   string `json:"color"`
	Mark    int    `json:"mark"`
	Correct bool   `json:"correct"`
	Count   int    `json:"count"`
}
type ChoiceAnswerResponse struct {
	Answers []ChoiceAnswer `json:"answers"`
	Marks   *int           `json:"marks"`
	Time    int            `json:"time"`
}
type TextAnswer struct {
	ID            string `json:"id"`
	CaseSensitive bool   `json:"caseSensitive"`
	Content       string `json:"content"`
	Answer        string `json:"answer"`
	Correct       bool   `json:"correct"`
	Mark          int    `json:"mark"`
}
type TextAnswerResponse struct {
	Answers []TextAnswer `json:"answers"`
	Marks   *int         `json:"marks"`
	Time    int          `json:"time"`
}
type ParagraphAnswerResponse struct {
	Answer string `json:"answer"`
	Marks  *int   `json:"marks"`
	Time   int    `json:"time"`
}
type MatchingAnswer struct {
	PromptID string `json:"prompt"`
	OptionID string `json:"option"`
	Correct  bool   `json:"correct"`
	Mark     int    `json:"mark"`
}
type MatchingAnswerResponse struct {
	Answers []MatchingAnswer `json:"answers"`
	Marks   *int             `json:"marks"`
	Time    int              `json:"time"`
}

type PoolAnswer struct {
	ID      string `json:"qid"`
	Type    string `json:"type"`
	Content any    `json:"content"`
}
type PoolAnswerResponse struct {
	Answers map[string]PoolAnswer `json:"answers"`
	Marks   int                   `json:"marks"`
	Time    int                   `json:"time"`
}

type AnswerPayload struct {
	Answers       any       `json:"answers"`
	ParticipantID uuid.UUID `json:"participant_id"`
	TotalMarks    int       `json:"marks"`
}

type Service interface {
	GetLiveQuizSessionBySessionID(ctx context.Context, sessionID uuid.UUID) (*SessionResponse, error)
	GetLiveQuizSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]SessionResponse, error)

	// ---------- Live Quiz Session related service methods ---------- //
	CreateLiveQuizSession(ctx context.Context, quizID uuid.UUID, id uuid.UUID, code string, hostID uuid.UUID) (*CreateLiveQuizSessionResponse, error)
	GetLiveQuizSessions(ctx context.Context, hub *Hub) ([]LiveQuizSessionResponse, error)
	GetLiveQuizSessionByID(ctx context.Context, id uuid.UUID) (*LiveQuizSessionResponse, error)
	GetLiveQuizSessionByQuizID(ctx context.Context, quizID uuid.UUID) (*LiveQuizSessionResponse, error)
	UpdateLiveQuizSession(ctx context.Context, req *UpdateLiveQuizSessionRequest, id uuid.UUID) (*LiveQuizSessionResponse, error)
	DeleteLiveQuizSession(ctx context.Context, id uuid.UUID) error

	CreateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error
	GetLiveQuizSessionCache(ctx context.Context, code string) (*Cache, error)
	UpdateLiveQuizSessionCache(ctx context.Context, code string, cache *Cache) error
	FlushLiveQuizSessionCache(ctx context.Context, code string) error
	DoesLiveQuizSessionCacheExist(ctx context.Context, code string) (bool, error)

	FlushAllLiveQuizSessionRelatedCache(ctx context.Context, code string) error

	// ---------- Participant related service methods ---------- //
	CreateParticipant(ctx context.Context, p *Participant) (*Participant, error)
	GetParticipantByID(ctx context.Context, id uuid.UUID) (*Participant, error)
	GetParticipantsByLiveQuizSessionID(ctx context.Context, lqsID uuid.UUID) ([]Participant, error)
	DoesParticipantExist(ctx context.Context, id uuid.UUID) (bool, error)
	UpdateParticipant(ctx context.Context, p *Participant) (*Participant, error)

	// ---------- Response related service methods ---------- //
	CreateResponse(ctx context.Context, code string, qid string, pid string, response any) error
	GetResponses(ctx context.Context, code string, qid string) ([]any, error)
	GetResponse(ctx context.Context, code string, qid string, pid string) (any, error)
	UpdateResponse(ctx context.Context, code string, qid string, pid string, response any) error
	FlushResponse(ctx context.Context, code string, qid string, pid string) error
	DoesResponseExist(ctx context.Context, code string, qid string, pid string) (bool, error)
	CountResponses(ctx context.Context, code string, qid string) (int, error)
	SaveResponse(ctx context.Context, response *Response) (*Response, error)

	// ---------- Calculation related service methods ---------- //
	GetAnswersResponseForHost(ctx context.Context, qid string, qType string, answers []any, answerCounts map[string]map[string]int) (any, error)
	CalculateChoice(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (ChoiceAnswerResponse, error)
	CalculateAndSaveChoiceResponse(ctx context.Context, options []any, answers []any, answerCounts map[string]int, time float64, timeLimit float64, timeFactor float64, response *Response) (ChoiceAnswerResponse, map[string]int, error)
	CalculateFillBlank(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (TextAnswerResponse, error)
	CalculateAndSaveFillBlankResponse(ctx context.Context, options []any, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (TextAnswerResponse, error)
	CalculateParagraph(ctx context.Context, status string, content string, answers []any, time float64, timeLimit float64, timeFactor float64) (any, error)
	CalculateAndSaveParagraphResponse(ctx context.Context, content string, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (any, error)
	CalculateMatching(ctx context.Context, status string, options []any, answers []any, time float64, timeLimit float64, timeFactor float64) (MatchingAnswerResponse, error)
	CalculateAndSaveMatchingResponse(ctx context.Context, options []any, answers []any, time float64, timeLimit float64, timeFactor float64, response *Response) (MatchingAnswerResponse, error)

	// ---------- Leaderboard related service methods ---------- //
	GetLeaderboard(ctx context.Context, lqsID uuid.UUID) ([]Participant, error)
}
