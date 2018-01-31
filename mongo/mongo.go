package mongo

import (
	"github.com/globalsign/mgo"
	"fmt"
)

type Connection struct {
	*mgo.Session
	uri string
}

func (c *Connection) Connect() error {
	session, err := mgo.Dial(c.uri)
	if err != nil {
		return err
	}
	c.Session = session

	return nil
}

func NewConnection(uri string) *Connection {
	return &Connection{Session: nil, uri: uri}
}

type Dber interface {
	DB(database string) *mgo.Database
}

func GenerateUri(host, port, db, user, pass string) string {
	var auth string
	if user != "" && pass != "" {
		auth = pass + ":" + user + "@"
	}

	return fmt.Sprintf("mongodb://%v%v:%v/%v", auth, host, port, db)
}
