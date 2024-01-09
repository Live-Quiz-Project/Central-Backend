package v1

import (
	"context"
	"log"
	"net/http"
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
		// origin := req.Header.Get("Origin")
		// return origin == "http://localhost:5173" || origin == "http://localhost:5174" || origin == "http://localhost:3000"
		return true
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

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		code = util.CodeGenerator(codes)
		lqs, err := h.Service.CreateLiveQuizSession(c, &req, lqsID, code, hostID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		count, err := h.quizService.GetQuestionCountByQuizID(c, lqs.QuizID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		orders := make([]int, count)
		if req.Config.ShuffleConfig.Question {
			orders = util.ShuffleNumbers(count)
		} else {
			for i := 0; i < count; i++ {
				orders[i] = i + 1
			}
		}
		log.Println(orders)

		err = h.Service.CreateLiveQuizSessionCache(context.Background(), code, &Cache{
			ID:              lqsID,
			QuizID:          lqs.QuizID,
			QuestionCount:   count,
			CurrentQuestion: -1,
			Question:        nil,
			Options:         nil,
			Answers:         nil,
			Status:          util.Idle,
			Config:          req.Config,
			Orders:          orders,
		})
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		h.hub.LiveQuizSessions[lqsID] = &LiveQuizSession{
			Session: Session{ID: lqs.ID,
				HostID:              hostID,
				QuizID:              lqs.QuizID,
				Status:              util.Idle,
				ExemptedQuestionIDs: nil},
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

	userID, err := uuid.Parse(uid.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	sessionID := c.Param("id")
	lqsID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	if userID != h.hub.LiveQuizSessions[lqsID].HostID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only the host can end the session"})
		return
	}

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.EndLQS,
			Payload: "Session has ended.",
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
		h.hub.Unregister <- cl
	}

	err = h.FlushLiveQuizSessionCache(c, h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	err = h.FlushLiveQuizSessionResponsesCache(c, h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	delete(h.hub.LiveQuizSessions, lqsID)

	c.JSON(http.StatusOK, gin.H{"message": "Successfully ended the session"})
}

func (h *Handler) CheckLiveQuizSessionAvailability(c *gin.Context) {
	code := c.Param("code")

	lqsCache, err := h.Service.GetLiveQuizSessionCache(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code {
			c.JSON(http.StatusOK, &CheckLiveQuizSessionAvailabilityResponse{
				ID:              s.ID,
				QuizID:          s.QuizID,
				Code:            s.Code,
				QuestionCount:   lqsCache.QuestionCount,
				CurrentQuestion: lqsCache.CurrentQuestion,
				Status:          lqsCache.Status,
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
}

func (h *Handler) JoinLiveQuizSession(c *gin.Context) {
	code := c.Param("code")

	uid := c.Query("id")
	userID, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var isHost bool
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code && s.HostID == userID {
			isHost = true
			break
		}
	}

	uname := c.Query("name")
	emoji := c.Query("emoji")
	color := "#" + c.Query("color")

	var lqsID uuid.UUID
	for _, s := range h.hub.LiveQuizSessions {
		if s.Code == code {
			lqsID = s.ID
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

	p := &Participant{
		ID:                uuid.New(),
		UserID:            &userID,
		LiveQuizSessionID: lqsID,
		Status:            util.Joined,
		Marks:             0,
		Name:              uname,
	}

	if !isHost {
		exists, eErr := h.DoesParticipantExists(c, userID, lqsID)
		if eErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": eErr.Error()})
			return
		}
		if exists {
			_, err := h.UpdateParticipantStatus(c, *p.UserID, p.LiveQuizSessionID, util.Joined)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			_, pErr := h.CreateParticipant(c, p)
			if pErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": pErr.Error()})
				return
			}
		}
	}

	cl := &Client{
		Conn:              conn,
		Message:           make(chan *Message, 10),
		ID:                userID,
		DisplayName:       uname,
		IsHost:            isHost,
		LiveQuizSessionID: lqsID,
		Status:            util.Joined,
		DisplayEmoji:      emoji,
		DisplayColor:      color,
		Marks:             0,
	}

	h.hub.Register <- cl

	go cl.writeMessage()

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.JoinedLQS,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            userID,
	}

	cl.readMessage(h)
}

func (h *Handler) GetHost(c *gin.Context) {
	sessionID := c.Param("id")
	lqsID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	var host Client
	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
		if cl.IsHost {
			host = *cl
			break
		}
	}

	if host.ID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No host found"})
		return
	}

	c.JSON(http.StatusOK, &Participant{
		ID:                host.ID,
		UserID:            &host.ID,
		LiveQuizSessionID: host.LiveQuizSessionID,
		Status:            host.Status,
		Name:              host.DisplayName,
		Marks:             host.Marks,
	})
}

func (h *Handler) GetParticipants(c *gin.Context) {
	sessionID := c.Param("id")
	lqsID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No such session exists"})
		return
	}

	participants, err := h.GetParticipantsByLiveQuizSessionID(c, lqsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, participants.Participants)
}

func (h *Handler) SendMessage(c *Client, ct Content) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    ct.Type,
			Payload: ct.Payload,
		},
		LiveQuizSessionID: c.LiveQuizSessionID,
		UserID:            c.ID,
	}
}

