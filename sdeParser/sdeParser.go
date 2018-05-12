package sdeParser

import (
	"io/ioutil"

	"fmt"

	"strings"

	"github.com/go-errors/errors"
	"github.com/sanyokbig/cats-industry-server/postgres"
	"gopkg.in/yaml.v2"
)

type SdeImporter struct {
	postgres *postgres.Connection
}

func NewSdeImporter(connection *postgres.Connection) *SdeImporter {
	return &SdeImporter{
		postgres: connection,
	}
}

func (p *SdeImporter) ImportActivities(sdePath string) error {
	bytes, err := ioutil.ReadFile(sdePath)
	if err != nil {
		return err
	}
	out := Activities{}
	err = yaml.Unmarshal(bytes, &out)
	if err != nil {
		return err
	}
	values := ``
	for _, a := range out {
		values += fmt.Sprintf(`(%v,'%v','%v','%v'),`, a.ID, escape(a.Name), escape(a.Description), a.Icon)
	}
	if len(values) == 0 {
		return errors.New("no values to insert")
	}
	values = values[:len(values)-1]
	query := fmt.Sprintf(`
		insert into ram_activities (id, name, description, icon) 
		values %v on conflict (id) do update set 
			name = excluded.name, 
			description = excluded.description, 
			icon = excluded.icon`, values)
	_, err = p.postgres.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (p *SdeImporter) ImportProductTypes(sdePath string) error {
	bytes, err := ioutil.ReadFile(sdePath)
	if err != nil {
		return err
	}
	out := ProductTypes{}
	err = yaml.Unmarshal(bytes, &out)
	if err != nil {
		return err
	}

	values := ``
	for id, t := range out {
		values += fmt.Sprintf(`(%v,'%v'),`, id, escape(t.Name.En))
	}
	if len(values) == 0 {
		return errors.New("no values to insert")
	}
	values = values[:len(values)-1]
	query := fmt.Sprintf(`
		insert into product_types (id, name) 
		values %v on conflict (id) do update set
		name = excluded.name`, values)
	_, err = p.postgres.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func escape(s string) string {
	return strings.Replace(s, "'", "''", -1)
}
