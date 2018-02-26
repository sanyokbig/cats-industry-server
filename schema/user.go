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
			SELECT user_id FROM users_characters WHERE character_id = $1
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
		INSERT INTO users_characters (user_id, character_id) VALUES ($1, $2)
	`, u.ID, characterID)
	if err != nil {
		return
	}
	rows.Close()
	return nil
}

func (u *User) AssignToGroup(db sqlx.Queryer, groupID uint) (err error) {
	rows, err := db.Queryx(`
		INSERT INTO users_groups (user_id, group_id) VALUES ($1, $2)
	`, u.ID, groupID)
	if err != nil {
		return
	}
	rows.Close()
	return nil
}

// Returns string representation of user roles
func (u User) GetRoles(db sqlx.Queryer) (roles *[]string, err error) {
	roles = &[]string{}
	rows, err := db.Queryx(`
		WITH roles_ids AS (
		  WITH groups_ids AS (
		      SELECT group_id
		      FROM users_groups
		      WHERE user_id = $1
		  )
		  SELECT role_id
		  FROM groups_roles
		  WHERE group_id IN (SELECT group_id
		                     FROM groups_ids)
		)
		-- Role names
		SELECT name
		FROM roles
		WHERE id IN (SELECT role_id
		             FROM roles_ids)
	`, u.ID)
	if err != nil {
		return nil, err
	}
	var role string
	for rows.Next() {
		err = rows.Scan(&role)
		*roles = append(*roles, role)
	}

	return roles, nil
}

type rolesGetter interface {
	GetRoles(userID uint) (*[]string, error)
}

// Prepares payload ready to be sent to client. Will not alter receiver.
func (u User) GetAuthPayload(db sqlx.Queryer, getter rolesGetter) (*Payload, error) {
	// Get full user info
	err := u.FindWithCharacters(db, u.ID)
	if err != nil {
		return nil, err
	}

	roles, err := getter.GetRoles(u.ID)
	if err != nil {
		return nil, err
	}

	return &Payload{
		"user":  u,
		"roles": roles,
	}, err
}
