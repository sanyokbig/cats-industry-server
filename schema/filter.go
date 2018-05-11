package schema

//easyjson:json
type GetParams struct {
	Filters []Filter `json:"filters"`
	Sort    Sort     `json:"sort"`
	Limit   Limit    `json:"limit"`
}

//easyjson:json
type Filter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

//easyjson:json
type Sort struct {
	Field string `json:"field"`
	Dir   int    `json:"dir"`
}

type Limit int
