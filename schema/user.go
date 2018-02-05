package schema

import "github.com/jmoiron/sqlx"

type User struct {
	Id uint `db:"id"`
}

func (u *User) FindByCharacter(db *sqlx.DB, characterID uint) error {
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
