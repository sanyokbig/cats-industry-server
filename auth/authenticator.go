package auth

import (
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
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

	tx, err := auth.db.Beginx()
	if err != nil {
		return errors.New("failed to begin tx: " + err.Error())
	}
	defer func() {
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Println("failed to rollback:", rbErr)
			}
			return
		}
		err = tx.Commit()
	}()

	// Create character owning token
	character, err := prepareCharacter(tx, owner)
	if err != nil {
		return err
	}

	// Save token to database
	err = token.Save(tx)
	if err != nil {
		return err
	}

	// If session have logged in user, add new character as an alt to user;
	// login with character otherwise
	if userID != 0 {
		// Session have user, see if user allowed to add alts
		allowed := auth.comms.Sentinel.Check(userID, "add_characters")
		if !allowed {
			return errors.New("not allowed to add new characters")
		}

		// Session have user, assign character as user alt
		err = assignCharacterToUser(tx, character, userID)
		if err != nil {
			return err
		}

	} else {
		// Session have no user. Login with this character
		userID, err = loginWithCharacter(tx, character)
		if err != nil {
			return err
		}

		// Store userID in session
		err = auth.comms.Sessions.Set(state, userID)
		if err != nil {
			return err
		}
	}

	err = notifyClientAboutAuth(tx, state, userID, auth.comms.Hub)
	if err != nil {
		return err
	}

	log.Println("login info sent to", clientID)

	return nil
}

func respError(w io.Writer, err error) {
	w.Write([]byte("something went horribly wrong :(\n\n" + err.Error()))
}
