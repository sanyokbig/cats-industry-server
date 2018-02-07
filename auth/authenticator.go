package auth

import (
	"cats-industry-server/comms"
	"cats-industry-server/postgres"
	"log"
	"net/http"

	"github.com/go-errors/errors"
)

var (
	ErrStateNotFoundInQuery = errors.New("state not found in request query")
	ErrCodeNotFoundInQuery  = errors.New("code not found in request query")
	ErrUnrecognizedState    = errors.New("unrecognized state")
)

type Authenticator struct {
	comms *comms.Comms
	db    *postgres.Connection

	pending map[string]string
	add     chan string
	remove  chan string
}

func New(comms *comms.Comms, conn *postgres.Connection) *Authenticator {
	return &Authenticator{
		comms:   comms,
		pending: map[string]string{},
		db:      conn,
	}
}

func (auth *Authenticator) Run() {
	for {
		select {
		case toAdd := <-auth.comms.Pending.Add:
			{
				auth.pending[toAdd.State] = toAdd.Client
			}
		case toRemove := <-auth.comms.Pending.Remove:
			{
				delete(auth.pending, toRemove)
			}
		}
	}
}

// Create token using passed code, get owner info, create new character and user if needed
func (auth *Authenticator) HandleSSORequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state := query.Get("state")
	code := query.Get("code")
	if state == "" {
		log.Println(ErrStateNotFoundInQuery)
		w.Write([]byte("something went horribly wrong :(\n\n" + ErrStateNotFoundInQuery.Error()))

		return
	}
	if code == "" {
		log.Println(ErrCodeNotFoundInQuery)
		w.Write([]byte("something went horribly wrong :(\n\n" + ErrCodeNotFoundInQuery.Error()))

		return
	}

	_, ok := auth.pending[state]
	if !ok {
		log.Println(ErrUnrecognizedState)
		w.Write([]byte("something went horribly wrong :(\n\n" + ErrUnrecognizedState.Error()))

		return
	}

	// Create token
	token, err := CreateToken(code)
	if err != nil {
		log.Println("failed to create token:", err)
	}

	user, err := prepareUser(auth.db.DB, token)
	if err != nil {
		log.Println(err)
		w.Write([]byte("something went horribly wrong :(\n\n" + err.Error()))
		return
	}

	err = auth.comms.Sessions.Set(state, user.ID)
	if err != nil {
		w.Write([]byte("something went horribly wrong :(\n\n" + err.Error()))
	}

	w.Write([]byte("<script>window.close()</script>"))

}
