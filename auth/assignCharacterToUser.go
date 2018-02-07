package auth

import (
	"cats-industry-server/schema"

	"github.com/jmoiron/sqlx"
)

/*

Assigns character to user
If character is owned by another user, assign all characters of that user to a target one and delete old user.

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
func assignCharacterToUser(db *sqlx.DB, character *schema.Character, userID uint) (err error) {
	return nil
}
