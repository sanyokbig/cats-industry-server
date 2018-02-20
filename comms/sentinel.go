package comms

type Sentinel interface {
	Check(userID uint, role string) bool
	SetRoles(userID uint, roles *[]string) error
	WarmUserRoles(userID uint) error
}