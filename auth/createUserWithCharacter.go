package auth

import (
	"github.com/sanyokbig/cats-industry-server/schema"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// Creates user and links him with character
func createUserForCharacter(db sqlx.Queryer, characterID uint) (user *schema.User, err error) {
	user = &schema.User{}
	err = user.Create(db)
	if err != nil {
		err = errors.New("failed to create user: " + err.Error())
		return
	}

	err = user.LinkWithCharacter(db, characterID)
	if err != nil {
		err = errors.New("failed to create user-character link: " + err.Error())
		return
	}
	return
}
