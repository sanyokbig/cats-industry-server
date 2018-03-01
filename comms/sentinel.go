package comms

type Sentinel interface {
	Check(userID uint, role string) bool
	GetRoles(userID uint) (*[]string, error)
}
