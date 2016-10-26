/*
Package server contains a simple STOMP server implementation.
*/
package server

import (
	"net"
	"time"

	"github.com/ventu-io/slf"
)

// The STOMP server has the concept of queues and topics. A message
// sent to a queue destination will be transmitted to the next available
// client that has subscribed. A message sent to a topic will be
// transmitted to all subscribers that are currently subscribed to the
// topic.
//
// Destinations that start with this prefix are considered to be queues.
// Destinations that do not start with this prefix are considered to be topics.
const QueuePrefix = "/queue"
const pwdCurr string = "github.com/go-stomp/stomp/server"

var log slf.StructuredLogger

// Default server parameters.
const (
	// Default address for listening for connections.
	DefaultAddr = ":61613"

	// Default read timeout for heart-beat.
	// Override by setting Server.HeartBeat.
	DefaultHeartBeat = time.Minute
)

var Version string

func init() {
	log = slf.WithContext(pwdCurr)
}

// Interface for authenticating STOMP clients.
type Authenticator interface {
	// Authenticate based on the given login and passcode, either of which might be nil.
	// Returns true if authentication is successful, false otherwise.
	Authenticate(login, passcode string) bool
}

type ServerConfig struct {
	Id         string
	Name       string
	Version    string
	ListenAddr string // TCP address to listen on, DefaultAddr if empty
	Heartbeat  int    //heart-beat interval in seconds
	Status     int    //queue status interval in seconds
	StatusLog  int    //log status interval in seconds
	IsDebug    bool   //log debug data for connections
}

// A Server defines parameters for running a STOMP server.
type Server struct {
	Authenticator Authenticator // Authenticates login/passcodes. If nil no authentication is performed
	QueueStorage  QueueStorage  // Implementation of queue storage. If nil, in-memory queues are used.
	Config        *ServerConfig
}

func (s *Server) Id() string {
	return s.Config.Id
}

func (s *Server) Name() string {
	return s.Config.Name
}

func (s *Server) Version() string {
	return s.Config.Version
}

func (s *Server) StatusDuration() time.Duration {
	return time.Duration(s.Config.Status) * time.Second
}

func (s *Server) StatusLogDuration() time.Duration {
	return time.Duration(s.Config.StatusLog) * time.Second
}

func (s *Server) HeartBeatDuration() time.Duration {
	return time.Duration(s.Config.Heartbeat) * time.Second
}

func NewServer(config *ServerConfig, a Authenticator) *Server {
	log.Infof("NewServer: %+v", config)
	return &Server{
		Config:        config,
		Authenticator: a,
	}
}

// ListenAndServe listens on the TCP network address addr and then calls Serve.
/*func ListenAndServe(addr string, a Authenticator) error {
	s := &Server{Addr: addr, Authenticator: a}
	defer log.Debug("ListenAndServe() function processed")

	return s.ListenAndServe()
}

// Serve accepts incoming TCP connections on the listener l, creating a new
// STOMP service thread for each connection.
func Serve(l net.Listener, a Authenticator) error {
	s := &Server{Authenticator: a, HeartBeat: 15 * time.Second}
	defer log.Debug("Serve() function processed")

	return s.Serve(l)
}*/

// ListenAndServe listens on the TCP network address s.Addr and
// then calls Serve to handle requests on the incoming connections.
// If s.Addr is blank, then DefaultAddr is used.
func (s *Server) ListenAndServe() error {
	addr := s.Config.ListenAddr
	if addr == "" {
		addr = DefaultAddr
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())

	return s.serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a new
// service thread for each connection. The service threads read
// requests and then process each request.

func (s *Server) serve(l net.Listener) error {
	proc := newRequestProcessor(s)
	return proc.Serve(l)
}
