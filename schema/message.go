package schema

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

//easyjson:json
type Payload map[string]interface{}

func NewPayload() Payload {
	return Payload{}
}

//easyjson:json
type Message struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}

func NewMessage() *Message {
	return &Message{
		Payload: NewPayload(),
	}
}

// Transforms payload map to selected struct via marshaling / unmarshaling
func (v *Payload) Deliver(target json.Unmarshaler) error {
	bytes, err := v.MarshalJSON()
	if err != nil {
		return err
	}

	err = target.UnmarshalJSON(bytes)
	if err != nil {
		return err
	}

	return nil
}

// Packs source data into payload
func (v *Payload) Pack(source json.Marshaler) error {
	bytes, err := source.MarshalJSON()
	if err != nil {
		log.Debugf("1 %v", source)
		return err
	}

	err = v.UnmarshalJSON(bytes)
	if err != nil {
		log.Debugf("2 %s", bytes)
		return err
	}

	return nil
}

func (p *Payload) SetAsDefaultAuthPayload() {
	*p = Payload{
		"user":  nil,
		"roles": &[]string{},
	}
}
