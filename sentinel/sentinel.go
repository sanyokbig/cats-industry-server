package sentinel

import (
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"github.com/sanyokbig/cats-industry-server/config"
	"log"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/schema"
)

type Sentinel struct {
	comms    *comms.Comms
	redis    *redis.Client
	postgres *postgres.Connection
	lifetime time.Duration
}

func NewSentinel(comms *comms.Comms, redis *redis.Client, postgres *postgres.Connection) *Sentinel {
	return &Sentinel{
		comms:    comms,
		redis:    redis,
		postgres: postgres,
		lifetime: time.Duration(config.RedisConfig.TTLDays) * 24 * time.Hour, // One week
	}
}

// Check if user have required role.
func (s *Sentinel) Check(userID uint, role string) bool {
	key := strconv.Itoa(int(userID))

	roles, err := s.redis.SMembersMap(key).Result()
	if err != nil {
		log.Println("check failed:", err)
		return false
	}

	log.Println(roles)

	return false
}

// Add roles to existing ones
func (s *Sentinel) AddRoles(userID uint, roles []string) error {
	return s.redis.SAdd(strconv.Itoa(int(userID)), roles).Err()
}

// Remove some roles
func (s *Sentinel) RemoveRoles(userID uint, roles []string) error {
	return s.redis.SRem(strconv.Itoa(int(userID)), roles).Err()
}

// Set passed roles to user regardless of previous ones
func (s *Sentinel) SetRoles(userID uint, roles []string) error {
	key := strconv.Itoa(int(userID))
	// Drop current roles
	if err := s.redis.Del(key).Err(); err != nil {
		return err
	}

	// Add new ones
	if err := s.redis.SAdd(key, roles).Err(); err != nil {
		return err
	}

	// Set ttl
	if err := s.redis.Expire(key, s.lifetime).Err(); err != nil {
		return err
	}

	return nil

}

// Get users roles from postgres and create entry in redis
func (s *Sentinel) PrepareUser(userID uint) error {
	roles, err := schema.User{ID: userID}.GetRoles(s.postgres)
	if err!= nil {
		return err
	}
	log.Println(roles)

	return nil
}
