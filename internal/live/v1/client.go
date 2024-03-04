package v1

import (
	"context"
	"encoding/json"
	"log"

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
			h.Converse(c, mstr)
		case util.KickParticipant:
			h.KickParticipant(c, mstr.Payload)
		case util.StartLQS:
			h.StartLiveQuizSession(c)
		case util.NextQuestion:
			h.NextQuestion(c)
		case util.DistQuestion:
			h.DistributeQuestion(c)
		case util.DistMedia:
			h.DistributeMedia(c)
		case util.DistOptions:
			h.DistributeOptions(c)
		case util.RevealAnswer:
			h.RevealAnswer(c)
		case util.Conclude:
			h.Conclude(c)
		case util.ToggleLock:
			h.ToggleLiveQuizSessionLock(c)
		case util.GetParticipants:
			h.GetParticipants(c)
		case util.SubmitAnswer:
			h.SubmitAnswer(c, mstr.Payload)
		case util.UnsubmitAnswer:
			h.UnsubmitAnswer(c)
		default:
			h.BroadcastMessage(c, mstr)
		}
	}
}
