package auth

import (
	"cats-industry-server/schema"
	"database/sql"
	"log"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// If token owner already in system, use login as him
func createUserWithCharacter(db *sqlx.DB, character *schema.Character) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		err = errors.New("failed to begin tx: " + err.Error())
		return
	}
	defer func() {
		if err != nil {
			log.Println("tx rollback")
			tx.Rollback()
			return
		}
		err = tx.Commit()
		if err != nil {
			err = errors.New("failed to commit tx: " + err.Error())
		}
	}()

	var shouldCreateLink bool

	// Find user of character
	user := &schema.User{}
	err = user.FindByCharacter(db, character.ID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New("failed to find user: " + err.Error())
		return
	}

	// Create new user if characters user not found
	if err == sql.ErrNoRows {
		shouldCreateLink = true
		err = user.Create(tx)

		if err != nil {
			err = errors.New("failed to create owner: " + err.Error())
			return
		}
	}

	// Create user-character link
	if shouldCreateLink {
		err = user.LinkWithCharacter(tx, character.ID)
		if err != nil {
			err = errors.New("failed to create user-character link: " + err.Error())
			return
		}
	}

	return
}
