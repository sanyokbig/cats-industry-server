package server

import (
	"net/http"
	"log"
	"cats-industry-server/server/ws"
	"cats-industry-server/auth"
)

func Run(port string) {
	hub := ws.NewHub()

	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})
	http.HandleFunc("/authRespond", func(w http.ResponseWriter, r *http.Request) {
		query:= r.URL.Query()
		log.Println(query)
		w.Write([]byte("<script>window.close()</script>"))
		auth.VerifyCode(query["code"][0])
	})

	log.Printf("listening on :%v\n", port)
	http.ListenAndServe(":"+port, nil)
}