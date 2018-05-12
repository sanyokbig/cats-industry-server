package sdeParser

import (
	"io/ioutil"

	"log"

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
	log.Println(out)
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
	log.Println(out)
	return nil
}
