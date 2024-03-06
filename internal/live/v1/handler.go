package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	q "github.com/Live-Quiz-Project/Backend/internal/quiz/v1"
	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
	Service
	quizService q.Service
}

func NewHandler(h *Hub, lServ Service, qServ q.Service) *Handler {
	return &Handler{
		hub:         h,
		Service:     lServ,
		quizService: qServ,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool {
		origin := req.Header.Get("Origin")
		allowOriginsEnv := os.Getenv("ALLOW_ORIGINS")
		allowOrigins := strings.Split(allowOriginsEnv, ",")
		for _, allowedOrigin := range allowOrigins {
			if allowedOrigin == origin {
				return true
			}
		}
		return false
	},
}

func (h *Handler) CreateLiveQuizSession(c *gin.Context) {
	var req CreateLiveQuizSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	latestQuizID, err := h.quizService.GetLatestQuizVersionByID(c, req.QuizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	lqsID := uuid.New()

	var code string
	codes := make([]string, 0)
	for _, s := range h.hub.LiveQuizSessions {
		if s.QuizID == req.QuizID {
			lqsID = s.ID
		}
		codes = append(codes, s.Code)
	}

	hostID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	quizTitle, err := h.quizService.GetQuizHistoryByID(c, *latestQuizID, hostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		code = util.CodeGenerator(codes)
		lqs, err := h.Service.CreateLiveQuizSession(c, *latestQuizID, lqsID, code, hostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		questions, err := h.quizService.GetQuestionsByQuizIDForLQS(c, *latestQuizID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		count := len(questions)
		orders := make([]int, count)
		if req.Config.ShuffleConfig.Question {
			orders = util.ShuffleNumbers(count)
		} else {
			for i := 0; i < count; i++ {
				orders[i] = i + 1
			}
		}

		answers, err := h.quizService.GetAnswersByQuizIDForLQS(c, *latestQuizID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		err = h.Service.CreateLiveQuizSessionCache(context.Background(), code, &Cache{
			LiveQuizSessionID: lqsID,
			HostID:            hostID,
			QuizTitle:         quizTitle.Title,
			QuizID:            lqs.QuizID,
			QuestionCount:     count,
			CurrentQuestion:   0,
			Questions:         questions,
			Answers:           answers,
			AnswerCounts:      make(map[string]map[string]int),
			Status:            util.Idle,
			Config:            req.Config,
			Locked:            false,
			Interrupted:       false,
			Orders:            orders,
			ResponseCount:     0,
			ParticipantCount:  0,
		})
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		h.hub.LiveQuizSessions[lqsID] = &LiveQuizSession{
			Session: Session{
				ID:                  lqs.ID,
				HostID:              hostID,
				QuizID:              lqs.QuizID,
				Status:              util.Ongoing,
				ExemptedQuestionIDs: nil,
			},
			Code:    lqs.Code,
			Clients: make(map[uuid.UUID]*Client),
		}
		c.JSON(http.StatusOK, &CreateLiveQuizSessionResponse{
			ID:     lqsID,
			QuizID: req.QuizID,
			Code:   code,
		})
		return
	}

	c.JSON(http.StatusOK, &CreateLiveQuizSessionResponse{
		ID:     h.hub.LiveQuizSessions[lqsID].ID,
		QuizID: h.hub.LiveQuizSessions[lqsID].QuizID,
		Code:   h.hub.LiveQuizSessions[lqsID].Code,
	})
}

func (h *Handler) GetLiveQuizSessions(c *gin.Context) {
	lqs, err := h.Service.GetLiveQuizSessions(c, h.hub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lqs)
}

func (h *Handler) EndLiveQuizSession(c *gin.Context) {
	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.EndLQS,
			Payload: nil,
		},
		LiveQuizSessionID: uuid.Nil,
		ClientID:          uuid.Nil,
		UserID:            nil,
	}

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	code := c.Param("code")
	var lqsID uuid.UUID
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code && s.HostID == userID {
			lqsID = s.ID
			break
		}
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	if userID != h.hub.LiveQuizSessions[lqsID].HostID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only the host can end the session"})
		return
	}

	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
		h.hub.Inject <- &Message{
			Content: Content{
				Type:    util.EndLQS,
				Payload: nil,
			},
			LiveQuizSessionID: lqsID,
			ClientID:          cl.ID,
			UserID:            cl.UserID,
		}
	}

	err = h.Service.FlushAllLiveQuizSessionRelatedCache(c, h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	delete(h.hub.LiveQuizSessions, lqsID)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully ended the session"})
}

func (h *Handler) CheckLiveQuizSessionAvailability(c *gin.Context) {
	code := c.Param("code")

	mod, err := h.Service.GetLiveQuizSessionCache(c, code)
	if err != nil && err.Error() != "redis: nil" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	count := 0
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code {
			for _, cl := range s.Clients {
				if cl.Status == util.Joined {
					count++
				}
			}
		}
	}

	if count > 60 || mod.Locked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is full or locked"})
		return
	}

	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code {
			c.JSON(http.StatusOK, &CheckLiveQuizSessionAvailabilityResponse{
				ID:              s.ID,
				QuizID:          s.QuizID,
				Code:            s.Code,
				QuizTitle:       mod.QuizTitle,
				QuestionCount:   mod.QuestionCount,
				CurrentQuestion: mod.CurrentQuestion,
				Status:          mod.Status,
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
}

func (h *Handler) JoinLiveQuizSession(c *gin.Context) {
	var err error
	code := c.Param("code")

	uid := c.Query("uid")
	var userID *uuid.UUID
	userID = nil
	if uid != "" {
		parsedUID, err := uuid.Parse(uid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = &parsedUID
	}

	var isHost bool
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code && userID != nil && s.HostID == *userID {
			isHost = true
			break
		}
	}

	uname := c.Query("name")
	emoji := c.Query("emoji")
	color := c.Query("color")

	var lqsID uuid.UUID
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code {
			lqsID = s.ID
			break
		}
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	conn, e := upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
		return
	}

	pid := c.Query("pid")
	participantID := uuid.New()
	if pid != "" {
		participantID, err = uuid.Parse(pid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
			return
		}
	}

	p := &Participant{
		ID:                participantID,
		UserID:            userID,
		LiveQuizSessionID: lqsID,
		Status:            util.Joined,
		Marks:             0,
		Name:              uname,
		Emoji:             emoji,
		Color:             color,
	}

	var pCount int
	if !isHost {
		exists, eErr := h.Service.DoesParticipantExist(c, p.ID)
		if eErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": eErr.Error()})
			return
		}
		if exists {
			p, err := h.Service.GetParticipantByID(c, p.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			p.Name = uname
			p.Emoji = emoji
			p.Color = color
			p.Status = util.Joined
			if _, err = h.Service.UpdateParticipant(c, p); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			_, pErr := h.Service.CreateParticipant(c, p)
			if pErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": pErr.Error()})
				return
			}
		}
	}

	participants, err := h.Service.GetParticipantsByLiveQuizSessionID(c, lqsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	pCount = len(participants)
	mod, err := h.Service.GetLiveQuizSessionCache(c, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	mod.ParticipantCount = pCount
	if err := h.Service.UpdateLiveQuizSessionCache(c, code, mod); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cl := &Client{
		Conn:              conn,
		Message:           make(chan *Message, 10),
		ID:                p.ID,
		UserID:            p.UserID,
		DisplayName:       p.Name,
		DisplayEmoji:      p.Emoji,
		DisplayColor:      p.Color,
		IsHost:            isHost,
		LiveQuizSessionID: lqsID,
		Status:            util.Joined,
	}
	h.hub.Register <- cl

	var answers any
	if !isHost && mod.CurrentQuestion > 0 && (mod.Status == util.Answering || mod.Status == util.RevealingAnswer) {
		res, err := h.Service.GetResponse(c, code, mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string), p.ID.String())
		if err != nil {
			log.Printf("Error occured here @399: %v", err)
			return
		}
		if res != nil {
			time, ok := res.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @708: %v", err)
				return
			}
			qAns, ok := mod.Answers[mod.Orders[mod.CurrentQuestion-1]-1].([]any)
			if !ok {
				log.Printf("Error occured @792: Type assertion failed")
				return
			}
			qType, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["type"].(string)
			if !ok {
				log.Printf("Error occured @708: %v", err)
				return
			}
			qTimeLimit, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["time_limit"].(float64)
			if !ok {
				log.Printf("Error occured @708: %v", err)
				return
			}
			qHaveTimeFactor, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["have_time_factor"].(bool)
			if !ok {
				log.Printf("Error occured @708: %v", err)
				return
			}
			qTimeFactor, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["time_factor"].(float64)
			if !ok {
				log.Printf("Error occured @708: %v", err)
				return
			}
			if !qHaveTimeFactor {
				qTimeFactor = 0
			}

			switch qType {
			case util.Choice, util.TrueFalse:
				opt, ok := res.(map[string]any)["options"].([]any)
				if !ok {
					log.Printf("Error occured @456: Type assertion failed")
					return
				}

				answers, err = h.Service.CalculateChoice(c, mod.Status, opt, qAns, time, qTimeLimit, qTimeFactor)
				if err != nil {
					log.Printf("Error occured @456: %v", err)
					return
				}
			case util.FillBlank:
				opt, ok := res.(map[string]any)["options"].([]any)
				if !ok {
					log.Printf("Error occured @456: Type assertion failed")
					return
				}

				answers, err = h.Service.CalculateFillBlank(c, mod.Status, opt, qAns, time, qTimeLimit, qTimeFactor)
				if err != nil {
					log.Printf("Error occured @456: %v", err)
					return
				}
			case util.Paragraph:
				opt, ok := res.(map[string]any)["options"].(string)
				if !ok {
					log.Printf("Error occured @456: Type assertion failed")
					return
				}

				answers, err = h.Service.CalculateParagraph(c, mod.Status, opt, qAns, time, qTimeLimit, qTimeFactor)
				if err != nil {
					log.Printf("Error occured @456: %v", err)
					return
				}
			case util.Matching:
				opt, ok := res.(map[string]any)["options"].([]any)
				if !ok {
					log.Printf("Error occured @456: Type assertion failed")
					return
				}

				answers, err = h.Service.CalculateMatching(c, mod.Status, opt, qAns, time, qTimeLimit, qTimeFactor)
				if err != nil {
					log.Printf("Error occured @456: %v", err)
					return
				}
			case util.Pool:
				opt, ok := res.(map[string]any)["options"].(map[string]any)
				if !ok {
					log.Printf("Error occured @1: Type assertion failed")
					return
				}

				ansRes := make(map[string]PoolAnswer, 0)
				var marksRes int
				var timeRes int

				for i, o := range opt {
					I, err := strconv.Atoi(i)
					if err != nil {
						log.Printf("Error occured @3: %v", err)
						return
					}
					sqID, ok := o.(map[string]any)["qid"].(string)
					if !ok {
						log.Printf("Error occured @10101: Type assertion failed")
						return
					}
					sqType, ok := o.(map[string]any)["type"].(string)
					if !ok {
						log.Printf("Error occured @2: Type assertion failed")
						return
					}

					switch sqType {
					case util.Choice, util.TrueFalse:
						opt, ok := o.(map[string]any)["content"].([]any)
						if !ok {
							log.Printf("Error occured @3: Type assertion failed")
							return
						}
						a := qAns[I].([]any)

						r, err := h.Service.CalculateChoice(c, mod.Status, opt, a, time, qTimeLimit, qTimeFactor)
						if err != nil {
							log.Printf("Error occured @4: %v", err)
							return
						}

						ansRes[i] = PoolAnswer{
							ID:      sqID,
							Type:    sqType,
							Content: r.Answers,
						}
						marksRes += *r.Marks
						timeRes = r.Time
					case util.FillBlank:
						opt, ok := o.(map[string]any)["content"].([]any)
						if !ok {
							log.Printf("Error occured @5: Type assertion failed")
							return
						}
						a := qAns[I].([]any)

						r, err := h.Service.CalculateFillBlank(c, mod.Status, opt, a, time, qTimeLimit, qTimeFactor)
						if err != nil {
							log.Printf("Error occured @6: %v", err)
							return
						}

						ansRes[i] = PoolAnswer{
							ID:      sqID,
							Type:    sqType,
							Content: r.Answers,
						}
						marksRes += *r.Marks
						timeRes = r.Time
					case util.Paragraph:
						opt, ok := o.(map[string]any)["content"]
						if !ok {
							log.Printf("Error occured @7: Type assertion failed")
							return
						}
						a := qAns[I].([]any)

						content := ""
						switch opt := opt.(type) {
						case string:
							content = opt
						case nil:
						}

						r, err := h.Service.CalculateParagraph(c, mod.Status, content, a, time, qTimeLimit, qTimeFactor)
						if err != nil {
							log.Printf("Error occured @8: %v", err)
							return
						}

						switch r := r.(type) {
						case string:
							ansRes[i] = PoolAnswer{
								ID:      sqID,
								Type:    sqType,
								Content: r,
							}
							marksRes += 0
							timeRes = 0
						case ParagraphAnswerResponse:
							ansRes[i] = PoolAnswer{
								ID:      sqID,
								Type:    sqType,
								Content: r.Answer,
							}
							marksRes += *r.Marks
							timeRes = r.Time
						case TextAnswerResponse:
							ansRes[i] = PoolAnswer{
								ID:      sqID,
								Type:    sqType,
								Content: r.Answers,
							}
							marksRes += *r.Marks
							timeRes = r.Time
						default:
						}
					case util.Matching:
						opt, ok := o.(map[string]any)["content"].([]any)
						if !ok {
							log.Printf("Error occured @9: Type assertion failed")
							return
						}
						a := qAns[I].([]any)

						r, err := h.Service.CalculateMatching(c, mod.Status, opt, a, time, qTimeLimit, qTimeFactor)
						if err != nil {
							log.Printf("Error occured @10: %v", err)
							return
						}

						ansRes[i] = PoolAnswer{
							ID:      sqID,
							Type:    sqType,
							Content: r.Answers,
						}
						marksRes += *r.Marks
						timeRes = r.Time
					}
				}

				answers = PoolAnswerResponse{
					Answers: ansRes,
					Marks:   marksRes,
					Time:    timeRes,
				}

				val, err := json.MarshalIndent(answers, "", "  ")
				if err != nil {
					log.Printf("Error occured @11: %v", err)
					return
				}
				log.Printf("Value: %v", string(val))
			}
		}
	}
	if isHost && mod.CurrentQuestion > 0 && mod.Status == util.RevealingAnswer {
		qid, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string)
		if !ok {
			log.Printf("Error occured @708: %v", err)
			return
		}
		qAns, ok := mod.Answers[mod.Orders[mod.CurrentQuestion-1]-1].([]any)
		if !ok {
			log.Printf("Error occured @792: Type assertion failed")
			return
		}
		qType, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["type"].(string)
		if !ok {
			log.Printf("Error occured @708: %v", err)
			return
		}

		answers, err = h.Service.GetAnswersResponseForHost(context.Background(), qid, qType, qAns, mod.AnswerCounts)
		if err != nil {
			log.Printf("Error occured @699: %v", err)
			return
		}
	}

	rank := -1
	if !isHost {
		rank, err = h.Service.GetRank(c, lqsID, p.ID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}
	}

	go cl.writeMessage()
	h.hub.Converse <- &Message{
		Content: Content{
			Type: util.JoinLQS,
			Payload: JoinedMessage{
				Code:    code,
				ID:      cl.ID,
				Name:    cl.DisplayName,
				Emoji:   cl.DisplayEmoji,
				Color:   cl.DisplayColor,
				Marks:   p.Marks,
				IsHost:  cl.IsHost,
				Answers: answers,
				Rank:    rank,
			},
		},
		LiveQuizSessionID: lqsID,
		ClientID:          cl.ID,
		UserID:            cl.UserID,
	}
	cl.readMessage(h)
}

