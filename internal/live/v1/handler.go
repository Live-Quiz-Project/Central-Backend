package v1

import (
	"log"
	"net/http"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
	Service
}

func NewHandler(h *Hub, s Service) *Handler {
	return &Handler{
		hub:     h,
		Service: s,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(req *http.Request) bool {
		origin := req.Header.Get("Origin")
		log.Println(origin)
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

		err = h.Service.CreateLiveQuizSessionCache(c, &LiveQuizSession{
			Session: Session{
				ID:                  lqs.ID,
				HostID:              hostID,
				QuizID:              lqs.QuizID,
				Status:              util.Idle,
				ExemptedQuestionIDs: nil,
			},
			Code: lqs.Code,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	if userID != h.hub.LiveQuizSessions[userID].HostID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only the host can end the session"})
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

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.EndLQS,
			Payload: "Session has ended.",
		},
		LiveQuizSessionID: lqsID,
		UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
	}

	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
		if err := cl.Conn.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't close client connection"})
			return
		}
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
	color := c.Query("color")

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
		Name:              uname,
		Marks:             0,
	}

	exists, eErr := h.DoesParticipantExists(c, userID)
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
	}

	_, pErr := h.CreateParticipant(c, p)
	if pErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": pErr.Error()})
		return
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

	h.hub.Broadcast <- &Message{
		Content: Content{
			Type:    util.JoinedLQS,
			Payload: nil,
		},
		LiveQuizSessionID: lqsID,
		UserID:            userID,
	}

	go cl.writeMessage()
	cl.readMessage(h)
}

func (h *Handler) GetHost(c *gin.Context) {}

