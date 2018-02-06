package auth

import (
	"cats-industry-server/schema"
	"log"

	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// If token owner already in system, use login as him
func prepareUser(db *sqlx.DB, token *Token) error {
	owner, err := token.GetOwner()
	if err != nil {
		log.Println(err)
		return errors.New("failed to get owner:")
	}

	// Find character in db
	character := &schema.Character{}
	err = character.Find(db, owner.CharacterID)
	if err != nil && err != sql.ErrNoRows {
		// Unexpected error
		log.Println(err)
		return errors.New("failed to find character")
	}

	// If character not found, create new one
	if err == sql.ErrNoRows {
		character = &schema.Character{
			ID:     owner.CharacterID,
			Name:   owner.CharacterName,
			IsMain: true,
		}
	}

	// Find user of character
	user := &schema.User{}
	err = user.FindByCharacter(db, character.ID)
	if err != nil && err != sql.ErrNoRows {
		// Unexpected error
		log.Println(err)
		return errors.New("failed to find user")
	}

	// Create new user if characters user not found
	if err == sql.ErrNoRows {
		err = user.Create(db)
		if err != nil {
			log.Println(err)
			return errors.New("failed to create user")
		}
	}

	log.Printf("%+v\n", user)

	return nil
}
