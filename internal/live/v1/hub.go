package v1

import (
	"github.com/Live-Quiz-Project/Backend/internal/util"
	"github.com/google/uuid"
)

type Hub struct {
	LiveQuizSessions map[uuid.UUID]*LiveQuizSession
	Register         chan *Client
	Unregister       chan *Client
	Broadcast        chan *Message
}

func NewHub() *Hub {
	return &Hub{
		LiveQuizSessions: make(map[uuid.UUID]*LiveQuizSession),
		Register:         make(chan *Client),
		Unregister:       make(chan *Client),
		Broadcast:        make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.LiveQuizSessions[cl.LiveQuizSessionID]; ok {
				lqs := h.LiveQuizSessions[cl.LiveQuizSessionID]
				if _, ok := lqs.Clients[cl.ID]; !ok {
					lqs.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.LiveQuizSessions[cl.LiveQuizSessionID]; ok {
				if _, ok := h.LiveQuizSessions[cl.LiveQuizSessionID].Clients[cl.ID]; ok {
					if len(h.LiveQuizSessions[cl.LiveQuizSessionID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content: Content{
								Type:    util.LeftLQS,
								Payload: nil,
							},
							LiveQuizSessionID: cl.LiveQuizSessionID,
							UserID:            cl.ID,
						}
					}
					delete(h.LiveQuizSessions[cl.LiveQuizSessionID].Clients, cl.ID)
					close(cl.Message)
					cl.Conn.Close()
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.LiveQuizSessions[m.LiveQuizSessionID]; ok {
				for _, cl := range h.LiveQuizSessions[m.LiveQuizSessionID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
