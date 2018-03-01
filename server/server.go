package server

import (
	"net/http"

	"github.com/go-redis/redis"
	"github.com/sanyokbig/cats-industry-server/auth"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/sentinel"
	"github.com/sanyokbig/cats-industry-server/session"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	Postgres     *postgres.Connection
	RedisClients *RedisClients
}

type RedisClients struct {
	Sessions *redis.Client
	Roles    *redis.Client
}

func (s *Server) Run(port string) {
	c := comms.New()

	hub := NewHub(c, s.Postgres)
	authenticator := auth.New(c, s.Postgres)
	sessions := session.New(c, s.RedisClients.Sessions)
	sent := sentinel.NewSentinel(c, s.RedisClients.Roles, s.Postgres)

	c.Hub = hub
	c.Sessions = sessions
	c.Sentinel = sent

	// Start accepting WS connections
	go hub.Run()

	// Start storing auth requests
	go authenticator.Run()

	// Update cached roles in case of changed groups/roles in database
	go sent.UpdateCache()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		authenticator.HandleSSORequest(w, r)
	})

	log.Infof("listening on :%v", port)
	http.ListenAndServe(":"+port, nil)
}
