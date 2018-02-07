package schema

import (
	"cats-industry-server/postgres"

	"github.com/jmoiron/sqlx"
)

// Application user, can have multiple characters.
type User struct {
	ID uint `db:"id"`
}

// Create new user and insert fresh id in struct
func (u *User) Create(db postgres.QueryRowxer) error {
	err := db.QueryRowx(`
		INSERT INTO users DEFAULT VALUES RETURNING id
	`).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) Find(db postgres.QueryRowxer, userID uint) error {
	err := db.QueryRowx(`
		SELECT * FROM users WHERE id = $1
	`, userID).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) FindByCharacter(db postgres.QueryRowxer, characterID uint) error {
	err := db.QueryRowx(`
		WITH link AS (
			SELECT user_id FROM users_characters WHERE character_id = $1
		) SELECT * FROM users WHERE id = (SELECT user_id FROM link)
	`, characterID).StructScan(u)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) LinkWithCharacter(db sqlx.Queryer, characterID uint) (err error) {
	rows, err := db.Queryx(`
		INSERT INTO users_characters (user_id, character_id) VALUES ($1, $2)
	`, u.ID, characterID)
	if err != nil {
		return
	}
	rows.Close()
	return nil
}
