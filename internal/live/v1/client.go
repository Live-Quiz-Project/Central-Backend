package v1

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn              *websocket.Conn
	Message           chan *Message
	ID                uuid.UUID  `json:"id"`
	UserID            *uuid.UUID `json:"uid"`
	DisplayName       string     `json:"display_name"`
	DisplayEmoji      string     `json:"display_emoji"`
	DisplayColor      string     `json:"display_color"`
	IsHost            bool       `json:"isHost"`
	LiveQuizSessionID uuid.UUID  `json:"lqsId"`
	Status            string     `json:"status"`
}

type Message struct {
	Content           Content    `json:"content"`
	LiveQuizSessionID uuid.UUID  `json:"live_quiz_session_id"`
	ClientID          uuid.UUID  `json:"client_id"`
	UserID            *uuid.UUID `json:"uid"`
}

type Content struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		m, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(m)
	}
}

func (c *Client) readMessage(h *Handler) {
	defer func() {
		log.Println("Closing connection")
		if !c.IsHost {
			p, err := h.Service.GetParticipantByID(context.Background(), c.ID)
			if err != nil {
				log.Printf("Error occured at defer readMessage: %v", err)
				h.hub.Unregister <- c
				c.Conn.Close()
			}
			p.Status = util.Left
			if _, err = h.Service.UpdateParticipant(context.Background(), p); err != nil {
				log.Printf("Error occured at defer readMessage: %v", err)
				h.hub.Unregister <- c
				c.Conn.Close()
			}
			h.hub.Unregister <- c
		}
		h.hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) || websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Websocket error occured: %v from isHost:%v", err, c.IsHost)

				if _, ok := h.hub.LiveQuizSessions[c.LiveQuizSessionID]; !ok {
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
					mod.ParticipantCount = pCount - 1

					err = h.Service.UpdateLiveQuizSessionCache(context.Background(), h.hub.LiveQuizSessions[c.LiveQuizSessionID].Code, mod)
					if err != nil {
						log.Printf("Error occured: %v", err)
						return
					}
				}

				break
			}
		}

		var mstr Content
		if err := json.Unmarshal(m, &mstr); err != nil {
			log.Printf("Error occured: %v @86", err)
			break
		}

		if _, ok := h.hub.LiveQuizSessions[c.LiveQuizSessionID]; !ok {
			log.Println("No such session exists")
			return
		}

		switch mstr.Type {
		case util.JoinLQS, util.LeaveLQS:
			c.Converse(h, mstr)
		case util.KickParticipant:
			c.KickParticipant(h, mstr.Payload)
		case util.StartLQS:
			c.StartLiveQuizSession(h)
		case util.NextQuestion:
			c.NextQuestion(h)
		case util.DistQuestion:
			c.DistributeQuestion(h)
		case util.DistMedia:
			c.DistributeMedia(h)
		case util.DistOptions:
			c.DistributeOptions(h)
		case util.RevealAnswer:
			c.RevealAnswer(h)
		case util.Conclude:
			c.Conclude(h)
		case util.ToggleLock:
			c.ToggleLiveQuizSessionLock(h)
		case util.GetParticipants:
			c.GetParticipants(h)
		case util.SubmitAnswer:
			c.SubmitAnswer(h, mstr.Payload)
		case util.UnsubmitAnswer:
			c.UnsubmitAnswer(h)
		default:
			c.BroadcastMessage(h, mstr)
		}
	}
}

func (c *Client) KickParticipant(h *Handler, payload any) {
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

func (c *Client) ToggleLiveQuizSessionLock(h *Handler) {
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

func (c *Client) GetParticipants(h *Handler) {
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

func (c *Client) StartLiveQuizSession(h *Handler) {
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
	go c.Countdown(h, 3, c.LiveQuizSessionID, done)
	<-done

	c.DistributeQuestion(h)
}

func (c *Client) DistributeQuestion(h *Handler) {
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
	go c.Countdown(h, 5, c.LiveQuizSessionID, done)
	<-done

	mediaType := mod.Questions[mod.Orders[mod.CurrentQuestion-1]-1].(map[string]any)["media_type"]
	if mediaType == "" {
		c.DistributeOptions(h)
	} else {
		c.DistributeMedia(h)
	}
}

func (c *Client) NextQuestion(h *Handler) {
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

	c.DistributeQuestion(h)
}

func (c *Client) DistributeMedia(h *Handler) {
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
		go c.Countdown(h, 15, c.LiveQuizSessionID, done)
		<-done
		c.DistributeOptions(h)
	}
}

func (c *Client) DistributeOptions(h *Handler) {
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
	go c.Countdown(h, timeLimit, c.LiveQuizSessionID, done)
	<-done

	c.RevealAnswer(h)
}

func (c *Client) RevealAnswer(h *Handler) {
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

func (c *Client) Conclude(h *Handler) {
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

func (c *Client) SubmitAnswer(h *Handler, payload any) {
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

func (c *Client) UnsubmitAnswer(h *Handler) {
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

func (c *Client) Converse(h *Handler, ct Content) {
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

func (c *Client) BroadcastMessage(h *Handler, ct Content) {
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

func (c *Client) Countdown(h *Handler, seconds int, lqsID uuid.UUID, cd chan<- struct{}) {
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
