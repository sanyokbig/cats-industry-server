package comms

// Used to communicate between processes using channels
type Comms struct {
	Pending
	Sessions
	Hub
}

func New() *Comms {
	return &Comms{
		Pending: Pending{
			Add:    make(chan PendingAdd),
			Remove: make(chan string),
		},
	}
}
