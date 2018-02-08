package auth

import (
	"cats-industry-server/schema"

	"log"

	"github.com/jmoiron/sqlx"
)

func loginWithCharacter(db *sqlx.DB, character *schema.Character) (err error) {
	// Prepare user
	user, err := prepareUser(db, character.ID)
	if err != nil {
		return err
	}

	log.Println(user)

	return nil
}
