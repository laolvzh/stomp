package client

import (
	"time"
)

// Contains information the client package needs from the
// rest of the STOMP server code.
type Config interface {
	// Method to authenticate a login and associated passcode.
	// Returns true if login/passcode is valid, false otherwise.
	Authenticate(login, passcode string) bool

	// Default duration for read/write heart-beat values. If this
	// returns zero, no heart-beat will take place. If this value is
	// larger than the maximu permitted value (which is more than
	// 11 days, but less than 12 days), then it is truncated to the
	// maximum permitted values.
	HeartBeat() time.Duration

	// show various debug information for connections
	IsDebug() bool

	// Maximum number of pending frames allowed before the read
	// go routine starts blocking.
	MaxPendingReads() int

	// Maximum number of pending frames allowed to a client.
	// before a disconnect occurs. If the client cannot keep
	// up with the server, we do not want the server to backlog
	// pending frames indefinitely.
	MaxPendingWrites() int
}
