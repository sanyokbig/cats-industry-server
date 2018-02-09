package auth

import (
	"cats-industry-server/postgres"
)

// Rebind all characters from source user to target user and delete former

/*
Example:

Starting users
	User1:
		Alexander
		Elsa
	User2:
		Ansgar
		Fergus

If active session is of User1 and next login done as Fergus, User2 is detected, all its users retrieved and assigned to User2, resulting in:

	User1:
		Alexander
		Elsa
		Ansgar
		Fergus

*/

func assimilateUser(db postgres.DB, sourceUser, targetUser uint) (err error) {
	// Unset main flag on assimilating characters
	rows, err := db.Queryx(`
		UPDATE characters SET is_main=false WHERE user_id = $1
	`, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	// Change links
	rows, err = db.Queryx(`
		UPDATE characters SET user_id = $1 WHERE user_id = $2
	`, targetUser, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	// Delete user
	rows, err = db.Queryx(`
		DELETE FROM users WHERE id = $1
	`, sourceUser)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}
