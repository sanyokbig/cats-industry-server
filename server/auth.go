package server

import (
	"net/http"
	"log"
)

func auth(w http.ResponseWriter, r *http.Request) {
	log.Println("auth")
}
