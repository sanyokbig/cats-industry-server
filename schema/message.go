package schema

type Payload map[string]interface{}

type Message struct {
	Type    string  `json:"type"`
	Payload Payload `json:"payload"`
}