func (h *Handler) KickParticipant(c *Client, payload any) {
	pid := payload.(map[string]any)["id"].(string)
	kickedID, err := uuid.Parse(pid)
	if err != nil {
		log.Println("Error occured: ", err)
		return
	}

	participants, err := h.Service.GetParticipantsByLiveQuizSessionID(context.Background(), c.LiveQuizSessionID)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	pCount := len(participants)

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.ParticipantCount = pCount
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Inject <- &Message{
		Content: Content{
			Type:    util.KickParticipant,
			Payload: nil,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          kickedID,
		UserID:            c.UserID,
	}
}

func (h *Handler) ToggleLiveQuizSessionLock(c *Client) {
	var code string
	for _, s := range h.hub.LiveQuizSessions {
		if s.ID == c.LiveQuizSessionID {
			code = s.Code
		}
	}

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	mod.Locked = !mod.Locked

	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.ToggleLock,
			Payload: mod.Locked,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) InterruptCountdown(c *gin.Context) {
	code := c.Param("code")

	mod, err := h.Service.GetLiveQuizSessionCache(c, code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	mod.Interrupted = true

	err = h.Service.UpdateLiveQuizSessionCache(c, code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully interrupted the countdown"})
}

func (h *Handler) GetParticipants(c *Client) {
	var err error
	var p []Participant

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	if mod.Status == util.Idle {
		p, err = h.Service.GetParticipantsByLiveQuizSessionID(context.Background(), c.LiveQuizSessionID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}
	} else {
		p, err = h.Service.GetLeaderboard(context.Background(), c.LiveQuizSessionID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}
	}

	if !c.IsHost && (((mod.Status == util.Questioning || mod.Status == util.Answering) && !mod.Config.LeaderboardConfig.DuringQuestions) || (mod.Status == util.RevealingAnswer && !mod.Config.LeaderboardConfig.AfterQuestions)) {
		p = []Participant{}
	}

	h.hub.Inject <- &Message{
		Content: Content{
			Type:    util.GetParticipants,
			Payload: p,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) StartLiveQuizSession(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.CurrentQuestion = 1
	mod.Status = util.Starting
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.StartLQS,
			Payload: nil,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}

	done := make(chan struct{})
	go h.Countdown(3, c.LiveQuizSessionID, done)
	<-done

	h.DistributeQuestion(c)
}

func (h *Handler) NextQuestion(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.CurrentQuestion += 1
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.DistributeQuestion(c)
}

func (h *Handler) DistributeQuestion(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.Status = util.Questioning
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.DistQuestion,
			Payload: mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1],
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}

	done := make(chan struct{})
	go h.Countdown(5, c.LiveQuizSessionID, done)
	<-done

	mediaType := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["media_type"]
	if mediaType == "" {
		h.DistributeOptions(c)
	} else {
		h.DistributeMedia(c)
	}
}

func (h *Handler) DistributeMedia(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.Status = util.Media
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.DistMedia,
			Payload: nil,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}

	mediaType := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["media_type"].(string)
	if mediaType == util.Image || mediaType == util.Equation {
		done := make(chan struct{})
		go h.Countdown(15, c.LiveQuizSessionID, done)
		<-done
		h.DistributeOptions(c)
	}
}

func (h *Handler) DistributeOptions(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.Status = util.Answering
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	timeLimit := int(mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["time_limit"].(float64))

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.DistOptions,
			Payload: timeLimit,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}

	done := make(chan struct{})
	go h.Countdown(timeLimit, c.LiveQuizSessionID, done)
	<-done

	h.RevealAnswer(c)
}

