package schema

type TokenType uint8

const (
	Char TokenType = iota + 1
	Corp
)

type Token struct {
	Id           uint
	UserId       uint
	Type         TokenType
	Expires      string
	AccessToken  string
	RefreshToken string
}
