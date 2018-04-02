package sentinel

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/config"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/schema"
	log "github.com/sirupsen/logrus"
)

// Sentinel is responsible for storing and checking user roles
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

// Same as check but receives sessionID instead of userID
func (s *Sentinel) CheckSession(sessionID string, role string) bool {
	userID, err := s.comms.Get(sessionID)
	if err != nil {
		log.Errorf("failed to get sessionID: %v", err)
		return false
	}

	return s.Check(userID, role)
}

// Check if user have required role.
func (s *Sentinel) Check(userID uint, role string) bool {
	key := strconv.Itoa(int(userID))
	err := s.ensureRolesCached(userID)
	if err != nil {
		log.Errorf("failed to ensure roles cached: %v", err)
		return false
	}
	// Check itself
	roles, err := s.redis.SMembersMap(key).Result()
	if err != nil {
		log.Errorf("check failed: %v", err)
		return false
	}
	log.Debug("got roles: %v", roles)

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

	// Stop here if zero roles passed, as an error will occur.
	// Print warning instead as this is not critical, but should not happen
	if len(rs) == 0 {
		log.Warningf("zero roles passed for user %v, cancelling roles set", key)
		return nil
	}

	if err := s.redis.SAdd(key, rs...).Err(); err != nil {
		return err
	}
	// Set ttl
	if err := s.redis.Expire(key, s.lifetime).Err(); err != nil {
		return err
	}

	return nil
}

// Returns list of roles form cache
func (s *Sentinel) GetRoles(userID uint) (*[]string, error) {
	key := strconv.Itoa(int(userID))
	err := s.ensureRolesCached(userID)
	if err != nil {
		log.Errorf("failed to ensure roles cached: %v", err)
		return nil, err
	}

	roles, err := s.redis.SMembers(key).Result()
	if err != nil {
		log.Errorf("check failed: %v", err)
		return nil, err
	}
	return &roles, nil
}

func (s *Sentinel) ensureRolesCached(userID uint) error {
	key := strconv.Itoa(int(userID))
	// First check if key exists. If not, generate roles cache
	exists, err := s.redis.Exists(key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		log.Infof("caching roles for user %v", key)
		err = s.cacheUserRoles(userID)
		if err != nil {
			return err
		}
	}

	return nil
}

// Get users roles from postgres and create entry in redis
func (s *Sentinel) cacheUserRoles(userID uint) error {
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
	log.Info("updating roles cache")
	iter := s.redis.Scan(0, "", 0).Iterator()
	updated, failed := 0, 0
	for iter.Next() {
		userID, err := strconv.Atoi(iter.Val())
		if err != nil {
			failed++
			log.Errorf("failed to parse userID %v: %v", iter.Val(), err)
			continue
		}
		err = s.cacheUserRoles(uint(userID))
		if err != nil {
			failed++
			log.Errorf("failed to cache roles for user %v: %v", userID, err)
			continue
		}
		updated++
	}
	log.Infof("done updating cache: %v updated, %v failed", updated, failed)

}
