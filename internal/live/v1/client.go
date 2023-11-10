package v1

import (
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
		h.hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error occured: %v", err)
				break
			}
		}
		log.Printf("Message received: %s", m)

		var mstr Content
		if err := json.Unmarshal(m, &mstr); err != nil {
			log.Printf("Error occured: %v", err)
			break
		}

		switch mstr.Type {
		case util.JoinedLQS, util.LeftLQS:
			h.SendMessage(c, mstr)
		// case util.StartLQS:
		// 	h.startLiveQuizSession(c, mstr.Payload)
		// case util.DistQuestion:
		// 	h.distributeQuestion(mstr.Payload, false)
		// case util.NextQuestion:
		// 	h.distributeQuestion(mstr.Payload, true)
		// case util.DistOptions:
		// 	h.distributeOptions(mstr.Payload, "1", "choice")
		// case util.RevealAnswer:
		// 	h.revealAnswer(mstr.Payload, "1", "choice")
		default:
			h.SendMessage(c, mstr)
		}
	}
}
