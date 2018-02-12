package auth

import (
	"cats-industry-server/comms"
	"cats-industry-server/postgres"
	"log"
	"net/http"

	"io"

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
func (auth *Authenticator) HandleSSORequest(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if err != nil {
			log.Println(err)
			respError(w, err)
		} else {
			w.Write([]byte("<script>window.close()</script>"))
		}
	}()

	// Get query params
	query := r.URL.Query()
	state := query.Get("state")
	code := query.Get("code")
	if state == "" {
		return ErrStateNotFoundInQuery
	}
	if code == "" {
		return ErrCodeNotFoundInQuery
	}

	// Get ws client id
	clientID, ok := auth.pending[state]
	if !ok {
		return ErrUnrecognizedState
	}

	// Create token
	token, err := CreateToken(code)
	if err != nil {
		return err
	}

	// Get owner using token
	owner, err := token.GetOwner()
	if err != nil {
		err = errors.New("failed to get owner: " + err.Error())
		return
	}

	// Get current logged in user
	userID, err := auth.comms.Sessions.Get(state)
	if err != nil {
		return err
	}

	// Create character owning token
	character, err := prepareCharacter(auth.db.DB, owner, userID)
	if err != nil {
		return err
	}
	// Character not saved, save later

	// If session have logged in user, add new character as an alt to user;
	// login with character otherwise
	if userID != 0 {
		// Session have user, assign character as user alt
		err = assignCharacterToUser(auth.db.DB, character, userID)
		if err != nil {
			return err
		}

	} else {
		// Session have no user. Login with this character
		userID, err = loginWithCharacter(auth.db.DB, character)
		if err != nil {
			return err
		}
		err = auth.comms.Sessions.Set(state, userID)
		if err != nil {
			return err
		}
	}

	err = notifyClientAboutAuth(auth.db, state, userID, auth.comms.Hub)
	if err != nil {
		return err
	}

	log.Println("login info sent to", clientID)

	return nil
}

func respError(w io.Writer, err error) {
	w.Write([]byte("something went horribly wrong :(\n\n" + err.Error()))
}
