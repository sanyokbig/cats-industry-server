package server

import (
	"net/http"
	"log"
	"cats-industry-server/auth"
	"cats-industry-server/postgres"
)

type Server struct {
	Postgres *postgres.Connection
}

func (s *Server) Run(port string) {
	hub := NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		log.Println(query)
		w.Write([]byte("<script>window.close()</script>"))
		token, err := auth.CreateToken(query["code"][0])
		if err != nil {
			log.Println("failed to create token:", err)
		}
		err = token.Save(s.Postgres.DB)
		if err != nil {
			log.Println(err)
		}
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}
