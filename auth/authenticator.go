package auth

import (
	"net/http"

	"io"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
	"github.com/sanyokbig/cats-industry-server/comms"
	"github.com/sanyokbig/cats-industry-server/postgres"
	log "github.com/sirupsen/logrus"
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

func rollback(tx *sqlx.Tx) {
	rbErr := tx.Rollback()
	if rbErr != nil {
		log.Errorf("failed to rollback: %v", rbErr)
	}
}

// Create token using passed code, get owner info, create new character and user if needed
func (auth *Authenticator) HandleSSORequest(w http.ResponseWriter, r *http.Request) (err error) {
	defer func() {
		if err != nil {
			log.Errorf("failed to handle sso request: %v", err)
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

	// Create character owning token
	character, err := prepareCharacter(tx, owner)
	if err != nil {
		rollback(tx)
		return err
	}

	// Save token to database
	err = token.Save(tx)
	if err != nil {
		rollback(tx)
		return err
	}

	// If session have logged in user, add new character as an alt to user;
	// login with character otherwise
	if userID != 0 {
		// Session have user, assign character as user alt
		err = assignCharacterToUser(tx, character, userID)
		if err != nil {
			rollback(tx)
			return err
		}

	} else {
		// Session have no user. Login with this character
		userID, err = loginWithCharacter(tx, character)
		if err != nil {
			rollback(tx)
			return err
		}

		// Store userID in session
		err = auth.comms.Sessions.Set(state, userID)
		if err != nil {
			rollback(tx)
			return err
		}
	}

	// Commit tx
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commig tx: %v", err)
		return errors.New("failed to commit tx")
	}

	err = notifyClientAboutAuth(auth.db, state, userID, auth.comms.Hub, auth.comms.Sentinel)
	if err != nil {
		return err
	}

	log.Infof("login info sent to %v", clientID)

	return nil
}

func respError(w io.Writer, err error) {
	w.Write([]byte("something went horribly wrong :(\n\n" + err.Error()))
}
