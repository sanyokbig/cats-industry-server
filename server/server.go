package server

import (
	"github.com/sanyokbig/cats-industry-server/auth"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/session"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/sanyokbig/cats-industry-server/sentinel"
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
	sent := sentinel.NewSentinel(c, s.RedisClients.Sessions, s.Postgres)

	c.Hub = hub
	c.Sessions = sessions
	c.Sentinel = sent

	go hub.Run()
	go authenticator.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		authenticator.HandleSSORequest(w, r)
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}
