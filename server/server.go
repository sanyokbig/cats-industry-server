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
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		body:= []byte{}
		r.Body.Read(body)

		log.Println(string(body))
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}