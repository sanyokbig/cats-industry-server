package auth

import (
	"github.com/sanyokbig/cats-industry-server/schema"

	"database/sql"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// Get user of character or create new one and link with character
func prepareUser(db sqlx.Queryer, characterID uint) (user *schema.User, err error) {
	// Find user of character
	user = &schema.User{}
	err = user.FindByCharacter(db, characterID)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	if err == sql.ErrNoRows {
		// User not found for character, create new one and link with character
		user, err = createUserForCharacter(db, characterID)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
