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

	// First check if key exists. If not, generate roles cache
	exists, err := s.redis.Exists(key).Result()
	if err != nil {
		log.Println("failed exist check:", err)
		return false
	}
	if exists == 0 {
		err = s.WarmUserRoles(userID)
		if err != nil {
			log.Println("failed to warm roles", err)
			return false
		}
	}

	// Check itself
	roles, err := s.redis.SMembersMap(key).Result()
	if err != nil {
		log.Println("check failed:", err)
		return false
	}
	log.Println(roles)

	_, ok := roles[role]
	return ok
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
func (s *Sentinel) SetRoles(userID uint, roles *[]string) error {
	key := strconv.Itoa(int(userID))
	// Drop current roles
	if err := s.redis.Del(key).Err(); err != nil {
		return err
	}

	// Add new ones
	rs := []interface{}{}
	for _, r := range *roles {
		rs = append(rs, r)
	}

	if len(rs) != 0 {
		if err := s.redis.SAdd(key, rs...).Err(); err != nil {
			return err
		}
	}

	// Set ttl
	if err := s.redis.Expire(key, s.lifetime).Err(); err != nil {
		return err
	}

	return nil
}

// Get users roles from postgres and create entry in redis
func (s *Sentinel) WarmUserRoles(userID uint) error {
	roles, err := schema.User{ID: userID}.GetRoles(s.postgres)
	if err != nil {
		return err
	}
	err = s.SetRoles(userID, roles)
	if err != nil {
		return err
	}
	return nil
}

// Runs through all stored user keys and updates roles list
func (s *Sentinel) UpdateCache() {
	iter := s.redis.Scan(0, "", 0).Iterator()
	for iter.Next() {
		userID, err := strconv.Atoi(iter.Val())
		if err != nil {
			log.Printf("failed to parse userID%v: %v", iter.Val(), err)
			continue
		}
		err = s.WarmUserRoles(uint(userID))
		if err != nil {
			log.Printf("failed to warm user roles %v: %v", userID, err)
			continue
		}
	}
}
