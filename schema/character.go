package schema

import (
	"cats-industry-server/postgres"

	"errors"
)

type Character struct {
	ID     uint   `db:"id"`
	Name   string `db:"name"`
	IsMain bool   `db:"is_main"`
	Skills []Skill
}

//easyjson:json
type Skill struct {
	ID           uint `json:"skill_id"`
	Skillpoints  uint `json:"skillpoints_in_skill"`
	TrainedLevel uint `json:"trained_skill_level"`
	ActiveLevel  uint `json:"active_skill_level"`
}

// Create new character and insert fresh id in struct
func (c *Character) Create(db postgres.NamedQuerier) error {
	rows, err := db.NamedQuery(`
		INSERT INTO characters (id, name, is_main) VALUES (:id, :name, :is_main) RETURNING id
	`, c)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("rows.Next failed, could not retrieve id")
	}

	err = rows.StructScan(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Character) Find(db postgres.QueryRowxer, characterID uint) error {
	err := db.QueryRowx(`
		SELECT * FROM characters WHERE id = $1
	`, characterID).StructScan(c)

	if err != nil {
		return err
	}

	return nil
}

func (c *Character) FindByUser(db postgres.QueryRowxer, userID uint) error {
	err := db.QueryRowx(`
		WITH link AS (
			SELECT character_id FROM users_characters WHERE user_id = $1
		) SELECT * FROM characters WHERE id = (SELECT character_id FROM link)
	`, userID).StructScan(c)

	if err != nil {
		return err
	}

	return nil
}
