package server

import (
	"net/http"
	"log"
	"cats-industry-server/server/ws"
)

func Run(port string) {
	hub := ws.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}