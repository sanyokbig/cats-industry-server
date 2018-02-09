package auth

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// Rebind all characters from source user to target user and delete former
func assimilateUser(db *sqlx.DB, sourceUser, targetUser uint) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			log.Println("tx rollback")
			rbErr := tx.Rollback()
			if rbErr != nil {
				log.Println("rollback failed:", rbErr)
			}
			return
		}
		err = tx.Commit()
	}()

	// Unset main flag on assimilating characters
	rows, err := tx.Queryx(`
		UPDATE characters SET is_main=false WHERE user_id = $1
	`, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	// Change links
	rows, err = tx.Queryx(`
		UPDATE characters SET user_id = $1 WHERE user_id = $2
	`, targetUser, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	// Delete user
	rows, err = tx.Queryx(`
		DELETE FROM users WHERE id = $1
	`, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}
