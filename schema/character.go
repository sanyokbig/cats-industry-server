package schema

import (
	"cats-industry-server/postgres"
	"errors"

	"database/sql"

	"github.com/jmoiron/sqlx"
)

//easyjson:json
type Character struct {
	ID     uint    `json:"id" db:"id"`
	UserID uint    `json:"-" db:"user_id"`
	Name   string  `json:"name" db:"name"`
	IsMain bool    `json:"is_main" db:"is_main"`
	Skills []Skill `json:"-"`
}

//easyjson:json
type CharactersList []Character

//easyjson:json
type Skill struct {
	ID           uint `json:"skill_id"`
	Skillpoints  uint `json:"skillpoints_in_skill"`
	TrainedLevel uint `json:"trained_skill_level"`
	ActiveLevel  uint `json:"active_skill_level"`
}

// Create new character and insert fresh id in struct
func (c *Character) Create(db postgres.NamedQueryer) error {
	rows, err := db.NamedQuery(`
		INSERT INTO characters (id, name, is_main) VALUES (:id, :name, :is_main)
	`, c)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return errors.New("rows.Next failed, could not retrieve id")
	}

	return nil
}

func (c *Character) Find(db sqlx.Queryer, characterID uint) error {
	err := db.QueryRowx(`
		SELECT * FROM characters WHERE id = $1
	`, characterID).StructScan(c)

	if err != nil {
		return err
	}

	return nil
}

func (cl *CharactersList) FindByUser(db sqlx.Queryer, userID uint) error {
	rows, err := db.Queryx(`
		WITH links as (
			SELECT character_id FROM users_characters WHERE user_id = $1
		)
		SELECT * FROM characters WHERE id in (SELECT character_id FROM links)
	`, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	ch := Character{}
	for rows.Next() {
		err = rows.StructScan(&ch)
		if err != nil {
			return err
		}
		*cl = append(*cl, ch)
	}

	return nil
}

func (c *Character) GetOwnerID(db sqlx.Queryer) (userID uint, err error) {
	err = db.QueryRowx(`
		SELECT user_id FROM users_characters WHERE character_id = $1
	`, c.ID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return userID, nil
}

func (c *Character) AssignToUser(db sqlx.Queryer, userID uint) (err error) {
	rows, err := db.Queryx(`
		UPDATE users_characters SET user_id = $1 WHERE character_id = $2
	`, userID, c.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	return nil
}

func (c *Character) UnsetMain(db sqlx.Queryer) (err error) {
	// Change links
	rows, err := db.Queryx(`
		UPDATE characters SET is_main = false WHERE id = $1
	`, c.ID)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}
