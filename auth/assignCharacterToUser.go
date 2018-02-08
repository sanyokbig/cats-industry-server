package auth

import (
	"cats-industry-server/schema"

	"log"

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
	// Get character owner
	ownerID, err := character.GetOwnerID(db)
	if err != nil {
		return err
	}

	log.Println("owner", ownerID, "user", userID)

	// Character already owned by this user, nothing to do here
	if ownerID == userID {
		return nil
	}

	if ownerID == 0 {
		// No owner, assigning as planned

		err = character.AssignToUser(db, userID)
		if err != nil {
			return err
		}
		return nil
	}

	// Otherwise, combine users
	err = assimilateUser(db, ownerID, userID)
	if err != nil {
		return err
	}

	return nil
}
