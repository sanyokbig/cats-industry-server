package schema

import (
	"encoding/json"
)

//easyjson:json
type Payload map[string]interface{}

//easyjson:json
type Message struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
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

func (p *Payload) SetAsDefaultAuthPayload() {
	p = &Payload{
		"user":  nil,
		"roles": &[]string{},
	}
}
