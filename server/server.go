package server

import (
	"github.com/sanyokbig/cats-industry-server/auth"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"github.com/sanyokbig/cats-industry-server/session"
	"log"
	"net/http"

	"github.com/go-redis/redis"
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

	go hub.Run()
	go authenticator.Run()
	go sessions.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		authenticator.HandleSSORequest(w, r)
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}
