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
	Clients []ServerClientStatus
	Queues  []QueueStatus
	Topics  []QueueStatus
	Time    string
}
