package auth

import (
	"cats-industry-server/schema"
	"log"

	"github.com/go-errors/errors"
	"github.com/jmoiron/sqlx"
)

// Creates user and links him with character
func createUserForCharacter(db *sqlx.DB, characterID uint) (user *schema.User, err error) {
	tx, err := db.Beginx()
	if err != nil {
		return nil, errors.New("failed to start tx: " + err.Error())
	}
	defer func() {
		if err != nil {
			log.Println("tx rollrack")
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Println("failed to rollback tx:", rbErr)
			}
			return
		}
		err = tx.Commit()
		if err != nil {
			err = errors.New("failed to commit tx: " + err.Error())
		}
	}()

	err = user.Create(tx)
	if err != nil {
		err = errors.New("failed to create owner: " + err.Error())
		return
	}

	err = user.LinkWithCharacter(tx, characterID)
	if err != nil {
		err = errors.New("failed to create user-character link: " + err.Error())
		return
	}
	return
}
