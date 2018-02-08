package session

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"

	"cats-industry-server/comms"
	"cats-industry-server/config"
)

// Session must be deleted after this time
var Lifetime time.Duration

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
	err = s.redis.Set(sessionID, 0, Lifetime).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Assign user to session
func (s *Sessions) Set(sessionID string, userID uint) (err error) {
	return s.redis.Set(sessionID, userID, Lifetime).Err()
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

	// Update ttl when key read
	err = s.redis.Expire(sessionID, Lifetime).Err()
	if err != nil {
		return 0, err
	}

	userIDint, err := strconv.Atoi(result)
	if err != nil {
		return 0, err
	}
	return uint(userIDint), nil
}

func (s *Sessions) Run() {
	Lifetime = time.Duration(config.RedisConfig.TTLDays) * 24 * time.Hour // One week
	for {
		select {}
	}
}
