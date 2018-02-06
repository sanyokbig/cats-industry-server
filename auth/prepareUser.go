package auth

import (
	"cats-industry-server/schema"
	"database/sql"
	"log"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// If token owner already in system, use login as him
func prepareUser(db *sqlx.DB, token *Token) (user *schema.User, err error) {
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

	owner, err := token.GetOwner()
	if err != nil {
		err = errors.New("failed to get owner: " + err.Error())
		return
	}

	// Find character in db
	character := &schema.Character{}
	err = character.Find(db, owner.CharacterID)
	if err != nil && err != sql.ErrNoRows {
		err = errors.New("failed to find character: " + err.Error())
		return
	}

	// If character not found, create new one
	if err == sql.ErrNoRows {
		character = &schema.Character{
			ID:     owner.CharacterID,
			Name:   owner.CharacterName,
			IsMain: true,
		}

		err = character.Create(tx)

		if err != nil {
			err = errors.New("failed to create character: " + err.Error())
			return
		}
	}

	var shouldCreateLink bool

	// Find user of character
	user = &schema.User{}
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
