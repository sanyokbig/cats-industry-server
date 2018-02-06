package comms

// Used to communicate between processes using channels
type Comms struct {
	Pending
	Sessions
}

func New() *Comms {
	return &Comms{
		Pending: Pending{
			Add:    make(chan PendingAdd),
			Remove: make(chan string),
		},
	}
}
