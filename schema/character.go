package schema

import (
	"github.com/jmoiron/sqlx"
)

type Character struct {
	ID     uint    `db:"id"`
	Name   string  `db:"name"`
	IsMain bool    `db:"is_main"`
	Skills []Skill `db:"skills"`
}

//easyjson:json
type Skill struct {
	ID           uint `json:"skill_id"`
	Skillpoints  uint `json:"skillpoints_in_skill"`
	TrainedLevel uint `json:"trained_skill_level"`
	ActiveLevel  uint `json:"active_skill_level"`
}

func (c *Character) Find(db *sqlx.DB, characterID uint) error {
	err := db.QueryRowx(`
		SELECT * FROM characters WHERE id = $1
	`, characterID).StructScan(c)

	if err != nil {
		return err
	}

	return nil
}
