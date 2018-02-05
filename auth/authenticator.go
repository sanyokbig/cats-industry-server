package auth

import (
	"cats-industry-server/comms"
	"cats-industry-server/postgres"
	"log"
	"net/http"
)

type Authenticator struct {
	comms *comms.Comms

	pending map[string]string
	db      *postgres.Connection
	add     chan string
	remove  chan string
}

type a interface {
}

func New(comms *comms.Comms, conn *postgres.Connection) *Authenticator {
	return &Authenticator{
		comms:   comms,
		pending: map[string]string{},
		db:      conn,
	}
}

func (a *Authenticator) Run() {
	for {
		select {
		case toAdd := <-a.comms.Pending.Add:
			{
				a.pending[toAdd.State] = toAdd.Client
			}
		case toRemove := <-a.comms.Pending.Remove:
			{
				delete(a.pending, toRemove)
			}
		}
		log.Println(a.pending)
	}
}

func (a *Authenticator) HandleSSORequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query["state"][0]
	client, ok := a.pending[state]
	if !ok {
		log.Println("unrecognized state")
		w.Write([]byte("something went horribly wrong :(\nunrecognized state"))
		return
	}

	log.Println(client)

	w.Write([]byte("<script>window.close()</script>"))
	token, err := CreateToken(query["code"][0])
	if err != nil {
		log.Println("failed to create token:", err)
	}
	err = token.Save(a.db.DB)
	if err != nil {
		log.Println(err)
	}
}
