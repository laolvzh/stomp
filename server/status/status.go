package status

type ServerClientSubscriptionStatus struct {
	ID   string
	Dest string
}

type ServerClientStatus struct {
	ID            int64
	Address       string
	Login         string
	Peer          string
	PeerName      string
	Time          string
	Subscriptions []ServerClientSubscriptionStatus
}

type QueueStatus struct {
	Dest              string
	MessageCount      int
	SubscriptionCount int
}

type ServerStatus struct {
	Clients      []ServerClientStatus
	Queues       []QueueStatus
	Topics       []QueueStatus
	Time         string
	Type         string `"json:type"`
	Id           string `"json:id"`
	Name         string `"json:name"`
	Subtype      string `json:"subtype"`
	Subsystem    string `json:"subsystem"`
	ComputerName string `json:"computer"`
	UserName     string `json:"user"`
	ProcessName  string `json:"process"`
	Version      string `json:"version"`
	Pid          int    `json:"pid"`
	Tid          int    `json:"tid"`
	Severity     int    `json:"severity"`
	Message      string `json:"message"`
}
