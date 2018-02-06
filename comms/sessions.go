package comms

type Sessions interface {
	Add() (sessionID string)
	Set(sessionID string, userID uint)
	Get(sessionID string) uint
}
