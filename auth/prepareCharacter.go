package auth

import (
	"github.com/sanyokbig/cats-industry-server/schema"
	"database/sql"

	"github.com/sanyokbig/cats-industry-server/postgres"

	"github.com/go-errors/errors"
)

// Looks for token-owning character, creates if none exists and returns it
func prepareCharacter(db postgres.DB, owner *schema.Owner) (character *schema.Character, err error) {
	// Find existing character
	character = &schema.Character{}
	err = character.Find(db, owner.CharacterID)

	// Unexpected error
	if err != nil && err != sql.ErrNoRows {
		err = errors.New("failed to find character: " + err.Error())
		return
	}

	// If character not found, create ghost one
	if err == sql.ErrNoRows {
		character = &schema.Character{
			ID:     owner.CharacterID,
			Name:   owner.CharacterName,
			IsMain: true,
		}

		err = character.Create(db)
		if err != nil {
			return nil, err
		}
	}

	return character, nil
}
