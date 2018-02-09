package schema

import (
	"github.com/jmoiron/sqlx"
)

// Application user, can have multiple characters.
type User struct {
	ID         uint            `json:"id" db:"id"`
	Characters *CharactersList `json:"characters"`
}

// Create new user and insert fresh id in struct
func (u *User) Create(db sqlx.Queryer) error {
	err := db.QueryRowx(`
		INSERT INTO users DEFAULT VALUES RETURNING id
	`).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Find(db sqlx.Queryer, userID uint) error {
	err := db.QueryRowx(`
		SELECT * FROM users WHERE id = $1
	`, userID).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

// Returns user of passed character
func (u *User) FindByCharacter(db sqlx.Queryer, characterID uint) error {
	err := db.QueryRowx(`
		WITH link AS (
			SELECT user_id FROM characters WHERE id = $1
		) SELECT * FROM users WHERE id = (SELECT user_id FROM link)
	`, characterID).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

// Returns user with owned characters joined
func (u *User) FindWithCharacters(db sqlx.Queryer, userID uint) error {
	err := u.Find(db, userID)
	if err != nil {
		return err
	}

	u.Characters = &CharactersList{}
	err = u.Characters.FindByUser(db, userID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) LinkWithCharacter(db sqlx.Queryer, characterID uint) (err error) {
	rows, err := db.Queryx(`
		UPDATE characters SET user_id = $1 WHERE id = $2
	`, u.ID, characterID)
	if err != nil {
		return
	}
	rows.Close()
	return nil
}
