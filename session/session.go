package session

import (
	"cats-industry-server/comms"
	"log"

	"github.com/satori/go.uuid"
)

// Session must be deleted after this time
const SessionLifetime int64 = 86400 * 7 // One week

type Sessions struct {
	comms *comms.Comms

	// sessionID : userID
	list map[string]uint
}

func New(comms *comms.Comms) *Sessions {
	sessions := &Sessions{
		comms: comms,
		list:  map[string]uint{},
	}

	comms.Sessions = sessions
	return sessions
}

// Create session with no user
func (s *Sessions) Add() (sessionID string) {
	sessionID = uuid.Must(uuid.NewV4()).String()
	s.list[sessionID] = 0

	return sessionID
}

// Assign user to session
func (s *Sessions) Set(sessionID string, userID uint) {
	s.list[sessionID] = userID
	log.Println("setting", sessionID, "to", userID)
}

// Get user of session
func (s *Sessions) Get(sessionID string) uint {
	id, ok := s.list[sessionID]

	if !ok {
		return 0
	}
	return id
}

func (s *Sessions) Run() {
	for {
		select {}
	}
}
