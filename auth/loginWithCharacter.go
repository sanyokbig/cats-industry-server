package auth

import (
	"github.com/jmoiron/sqlx"
)

func loginWithCharacter(db sqlx.Queryer, character *GhostCharacter) (userID uint, err error) {
	// Prepare user
	user, err := prepareUser(db, character.ID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
