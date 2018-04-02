package server

import (
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	log "github.com/sirupsen/logrus"

	"github.com/sanyokbig/cats-industry-server/schema"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	comms    *comms.Comms
	postgres *postgres.Connection
	// Sessions of registered clients.
	sessions map[string]map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub(comms *comms.Comms, connection *postgres.Connection) *Hub {
	hub := &Hub{
		comms:      comms,
		postgres:   connection,
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		sessions:   make(map[string]map[*Client]bool),
	}

	comms.Hub = hub

	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.sessions[client.sid] == nil {
				h.sessions[client.sid] = map[*Client]bool{}
			}
			h.sessions[client.sid][client] = true

		case client := <-h.unregister:
			session := h.sessions[client.sid]
			if session == nil {
				continue
			}

			if _, ok := session[client]; ok {
				close(client.send)
				delete(session, client)
			}

			if len(session) == 0 {
				delete(h.sessions, client.sid)
			}

		case message := <-h.broadcast:
			for _, session := range h.sessions {
				for client := range session {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(session, client)
					}
				}
			}
		}

		for k, v := range h.sessions {
			log.Debugf("s: %v", k)
			for s := range v {
				log.Debugf("		c: %v", s.id)
			}
		}
	}
}

func (h *Hub) SendToSession(session string, message *schema.Message) {
	clients := h.sessions[session]

	for c := range clients {
		c.Respond(message)
	}
}
