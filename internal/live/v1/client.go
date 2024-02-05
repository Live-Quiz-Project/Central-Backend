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
	ID                uuid.UUID `json:"uid"`
	DisplayName       string    `json:"display_name"`
	DisplayEmoji      string    `json:"display_emoji"`
	DisplayColor      string    `json:"display_color"`
	IsHost            bool      `json:"isHost"`
	LiveQuizSessionID uuid.UUID `json:"lqsId"`
	Status            string    `json:"status"`
	Marks             int       `json:"marks"`
}

type Message struct {
	Content           Content   `json:"content"`
	LiveQuizSessionID uuid.UUID `json:"live_quiz_session_id"`
	UserID            uuid.UUID `json:"uid"`
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
		// if c.IsHost {
		// h.hub.Broadcast <- &Message{
		// 	Content: Content{
		// 		Type:    util.EndLQS,
		// 		Payload: "Host has left the session.",
		// 	},
		// 	LiveQuizSessionID: c.LiveQuizSessionID,
		// 	UserID:            c.ID,
		// }
		// h.UnregisterParticipants(c)
		// } else {
		if !c.IsHost {
			_, err := h.Service.UpdateParticipantStatus(context.Background(), c.ID, c.LiveQuizSessionID, util.Left)
			if err != nil {
				log.Printf("Error occured at defer readMessage: %v", err)
			}
		}
		h.hub.Unregister <- c
		// }
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error occured: %v from isHost:%v", err, c.IsHost)
				break
			}
		}

		var mstr Content
		if err := json.Unmarshal(m, &mstr); err != nil {
			log.Printf("Error occured: %v", err)
			break
		}

		switch mstr.Type {
		case util.JoinedLQS, util.LeftLQS:
			h.SendMessage(c, mstr)
		case util.StartLQS:
			h.StartLiveQuizSession(c.LiveQuizSessionID)
		case util.DistQuestion, util.NextQuestion:
			h.DistributeQuestion(c.LiveQuizSessionID)
		case util.DistOptions:
			h.DistributeOptions(c.LiveQuizSessionID)
		case util.RevealAnswer:
			h.RevealAnswer(c.LiveQuizSessionID)
		default:
			h.SendMessage(c, mstr)
		}
	}
}