func (h *Handler) RevealAnswer(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured @691: %v", err)
		return
	}

	res, err := h.Service.GetResponses(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string))
	if err != nil {
		log.Printf("Error occured @697: %v", err)
		return
	}

	qid, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string)
	if !ok {
		log.Printf("Error occured @703: %v", err)
		return
	}
	qType, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["type"].(string)
	if !ok {
		log.Printf("Error occured @708: %v", err)
		return
	}
	qHaveTimeFactor, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["have_time_factor"].(bool)
	if !ok {
		log.Printf("Error occured @713: %v", err)
		return
	}
	qTimeFactor, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["time_factor"].(float64)
	if !ok {
		log.Printf("Error occured @718: %v", err)
		return
	}
	if !qHaveTimeFactor {
		qTimeFactor = 0
	}
	qTimeLimit, ok := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["time_limit"].(float64)
	if !ok {
		log.Printf("Error occured @723: %v", err)
		return
	}
	qAns, ok := mod.Answers[mod.Orders[mod.CurrentQuestion-1]-1].([]any)
	if !ok {
		log.Printf("Error occured @792: Type assertion failed")
		return
	}

	ansCounts := make(map[string]int)
	if qType == util.Choice || qType == util.TrueFalse {
		for _, a := range qAns {
			ansCounts[a.(map[string]any)["id"].(string)] = 0
		}
	}

	var rpl []AnswerPayload
	switch qType {
	case util.Choice, util.TrueFalse:
		for _, r := range res {
			co, ok := r.(map[string]any)["options"].([]any)
			if !ok {
				log.Printf("Error occured @729: Type assertion failed")
				return
			}
			time, ok := r.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @734: Type assertion failed")
				return
			}
			questionID, err := uuid.Parse(qid)
			if err != nil {
				log.Printf("Error occured @741: %v", err)
				return
			}
			pid, ok := r.(map[string]any)["pid"].(string)
			if !ok {
				log.Printf("Error occured @747: %v", err)
				return
			}
			participantID, err := uuid.Parse(pid)
			if err != nil {
				log.Printf("Error occured @752: %v", err)
				return
			}

			var cAnsRes ChoiceAnswerResponse

			cAnsRes, ansCounts, err = h.Service.CalculateAndSaveChoiceResponse(context.Background(), co, qAns, ansCounts, time, qTimeLimit, qTimeFactor, &Response{
				ID:                uuid.New(),
				LiveQuizSessionID: c.LiveQuizSessionID,
				QuestionID:        questionID,
				ParticipantID:     participantID,
				Type:              qType,
			})
			if err != nil {
				log.Printf("Error occured @792: %v", err)
				return
			}

			rpl = append(rpl, AnswerPayload{
				Answers:       cAnsRes,
				ParticipantID: participantID,
			})
		}
	case util.FillBlank:
		for _, r := range res {
			to, ok := r.(map[string]any)["options"].([]any)
			if !ok {
				log.Printf("Error occured @729: Type assertion failed")
				return
			}
			time, ok := r.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @734: Type assertion failed")
				return
			}
			questionID, err := uuid.Parse(qid)
			if err != nil {
				log.Printf("Error occured @741: %v", err)
				return
			}
			pid, ok := r.(map[string]any)["pid"].(string)
			if !ok {
				log.Printf("Error occured @747: %v", err)
				return
			}
			participantID, err := uuid.Parse(pid)
			if err != nil {
				log.Printf("Error occured @752: %v", err)
				return
			}

			fbAnsRes, err := h.Service.CalculateAndSaveFillBlankResponse(context.Background(), to, qAns, time, qTimeLimit, qTimeFactor, &Response{
				ID:                uuid.New(),
				LiveQuizSessionID: c.LiveQuizSessionID,
				QuestionID:        questionID,
				ParticipantID:     participantID,
				Type:              qType,
			})
			if err != nil {
				log.Printf("Error occured @792: %v", err)
				return
			}

			rpl = append(rpl, AnswerPayload{
				Answers:       fbAnsRes,
				ParticipantID: participantID,
			})
		}
	case util.Paragraph:
		for _, r := range res {
			answer, ok := r.(map[string]any)["options"].(string)
			if !ok {
				log.Printf("Error occured @729: Type assertion failed")
				return
			}
			time, ok := r.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @734: Type assertion failed")
				return
			}
			questionID, err := uuid.Parse(qid)
			if err != nil {
				log.Printf("Error occured @741: %v", err)
				return
			}
			pid, ok := r.(map[string]any)["pid"].(string)
			if !ok {
				log.Printf("Error occured @747: %v", err)
				return
			}
			participantID, err := uuid.Parse(pid)
			if err != nil {
				log.Printf("Error occured @752: %v", err)
				return
			}

			pAnsRes, err := h.Service.CalculateAndSaveParagraphResponse(context.Background(), answer, qAns, time, qTimeLimit, qTimeFactor, &Response{
				ID:                uuid.New(),
				LiveQuizSessionID: c.LiveQuizSessionID,
				QuestionID:        questionID,
				ParticipantID:     participantID,
				Type:              qType,
			})
			if err != nil {
				log.Printf("Error occured @792: %v", err)
				return
			}

			rpl = append(rpl, AnswerPayload{
				Answers:       pAnsRes,
				ParticipantID: participantID,
			})
		}
	case util.Matching:
		for _, r := range res {
			mo, ok := r.(map[string]any)["options"].([]any)
			if !ok {
				log.Printf("Error occured @729: Type assertion failed")
				return
			}
			time, ok := r.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @734: Type assertion failed")
				return
			}
			questionID, err := uuid.Parse(qid)
			if err != nil {
				log.Printf("Error occured @741: %v", err)
				return
			}
			pid, ok := r.(map[string]any)["pid"].(string)
			if !ok {
				log.Printf("Error occured @747: %v", err)
				return
			}
			participantID, err := uuid.Parse(pid)
			if err != nil {
				log.Printf("Error occured @752: %v", err)
				return
			}

			mAnsRes, err := h.Service.CalculateAndSaveMatchingResponse(context.Background(), mo, qAns, time, qTimeLimit, qTimeFactor, &Response{
				ID:                uuid.New(),
				LiveQuizSessionID: c.LiveQuizSessionID,
				QuestionID:        questionID,
				ParticipantID:     participantID,
				Type:              qType,
			})
			if err != nil {
				log.Printf("Error occured @792: %v", err)
				return
			}

			rpl = append(rpl, AnswerPayload{
				Answers:       mAnsRes,
				ParticipantID: participantID,
			})
		}
	case util.Pool:
		for _, r := range res {
			time, ok := r.(map[string]any)["time"].(float64)
			if !ok {
				log.Printf("Error occured @734: Type assertion failed")
				return
			}
			pid, ok := r.(map[string]any)["pid"].(string)
			if !ok {
				log.Printf("Error occured @747: %v", err)
				return
			}
			participantID, err := uuid.Parse(pid)
			if err != nil {
				log.Printf("Error occured @752: %v", err)
				return
			}

			options, ok := r.(map[string]any)["options"].(map[string]any)
			if !ok {
				log.Printf("Error occured @1093: Type assertion failed")
				return
			}

			ansRes := make(map[string]PoolAnswer, 0)
			var marksRes int
			var timeRes int

			for i, o := range options {
				I, err := strconv.Atoi(i)
				if err != nil {
					log.Printf("Error occured @3: %v", err)
					return
				}
				sqType, ok := o.(map[string]any)["type"].(string)
				if !ok {
					log.Printf("Error occured @1109: Type assertion failed")
					return
				}
				sqID, ok := o.(map[string]any)["qid"].(string)
				if !ok {
					log.Printf("Error occured @1114: Type assertion failed")
					return
				}
				subqID, err := uuid.Parse(sqID)
				if err != nil {
					log.Printf("Error occured @1114: %v", err)
					return
				}

				ac := make(map[string]int)
				ans, ok := qAns[I].([]any)
				if !ok {
					log.Printf("Error occured @1114: Type assertion failed")
					return
				}
				if sqType == util.Choice || sqType == util.TrueFalse {
					for _, a := range ans {
						ac[a.(map[string]any)["id"].(string)] = 0
					}
				}

				switch sqType {
				case util.Choice, util.TrueFalse:
					sqContent, ok := o.(map[string]any)["content"].([]any)
					if !ok {
						log.Printf("Error occured @1114: Type assertion failed")
						return
					}

					var cAnsRes ChoiceAnswerResponse

					cAnsRes, ac, err = h.Service.CalculateAndSaveChoiceResponse(context.Background(), sqContent, ans, ac, time, qTimeLimit, qTimeFactor, &Response{
						ID:                uuid.New(),
						LiveQuizSessionID: c.LiveQuizSessionID,
						QuestionID:        subqID,
						ParticipantID:     participantID,
						Type:              sqType,
					})
					if err != nil {
						log.Printf("Error occured @792: %v", err)
						return
					}

					ansRes[i] = PoolAnswer{
						ID:      sqID,
						Type:    sqType,
						Content: cAnsRes.Answers,
					}
					marksRes += *cAnsRes.Marks
					timeRes = cAnsRes.Time
					mod.AnswerCounts[sqID] = ac
				case util.FillBlank:
					sqContent, ok := o.(map[string]any)["content"].([]any)
					if !ok {
						log.Printf("Error occured @1114: Type assertion failed")
						return
					}

					fbAnsRes, err := h.Service.CalculateAndSaveFillBlankResponse(context.Background(), sqContent, ans, time, qTimeLimit, qTimeFactor, &Response{
						ID:                uuid.New(),
						LiveQuizSessionID: c.LiveQuizSessionID,
						QuestionID:        subqID,
						ParticipantID:     participantID,
						Type:              sqType,
					})
					if err != nil {
						log.Printf("Error occured @792: %v", err)
						return
					}

					ansRes[i] = PoolAnswer{
						ID:      sqID,
						Type:    sqType,
						Content: fbAnsRes.Answers,
					}
					marksRes += *fbAnsRes.Marks
					timeRes = fbAnsRes.Time
				case util.Paragraph:
					sqContent, ok := o.(map[string]any)["content"]
					if !ok {
						log.Printf("Error occured @1114: Type assertion failed")
						return
					}

					content := ""
					switch sqContent := sqContent.(type) {
					case string:
						content = sqContent
					case nil:
					}

					pAnsRes, err := h.Service.CalculateAndSaveParagraphResponse(context.Background(), content, ans, time, qTimeLimit, qTimeFactor, &Response{
						ID:                uuid.New(),
						LiveQuizSessionID: c.LiveQuizSessionID,
						QuestionID:        subqID,
						ParticipantID:     participantID,
						Type:              sqType,
					})
					if err != nil {
						log.Printf("Error occured @792: %v", err)
						return
					}

					switch pAnsRes := pAnsRes.(type) {
					case string:
						ansRes[i] = PoolAnswer{
							ID:      sqID,
							Type:    sqType,
							Content: pAnsRes,
						}
						marksRes += 0
						timeRes = 0
					case TextAnswerResponse:
						ansRes[i] = PoolAnswer{
							ID:      sqID,
							Type:    sqType,
							Content: pAnsRes.Answers,
						}
						marksRes += *pAnsRes.Marks
						timeRes = pAnsRes.Time
					default:
					}
				case util.Matching:
					sqContent, ok := o.(map[string]any)["content"].([]any)
					if !ok {
						log.Printf("Error occured @1114: Type assertion failed")
						return
					}

					mAnsRes, err := h.Service.CalculateAndSaveMatchingResponse(context.Background(), sqContent, ans, time, qTimeLimit, qTimeFactor, &Response{
						ID:                uuid.New(),
						LiveQuizSessionID: c.LiveQuizSessionID,
						QuestionID:        subqID,
						ParticipantID:     participantID,
						Type:              sqType,
					})
					if err != nil {
						log.Printf("Error occured @792: %v", err)
						return
					}

					ansRes[i] = PoolAnswer{
						ID:      sqID,
						Type:    sqType,
						Content: mAnsRes.Answers,
					}
					marksRes += *mAnsRes.Marks
					timeRes = mAnsRes.Time
				}
			}

			rpl = append(rpl, AnswerPayload{
				Answers: PoolAnswerResponse{
					Answers: ansRes,
					Marks:   marksRes,
					Time:    timeRes,
				},
				ParticipantID: participantID,
			})
		}
	}

	ps, err := h.Service.GetParticipantsByLiveQuizSessionID(context.Background(), c.LiveQuizSessionID)
	if err != nil {
		log.Printf("Error occured @818: %v", err)
		return
	}

	for _, p := range ps {
		injected := false
		for _, r := range rpl {
			if p.ID == r.ParticipantID {
				h.hub.Inject <- &Message{
					Content: Content{
						Type:    util.RevealAnswer,
						Payload: r.Answers,
					},
					LiveQuizSessionID: c.LiveQuizSessionID,
					ClientID:          p.ID,
					UserID:            c.UserID,
				}
				injected = true
				break
			}
		}
		if !injected {
			h.hub.Inject <- &Message{
				Content: Content{
					Type:    util.RevealAnswer,
					Payload: nil,
				},
				LiveQuizSessionID: c.LiveQuizSessionID,
				ClientID:          p.ID,
				UserID:            c.UserID,
			}
		}

		h.hub.Inject <- &Message{
			Content: Content{
				Type:    util.UpdateMarks,
				Payload: p.Marks,
			},
			LiveQuizSessionID: c.LiveQuizSessionID,
			ClientID:          p.ID,
			UserID:            c.UserID,
		}
	}

	mod.Status = util.RevealingAnswer
	mod.AnswerCounts[qid] = ansCounts

	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured @699: %v", err)
		return
	}

	var hostPID uuid.UUID
	for _, cl := range h.hub.LiveQuizSessions[c.LiveQuizSessionID].Clients {
		if cl.IsHost {
			hostPID = cl.ID
		}
	}

	correctAns, err := h.Service.GetAnswersResponseForHost(context.Background(), qid, qType, qAns, mod.AnswerCounts)
	if err != nil {
		log.Printf("Error occured @699: %v", err)
		return
	}

	h.hub.Inject <- &Message{
		Content: Content{
			Type:    util.RevealAnswer,
			Payload: correctAns,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          hostPID,
		UserID:            c.UserID,
	}
}