func (h *Handler) UnregisterParticipants(c *Client) {
	h.Service.UnregisterParticipants(context.Background(), c.LiveQuizSessionID)

	for _, cl := range h.hub.LiveQuizSessions[c.LiveQuizSessionID].Clients {
		h.hub.Unregister <- cl
	}

	delete(h.hub.LiveQuizSessions, c.LiveQuizSessionID)
}

func (h *Handler) StartLiveQuizSession(lqsID uuid.UUID) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.StartLQS,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		return
	}

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, &Cache{
		ID:              mod.ID,
		QuizID:          mod.QuizID,
		QuestionCount:   mod.QuestionCount,
		Question:        nil,
		Options:         nil,
		Answers:         nil,
		CurrentQuestion: mod.CurrentQuestion,
		Status:          util.Starting,
		Config:          mod.Config,
		Orders:          mod.Orders,
	})
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	done := make(chan struct{})
	go h.Countdown(3, lqsID, done)
	<-done

	h.DistributeQuestion(mod.ID)
}

func (h *Handler) DistributeQuestion(lqsID uuid.UUID) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.DistQuestion,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		return
	}

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	q, err := h.quizService.GetQuestionByQuizIDAndOrder(context.Background(), h.hub.LiveQuizSessions[lqsID].QuizID, mod.CurrentQuestion+1)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, &Cache{
		ID:              mod.ID,
		QuizID:          mod.QuizID,
		QuestionCount:   mod.QuestionCount,
		Question:        q,
		Options:         nil,
		Answers:         nil,
		CurrentQuestion: mod.CurrentQuestion + 1,
		Status:          util.Questioning,
		Config:          mod.Config,
		Orders:          mod.Orders,
	})
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	h.Service.CreateLiveQuizSessionResponsesCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, nil)

	done := make(chan struct{})
	go h.Countdown(5, lqsID, done)
	<-done

	log.Println("Distributing options")
	h.DistributeOptions(lqsID)
}

func (h *Handler) DistributeOptions(lqsID uuid.UUID) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.DistOptions,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		log.Println("No such session exists")
		return
	}

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}
	log.Println(mod.Status)

	switch mod.Question.Type {
	case util.Choice, util.TrueFalse:
		options, err := h.quizService.GetChoiceOptionsByQuestionID(context.Background(), mod.Question.ID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, &Cache{
			ID:              mod.ID,
			QuizID:          mod.QuizID,
			QuestionCount:   mod.QuestionCount,
			Question:        mod.Question,
			Options:         options,
			Answers:         nil,
			CurrentQuestion: mod.CurrentQuestion,
			Status:          util.Answering,
			Config:          mod.Config,
			Orders:          mod.Orders,
		})
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}
	}

	done := make(chan struct{})
	go h.Countdown(mod.Question.TimeLimit, lqsID, done)
	<-done

	// res, er := h.Service.GetLiveQuizSessionResponseCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	// if er != nil {
	// 	log.Printf("Error occured: %v", er)
	// 	return
	// }

	h.RevealAnswer(lqsID)
}

func (h *Handler) RevealAnswer(lqsID uuid.UUID) {
	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.RevealAnswer,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
		return
	}

	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	switch mod.Question.Type {
	case util.Choice, util.TrueFalse:
		answers, err := h.quizService.GetChoiceAnswersByQuestionID(context.Background(), mod.Question.ID)
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code, &Cache{
			ID:              mod.ID,
			QuizID:          mod.QuizID,
			QuestionCount:   mod.QuestionCount,
			Question:        mod.Question,
			Options:         mod.Options,
			Answers:         answers,
			CurrentQuestion: mod.CurrentQuestion,
			Status:          util.RevealingAnswer,
			Config:          mod.Config,
			Orders:          mod.Orders,
		})
		if err != nil {
			log.Printf("Error occured: %v", err)
			return
		}

		// get own response and check if correct
	}
}

func (h *Handler) GetLiveQuizSessionCache(c *gin.Context) {
	sessionID := c.Param("id")
	lqsID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	cache, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	c.JSON(http.StatusOK, cache)
}

func (h *Handler) GetLiveQuizSessionResponsesCache(c *gin.Context) {
	sessionID := c.Param("id")
	lqsID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	cache, err := h.Service.GetLiveQuizSessionResponsesCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	c.JSON(http.StatusOK, cache)
}

func (h *Handler) Countdown(seconds int, lqsID uuid.UUID, cd chan<- struct{}) {
	mod, err := h.Service.GetLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[lqsID].Code)
	if err != nil {
		log.Printf("Error occured: %v", err)
		return
	}

	for i := seconds; i > 0; i-- {
		if _, ok := h.hub.LiveQuizSessions[lqsID]; ok {
			h.hub.Broadcast <- &Message{
				Content: Content{
					Type: util.Countdown,
					Payload: &CountDownPayload{
						LiveQuizSessionID: lqsID,
						TimeLeft:          i,
						QuestionCount:     mod.QuestionCount,
						CurrentQuestion:   mod.CurrentQuestion,
						Status:            mod.Status,
					},
				},
				LiveQuizSessionID: lqsID,
				UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
			}
			if i > 0 {
				time.Sleep(time.Second)
			}
		}
	}
	close(cd)
}
