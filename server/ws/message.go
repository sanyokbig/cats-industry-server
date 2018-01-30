package ws

type Payload map[string]interface{}

type Message struct {
	Type    string  `json:"type"`
	Payload Payload `json:"data"`
}