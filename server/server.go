package server

import (
	"net/http"
	"log"
)

func Run(port string) {
	http.HandleFunc("/ws", ws)
	http.HandleFunc("/auth", auth)
	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}