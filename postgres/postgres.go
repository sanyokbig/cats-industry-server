package postgres

import (
	"github.com/jmoiron/sqlx"

	"fmt"

	_ "github.com/lib/pq"
)

type Connection struct {
	*sqlx.DB
	connStr string
}

func (c *Connection) Connect() error {
	conn, err := sqlx.Open("postgres", c.connStr)
	if err != nil {
		return err
	}

	err = conn.Ping()
	if err != nil {
		conn.Close()
		return err
	}
	c.DB = conn

	return nil
}

func NewConnection(host, port, db, user, password string) *Connection {
	connStr := fmt.Sprintf("host='%s' port='%s' dbname='%s' user='%s' password='%s' sslmode=disable", host, port, db, user, password)
	return &Connection{DB: new(sqlx.DB), connStr: connStr}
}

type NamedQueryer interface {
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

type DB interface {
	NamedQueryer
	sqlx.Queryer
}
