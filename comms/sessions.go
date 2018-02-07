package comms

type Sessions interface {
	New() (sessionID string, err error)
	Set(sessionID string, userID uint) error
	Get(sessionID string) (userID uint, err error)
}
