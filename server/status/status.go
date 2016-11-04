package status

type ServerClientSubscriptionStatus struct {
	ID   string
	Dest string
}

type ServerClientStatus struct {
	ID                    int64
	Address               string
	Login                 string
	Peer                  string
	PeerName              string
	Time                  string
	SkippedWrites         int64
	CurrentSkippedWrites  int
	SentFrames            int64
	CurrentSentFrames     int
	ReceivedFrames        int64
	CurrentReceivedFrames int
	Subscriptions         []ServerClientSubscriptionStatus
}

type QueueStatus struct {
	Dest              string
	MessageCount      int
	TotalCount        int64
	CurrentCount      int
	SubscriptionCount int
}

type TopicStatus struct {
	Dest              string
	TotalCount        int64
	CurrentCount      int
	SubscriptionCount int
}

type ServerStatus struct {
	Clients                   []*ServerClientStatus
	Queues                    []*QueueStatus
	Topics                    []*TopicStatus
	Time                      string  `json:"utc"`
	Type                      string  `json:"type"`
	Id                        string  `json:"id"`
	Name                      string  `json:"name"`
	Subtype                   string  `json:"subtype"`
	Subsystem                 string  `json:"subsystem"`
	ComputerName              string  `json:"computer"`
	UserName                  string  `json:"user"`
	ProcessName               string  `json:"process"`
	Version                   string  `json:"version"`
	Pid                       int     `json:"pid"`
	Tid                       int     `json:"tid"`
	Severity                  int     `json:"severity"`
	Message                   string  `json:"message"`
	EnqueueCount              int     `json:"enqueueCount"`
	RequeueCount              int     `json:"requeueCount"`
	ConnectCount              int     `json:"connectCount"`
	DisconnectCount           int     `json:"disconnectCount"`
	CurrentEnqueueCount       int     `json:"currentEnqueueCount"`
	CurrentRequeueCount       int     `json:"currentRequeueCount"`
	CurrentConnectCount       int     `json:"currentConnectCount"`
	CurrentDisconnectCount    int     `json:"currentDisconnectCount"`
	TotalCurrentCount         int     `json:"totalCurrentCount"`
	TotalQueueCount           int     `json:"totalQueueCount"`
	TotalCurrentSkippedWrites int     `json:"totalSkippedWrites"`
	MessageRate               float64 `json:"messageRate"`
}
