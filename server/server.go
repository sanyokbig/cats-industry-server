package server

import (
	"cats-industry-server/auth"
	"cats-industry-server/comms"
	"cats-industry-server/postgres"
	"log"
	"net/http"
)

type Server struct {
	Postgres *postgres.Connection
}

func (s *Server) Run(port string) {
	c := comms.New()

	hub := NewHub(c)
	authenticator := auth.New(c, s.Postgres)

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
