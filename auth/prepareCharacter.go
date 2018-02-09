package auth

import (
	"cats-industry-server/schema"
	"database/sql"

	"cats-industry-server/postgres"

	"github.com/go-errors/errors"
)

// Looks for token-owning character, creates if none exists and returns it
func prepareCharacter(db postgres.DB, owner *schema.Owner, userID uint) (character *schema.Character, err error) {
	// Find existing character
	character = &schema.Character{}
	err = character.Find(db, owner.CharacterID)

	// Unexpected error
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
			UserID: userID,
		}

		err = character.Create(db)

		if err != nil {
			err = errors.New("failed to create character: " + err.Error())
			return
		}
	}

	return character, nil
}
