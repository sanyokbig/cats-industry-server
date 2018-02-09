package auth

import (
	"cats-industry-server/postgres"
)

/*

Assigns character to user
If character is owned by another user, assimilate old user with new
*/
func assignCharacterToUser(db postgres.DB, character *GhostCharacter, userID uint) (err error) {
	// Get character owner
	ownerID, err := character.GetOwnerID(db)
	if err != nil {
		return err
	}

	// Character already owned by this user, nothing to do here
	if ownerID == userID {
		return nil
	} else if ownerID == 0 {
		// No owner, assigning as planned
		err = character.UnsetMain(db)
		if err != nil {
			return err
		}
		err = character.AssignToUser(db, userID)
		if err != nil {
			return err
		}
		return nil
	} else {
		// Otherwise, combine users
		err = assimilateUser(db, ownerID, userID)
		if err != nil {
			return err
		}
	}

	return nil
}