func (h *Handler) GetParticipants(c *gin.Context) {
	// 	var p []GetParticipantsResponse

	// 	sessionID := c.Param("id")
	// 	lqsID, err := uuid.Parse(sessionID)

	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
	// 		return
	// 	}
	// 	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
	// 		p = make([]GetParticipantsResponse, 0)
	// 		c.JSON(http.StatusOK, p)
	// 		return
	// 	}

	//	for _, c := range h.hub.LiveQuizSessions[lqsID].Clients {
	//		if !c.IsHost {
	//			p = append(p, GetParticipantsResponse{
	//				ID:           c.ID,
	//				Name:         c.Name,
	//				Emoji:        c.Emoji,
	//				ProfileColor: c.ProfileColor,
	//			})
	//		}
	//	}
	//
	// c.JSON(http.StatusOK, p)
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

// func (h *Handler) startLiveQuizSession(c *Client, payload any) {
// 	h.hub.Broadcast <- &Message{
// 		Content: Content{
// 			Type:    util.StartLQS,
// 			Payload: nil,
// 		},
// 		LiveQuizSessionID: c.LiveQuizSessionID,
// 		UserID:            c.ID,
// 	}

// 	p, ok := payload.(map[string]any)
// 	if !ok {
// 		log.Printf("Error occured: %v", "Payload is not a map")
// 		return
// 	}

// 	sessionID := p["id"].(string)
// 	lqsID, err := uuid.Parse(sessionID)
// 	if err != nil {
// 		log.Printf("Error occured: %v", err)
// 		return
// 	}

// 	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
// 		return
// 	}

// 	// todo: update mod
// 	// h.hub.Moderate <- &Moderator{
// 	// 	LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 	// 	QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 	// 	CurrentQuestion:   1,
// 	// 	Status:            util.Starting,
// 	// }

// 	done := make(chan struct{})
// 	go h.countdown(3, lqsID, done)
// 	<-done

// 	h.distributeQuestion(payload, false)
// }

// func (h *Handler) distributeQuestion(payload any, isNext bool) {
// 	p, ok := payload.(map[string]any)
// 	if !ok {
// 		log.Printf("Error occured: %v", "Payload is not a map")
// 		return
// 	}

// 	sessionID := p["id"].(string)
// 	lqsID, err := uuid.Parse(sessionID)
// 	if err != nil {
// 		log.Printf("Error occured: %v", err)
// 		return
// 	}

// 	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
// 		return
// 	}
// 	var hostID uuid.UUID
// 	quizID := h.hub.LiveQuizSessions[lqsID].QuizID
// 	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
// 		if cl.IsHost {
// 			hostID = cl.ID
// 			break
// 		}
// 	}

// 	//	todo: update mod
// 	// h.hub.Moderate <- &Moderator{
// 	// 	LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 	// 	QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 	// 	CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 	// 	Status:            util.Questioning,
// 	// }

// 	// query := fmt.Sprintf(`
// 	// 	SELECT
// 	// 		*
// 	// 	FROM
// 	// 		question
// 	// 	WHERE
// 	// 		quiz_id = '%s'
// 	// 		AND "order" = '%d';`,
// 	// 	quizID,
// 	// 	h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 	// )
// 	// row := db.DB.QueryRow(query)
// 	// type question struct {
// 	// 	ID             string `json:"id"`
// 	// 	QuizID         string `json:"quizId"`
// 	// 	IsParent       bool   `json:"isParent"`
// 	// 	ParentID       string `json:"parentId"`
// 	// 	Type           string `json:"type"`
// 	// 	Order          int    `json:"order"`
// 	// 	Content        string `json:"content"`
// 	// 	Note           string `json:"note"`
// 	// 	Media          string `json:"media"`
// 	// 	TimeLimit      int    `json:"timeLimit"`
// 	// 	HaveTimeFactor bool   `json:"haveTimeFactor"`
// 	// 	TimeFactor     int    `json:"timeFactor"`
// 	// 	FontSize       int    `json:"fontSize"`
// 	// 	SelectUpTo     int    `json:"selectUpTo"`
// 	// }
// 	// var q question
// 	// err := row.Scan(
// 	// 	&q.ID,
// 	// 	&q.QuizID,
// 	// 	&q.IsParent,
// 	// 	&q.ParentID,
// 	// 	&q.Type,
// 	// 	&q.Order,
// 	// 	&q.Content,
// 	// 	&q.Note,
// 	// 	&q.Media,
// 	// 	&q.TimeLimit,
// 	// 	&q.HaveTimeFactor,
// 	// 	&q.TimeFactor,
// 	// 	&q.FontSize,
// 	// 	&q.SelectUpTo,
// 	// )
// 	// if err != nil {
// 	// 	log.Printf("Error occurred while scanning question data: %v", err)
// 	// 	return
// 	// }

// 	// type responsePayload struct {
// 	// 	Question  question      `json:"question"`
// 	// 	Moderator Moderator `json:"mod"`
// 	// }
// 	// pl := responsePayload{
// 	// 	Question: q,
// 	// 	Moderator: Moderator{
// 	// 		LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 	// 		QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 	// 		CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 	// 		Status:            util.Questioning,
// 	// 	},
// 	// }

// 	// h.hub.Broadcast <- &Message{
// 	// 	Content: Content{
// 	// 		Type:    util.DistQuestion,
// 	// 		Payload: pl,
// 	// 	},
// 	// 	LiveQuizSessionID: lqsID,
// 	// 	UserID:            hostID,
// 	// }

// 	done := make(chan struct{})
// 	go h.countdown(5, lqsID, done)
// 	<-done

// 	h.distributeOptions(payload, q.ID, q.Type)
// }

// func (h *Handler) distributeOptions(payload any, qID string, qType string) {
// 	p, ok := payload.(map[string]any)
// 	if !ok {
// 		log.Printf("Error occured: %v", "Payload is not a map")
// 		return
// 	}

// 	sessionID := p["id"].(string)
// 	lqsID, err := uuid.Parse(sessionID)
// 	if err != nil {
// 		log.Printf("Error occured: %v", err)
// 		return
// 	}

// 	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
// 		return
// 	}
// 	var hostID uuid.UUID
// 	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
// 		if cl.IsHost {
// 			hostID = cl.ID
// 			break
// 		}
// 	}
// 	quizID := h.hub.LiveQuizSessions[lqsID].QuizID

// 	// todo: update mod
// 	// h.hub.Moderate <- &Moderator{
// 	// 	LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 	// 	QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 	// 	CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 	// 	Status:            lqsUtil.Answering,
// 	// }

// 	// qQuery := fmt.Sprintf(`
// 	// 	SELECT
// 	// 		q.time_limit,
// 	// 		q.selected_up_to
// 	// 	FROM
// 	// 		question q
// 	// 	WHERE
// 	// 		q.id = '%s';`,
// 	// 	qID,
// 	// )
// 	// type question struct {
// 	// 	TimeLimit  int `json:"timeLimit"`
// 	// 	SelectUpTo int `json:"selectUpTo"`
// 	// }
// 	// var q question
// 	// err := db.DB.QueryRow(qQuery).Scan(
// 	// 	&q.TimeLimit,
// 	// 	&q.SelectUpTo,
// 	// )
// 	// if err != nil {
// 	// 	log.Printf("Error occurred while scanning question data: %v", err)
// 	// 	return
// 	// }

// 	switch qType {
// 	case util.Choice, util.TrueFalse:
// 		// oQuery := fmt.Sprintf(`
// 		// 	SELECT
// 		// 		o.*,
// 		// 		q.content AS question_content
// 		// 	FROM
// 		// 		question q
// 		// 		LEFT JOIN option_choice o ON q.id = o.question_id
// 		// 	WHERE
// 		// 		q.question_id = '%s'
// 		// 		AND q.quiz_id = '%s';`,
// 		// 	qID,
// 		// 	quizID,
// 		// )
// 		// type option struct {
// 		// 	ID              string `json:"id"`
// 		// 	QuestionID      string `json:"qId"`
// 		// 	Order           int    `json:"order"`
// 		// 	Content         string `json:"content"`
// 		// 	Mark            int    `json:"mark"`
// 		// 	Color           string `json:"color"`
// 		// 	IsCorrect       bool   `json:"isCorrect"`
// 		// 	QuestionContent string `json:"qContent"`
// 		// }
// 		// rows, err := db.DB.Query(oQuery)
// 		// if err != nil {
// 		// 	log.Printf("Error occurred while executing the SQL query: %v", err)
// 		// 	return
// 		// }
// 		// defer rows.Close()

// 		// options := make([]option, 0)
// 		// for rows.Next() {
// 		// 	var o option
// 		// 	err := rows.Scan(
// 		// 		&o.ID,
// 		// 		&o.QuestionID,
// 		// 		&o.Order,
// 		// 		&o.Content,
// 		// 		&o.Mark,
// 		// 		&o.Color,
// 		// 		&o.IsCorrect,
// 		// 		&o.QuestionContent,
// 		// 	)
// 		// 	if err != nil {
// 		// 		log.Printf("Error occurred while scanning option data: %v", err)
// 		// 		return
// 		// 	}
// 		// 	options = append(options, o)
// 		// }

// 		// type responsePayload struct {
// 		// 	Options    []option  `json:"options"`
// 		// 	SelectUpTo int       `json:"selectUpTo"`
// 		// 	Moderator  Moderator `json:"mod"`
// 		// }
// 		// pl := responsePayload{
// 		// 	Options:    options,
// 		// 	SelectUpTo: q.SelectUpTo,
// 		// 	Moderator: Moderator{
// 		// 		LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 		// 		QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 		// 		CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 		// 		Status:            lqsUtil.Answering,
// 		// 	},
// 		// }

// 		// h.hub.Broadcast <- &Message{
// 		// 	Content: Content{
// 		// 		Type:    lqsUtil.DistOptions,
// 		// 		Payload: pl,
// 		// 	},
// 		// 	LiveQuizSessionID: lqsID,
// 		// 	UserID:            hostID,
// 		// }

// 		// done := make(chan struct{})
// 		// go h.countdown(q.TimeLimit, lqsID, done)
// 		// <-done
// 	}

// 	h.revealAnswer(payload, qID, qType)
// }

// func (h *Handler) revealAnswer(payload any, qID string, qType string) {
// 	log.Println("Reveal answer")
// 	p, ok := payload.(map[string]any)
// 	if !ok {
// 		log.Printf("Error occured: %v", "Payload is not a map")
// 		return
// 	}

// 	sessionID := p["id"].(string)
// 	lqsID, err := uuid.Parse(sessionID)
// 	if err != nil {
// 		log.Printf("Error occured: %v", err)
// 		return
// 	}

// 	if _, ok := h.hub.LiveQuizSessions[lqsID]; !ok {
// 		return
// 	}
// 	var hostID uuid.UUID
// 	for _, cl := range h.hub.LiveQuizSessions[lqsID].Clients {
// 		if cl.IsHost {
// 			hostID = cl.ID
// 			break
// 		}
// 	}
// 	quizID := h.hub.LiveQuizSessions[lqsID].QuizID

// 	switch qType {
// 	case util.Choice, util.TrueFalse:
// 		// 	query := fmt.Sprintf(`
// 		// 		SELECT
// 		// 			o.content,
// 		// 			o.color,
// 		// 			q.content AS question_content
// 		// 		FROM
// 		// 			question q
// 		// 			INNER JOIN option_choice o ON q.id = o.question_id
// 		// 		WHERE
// 		// 			q.question_id = '%s'
// 		// 			AND q.quiz_id = '%s'
// 		// 			AND o.is_correct = TRUE;`,
// 		// 		qID,
// 		// 		quizID,
// 		// 	)
// 		// 	type answer struct {
// 		// 		Content         string `json:"content"`
// 		// 		Color           string `json:"color"`
// 		// 		QuestionContent string `json:"qContent"`
// 		// 	}
// 		// 	var a answer
// 		// 	err := db.DB.QueryRow(query).Scan(
// 		// 		&a.Content,
// 		// 		&a.Color,
// 		// 		&a.QuestionContent,
// 		// 	)
// 		// 	if err != nil {
// 		// 		log.Printf("Error occurred while scanning option data: %v", err)
// 		// 		return
// 		// 	}

// 		// 	type responsePayload struct {
// 		// 		Answer    answer    `json:"answer"`
// 		// 		Moderator Moderator `json:"mod"`
// 		// 	}
// 		// 	pl := responsePayload{
// 		// 		Answer: a,
// 		// 		Moderator: Moderator{
// 		// 			LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 		// 			QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 		// 			CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 		// 			Status:            lqsUtil.RevealingAnswer,
// 		// 		},
// 		// 	}

// 		// 	h.hub.Broadcast <- &Message{
// 		// 		Content: Content{
// 		// 			Type:    lqsUtil.RevealAnswer,
// 		// 			Payload: pl,
// 		// 		},
// 		// 		LiveQuizSessionID: lqsID,
// 		// 		UserID:            hostID,
// 		// 	}
// 	}

// 	// h.hub.Moderate <- &Moderator{
// 	// 	LiveQuizSessionID: h.hub.LiveQuizSessions[lqsID].ID,
// 	// 	QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 	// 	CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion + 1,
// 	// 	Status:            lqsUtil.RevealingAnswer,
// 	// }
// }

// func (h *Handler) countdown(seconds int, lqsID uuid.UUID, cd chan<- struct{}) {
// 	for i := seconds; i > 0; i-- {
// 		if _, ok := h.hub.LiveQuizSessions[lqsID]; ok {
// 			h.hub.Broadcast <- &Message{
// 				Content: Content{
// 					Type: util.Countdown,
// 					Payload: &CountDownPayload{
// 						LiveQuizSessionID: lqsID,
// 						TimeLeft:          i,
// 						// QuestionCount:     h.hub.LiveQuizSessions[lqsID].Moderator.QuestionCount,
// 						// CurrentQuestion:   h.hub.LiveQuizSessions[lqsID].Moderator.CurrentQuestion,
// 						// Status:            h.hub.LiveQuizSessions[lqsID].Moderator.Status,
// 					},
// 				},
// 				LiveQuizSessionID: lqsID,
// 				UserID:            h.hub.LiveQuizSessions[lqsID].HostID,
// 			}
// 			if i > 0 {
// 				time.Sleep(time.Second)
// 			}
// 		}
// 	}
// 	close(cd)
// }
