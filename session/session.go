package session

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"

	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/config"
)

type Sessions struct {
	comms *comms.Comms
	redis *redis.Client

	lifetime time.Duration
}

func New(comms *comms.Comms, client *redis.Client) *Sessions {
	return &Sessions{
		comms: comms,
		redis: client,
		lifetime: time.Duration(config.RedisConfig.TTLDays) * 24 * time.Hour, // One week
	}
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
	err = s.redis.Set(sessionID, 0, s.lifetime).Err()
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

// Assign user to session
func (s *Sessions) Set(sessionID string, userID uint) (err error) {
	return s.redis.Set(sessionID, userID, s.lifetime).Err()
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
	err = s.redis.Expire(sessionID, s.lifetime).Err()
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
	for {
		select {}
	}
}
