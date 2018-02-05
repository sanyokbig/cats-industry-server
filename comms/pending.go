package comms

type Pending struct {
	Add    chan PendingAdd
	Remove chan string
}

type PendingAdd struct {
	Client, State string
}
