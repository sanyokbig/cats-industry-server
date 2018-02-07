package session

import (
	"cats-industry-server/comms"

	"strconv"

	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
)

// Session must be deleted after this time
const SessionLifetime int64 = 86400 * 7 // One week

type Sessions struct {
	comms *comms.Comms
	redis *redis.Client
}

func New(comms *comms.Comms, client *redis.Client) *Sessions {
	sessions := &Sessions{
		comms: comms,
		redis: client,
	}

	comms.Sessions = sessions
	return sessions
}

// Create session with no user
func (s *Sessions) New() (sessionID string, err error) {
	// Generate new SessionID
	newSessionID, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	// Store empty session
	sessionID = newSessionID.String()
	err = s.redis.Set(sessionID, 0, 0).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Assign user to session
func (s *Sessions) Set(sessionID string, userID uint) (err error) {
	return s.redis.Set(sessionID, userID, 0).Err()
}

// Get user of session
func (s *Sessions) Get(sessionID string) (userID uint, err error) {
	result, err := s.redis.Get(sessionID).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	userIDint, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return uint(userIDint), nil
}

func (s *Sessions) Run() {
	for {
		select {}
	}
}
