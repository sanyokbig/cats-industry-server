package schema

import (
	"github.com/sanyokbig/cats-industry-server/postgres"

	"database/sql"

	"github.com/jmoiron/sqlx"
	"strconv"
)

//easyjson:json
type Character struct {
	ID     uint    `json:"id" db:"id"`
	Name   string  `json:"name" db:"name"`
	IsMain bool    `json:"is_main" db:"is_main"`
	Skills []Skill `json:"-"`

	IsMailing    bool `json:"is_mailing"`
	IsIndustrial bool `json:"is_industrial"`
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
		INSERT INTO users_characters (user_id, character_id) VALUES ($1, $2)
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

type tokenResult struct {
	isMailing    bool
	isIndustrial bool
}

// Define what type of token character provided.
func (cl *CharactersList) LoadTokenStatus(db sqlx.Queryer) (err error) {
	results := map[uint]*tokenResult{}
	// Get all char ids
	charIds := ""
	for _, c := range *cl {
		charIds += strconv.Itoa(int(c.ID))
		results[c.ID] = &tokenResult{}
	}

	// Load all tokens for these characters
	rows, err := db.Queryx(`
		SELECT character_id, scopes FROM tokens WHERE character_id in ($1)
		`, charIds)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Run through result and set flags
	for rows.Next() {
		c, sc := uint(0), ""
		err = rows.Scan(&c, &sc)
		if err != nil {
			return err
		}
		// Get type of token scope set
		name, ok := ScopeSetsReversed[sc]
		if ok {
			if name == "industrial" {
				(results[c]).isIndustrial = true
			} else if name == "mailing" {
				(results[c]).isMailing = true
			}
		}
	}

	for i, char := range *cl {
		r, ok := results[char.ID]
		if ok {
			char.IsIndustrial = r.isIndustrial
			char.IsMailing = r.isMailing
		}
		(*cl)[i] = char
	}

	return nil
}