func (h *Handler) Conclude(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	mod.Status = util.Concluding
	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	participants, err := h.Service.GetParticipantsByLiveQuizSessionID(context.Background(), c.LiveQuizSessionID)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	for _, p := range participants {
		rank, err := h.Service.GetRank(context.Background(), p.LiveQuizSessionID, p.ID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		h.hub.Inject <- &Message{
			Content: Content{
				Type:    util.Conclude,
				Payload: rank,
			},
			LiveQuizSessionID: p.LiveQuizSessionID,
			ClientID:          p.ID,
			UserID:            p.UserID,
		}
	}

	h.hub.Inject <- &Message{
		Content: Content{
			Type:    util.Conclude,
			Payload: nil,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) SubmitAnswer(c *Client, payload any) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	code := h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code
	pid := c.ID.String()
	qid := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string)

	exist, err := h.Service.DoesResponseExist(context.Background(), code, qid, pid)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	if exist {
		if err := h.Service.UpdateResponse(context.Background(), code, qid, pid, payload); err != nil {
			log.Printf("Error occured: %v", err)
			return
		}
	} else {
		if err := h.Service.CreateResponse(context.Background(), code, qid, pid, payload); err != nil {
			log.Printf("Error occured at CreateResponse: %v", err)
			return
		}
	}

	count, err := h.Service.CountResponses(context.Background(), code, qid)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	if !mod.Config.ParticipantConfig.Reanswer && count == mod.ParticipantCount && mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["type"].(string) != util.Pool {
		mod.Interrupted = true
	}
	mod.ResponseCount = count

	if err = h.Service.UpdateLiveQuizSessionCache(context.Background(), code, mod); err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Converse <- &Message{
		Content: Content{
			Type:    util.SubmitAnswer,
			Payload: count,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) UnsubmitAnswer(c *Client) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	code := h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code
	pid := c.ID.String()
	qid := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["id"].(string)

	if err := h.Service.FlushResponse(context.Background(), code, qid, pid); err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	count, err := h.Service.CountResponses(context.Background(), code, qid)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	mod.ResponseCount = count

	if err = h.Service.UpdateLiveQuizSessionCache(context.Background(), code, mod); err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.hub.Converse <- &Message{
		Content: Content{
			Type:    util.UnsubmitAnswer,
			Payload: count,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) UpdateModerator(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	uid, ok := c.Get("uid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	}

	var userID *uuid.UUID
	if uid.(string) == "NOT_HOST" {
		userID = nil
	} else {
		parsedUID, err := uuid.Parse(uid.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userID = &parsedUID
	}

	isHost := false
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code && userID != nil && s.HostID == *userID {
			isHost = true
			break
		}
	}

	mod, err := h.Service.GetLiveQuizSessionCache(c, code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	if !isHost {
		mod.Answers = make([]any, 0)
	}

	c.JSON(http.StatusOK, mod)
}

func (h *Handler) BroadcastMessage(c *Client, ct Content) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    ct.Type,
			Payload: ct.Payload,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) Converse(c *Client, ct Content) {
	h.hub.Converse <- &Message{
		Content: Content{
			Type:    ct.Type,
			Payload: ct.Payload,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		ClientID:          c.ID,
		UserID:            c.UserID,
	}
}

func (h *Handler) Countdown(seconds int, lqsID uuid.UUID, cd chan<- struct{}) {
	for i := float64(seconds) * 10; i > 0; i -= 1 {
		if _, ok := h.hub.LiveQuizSessions[lqsID]; ok {
			mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
			if err != nil {
				log.Printf("Error occured: %v", err)
				break
			}

			if mod.Interrupted {
				mod.Interrupted = false
				err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, mod)
				if err != nil {
					log.Printf("Error occured: %v", err)
					break
				}
				break
			}
			h.hub.Broadcast <- &Message{
				Content: Content{
					Type: util.Countdown,
					Payload: CountDownPayload{
						TimeLeft:        float64(i / 10),
						CurrentQuestion: mod.CurrentQuestion,
						Status:          mod.Status,
					},
				},
				LiveQuizSessionID: lqsID,
				ClientID:          h.hub.LiveQuizSessions[lqsID].ID,
				UserID:            &h.hub.LiveQuizSessions[lqsID].HostID,
			}
			if i > 0 {
				time.Sleep(time.Second / 10)
			}
		}
	}
	close(cd)
}
