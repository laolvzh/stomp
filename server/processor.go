package server

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-stomp/stomp/frame"
	"github.com/go-stomp/stomp/server/client"
	"github.com/go-stomp/stomp/server/queue"
	"github.com/go-stomp/stomp/server/status"
	"github.com/go-stomp/stomp/server/topic"
)

//var log slf.StructuredLogger

//func init() {
//log = slf.WithContext("processor")
//}

type requestProcessor struct {
	server                 *Server
	config                 *config
	ch                     chan client.Request
	tm                     *topic.Manager
	qm                     *queue.Manager
	connections            map[int64]*client.Conn
	stop                   bool // has stop been requested
	connectCount           int
	disconnectCount        int
	enqueueCount           int
	requeueCount           int
	currentConnectCount    int
	currentDisconnectCount int
	currentEnqueueCount    int
	currentRequeueCount    int

	currentEnqueueCountLog int
	currentQueueCountLog   int
	currentSkippedCount    int
}

func newRequestProcessor(server *Server) *requestProcessor {
	config := newConfig(server)
	proc := &requestProcessor{
		server:      server,
		config:      config,
		ch:          make(chan client.Request, config.MaxPendingWrites()*16), //HACK: arbitrary coeff
		tm:          topic.NewManager(),
		connections: make(map[int64]*client.Conn),
	}

	if server.QueueStorage == nil {
		proc.qm = queue.NewManager(queue.NewMemoryQueueStorage())
	} else {
		proc.qm = queue.NewManager(server.QueueStorage)
	}

	return proc
}

func (proc *requestProcessor) createStatus() *status.ServerStatus {
	//clients
	totalCurrentSkippedWrites := 0
	clients := make([]*status.ServerClientStatus, 0)
	for _, conn := range proc.connections {
		connStatus := conn.GetStatus()
		totalCurrentSkippedWrites += connStatus.CurrentSkippedWrites
		clients = append(clients, connStatus)
	}

	//
	totalQueueCount := 0
	totalCurrentCount := 0
	queues := proc.qm.GetStatus()

	for _, qs := range queues {
		totalQueueCount += qs.MessageCount
		totalCurrentCount += qs.CurrentCount
	}

	//
	topics := proc.tm.GetStatus()
	for _, ts := range topics {
		totalCurrentCount += ts.CurrentCount
	}

	hostname, _ := os.Hostname()

	rate := float64(proc.currentEnqueueCount+proc.currentRequeueCount) / float64(proc.server.Config.Status)

	serverStatus := &status.ServerStatus{
		Clients:                   clients,
		Queues:                    queues,
		Topics:                    topics,
		Time:                      time.Now().Format(time.RFC3339),
		Type:                      "status",
		Id:                        proc.server.Id(),
		Name:                      proc.server.Name(),
		Version:                   proc.server.Version(),
		Subtype:                   "server",
		Subsystem:                 "",
		ComputerName:              hostname,
		UserName:                  fmt.Sprintf("%d", os.Getuid()),
		ProcessName:               os.Args[0],
		Pid:                       os.Getpid(),
		Severity:                  20,
		EnqueueCount:              proc.enqueueCount,
		RequeueCount:              proc.requeueCount,
		ConnectCount:              proc.connectCount,
		DisconnectCount:           proc.disconnectCount,
		CurrentEnqueueCount:       proc.currentEnqueueCount,
		CurrentRequeueCount:       proc.currentRequeueCount,
		CurrentConnectCount:       proc.currentConnectCount,
		CurrentDisconnectCount:    proc.currentDisconnectCount,
		TotalQueueCount:           totalQueueCount,
		TotalCurrentCount:         totalCurrentCount,
		TotalCurrentSkippedWrites: totalCurrentSkippedWrites,
		MessageRate:               rate,
	}

	//diffs
	proc.currentEnqueueCountLog += proc.currentEnqueueCount
	proc.currentSkippedCount += totalCurrentSkippedWrites

	//state
	proc.currentQueueCountLog = totalQueueCount

	proc.currentEnqueueCount = 0
	proc.currentRequeueCount = 0
	proc.currentConnectCount = 0
	proc.currentDisconnectCount = 0

	return serverStatus
}

func (proc *requestProcessor) createStatusFrame() *frame.Frame {
	f := frame.New("MESSAGE", frame.ContentType, "application/json")
	status := proc.createStatus()
	//bytes, err := json.MarshalIndent(status, "", "  ")
	bytes, err := json.Marshal(status)
	//log.Debugf("createStatusFrame %v", string(bytes))
	if err != nil {
		f.Body = []byte(fmt.Sprintf("error %v\n", err))
	} else {
		f.Body = bytes
	}
	return f
}

func (proc *requestProcessor) sendStatusFrame() {
	topic := proc.tm.Find("/topic/go-stomp.status")
	f := proc.createStatusFrame()
	f.Header.Add(frame.Destination, "/topic/go-stomp.status")
	//log.Debugf("status frame %v", f.Dump())
	topic.Enqueue(f)
}

func (proc *requestProcessor) Serve(l net.Listener) error {
	go proc.Listen(l)

	ticker := time.NewTicker(proc.server.StatusDuration())
	infoTicker := time.NewTicker(proc.server.StatusLogDuration())

	for {
		select {
		case _ = <-infoTicker.C:

			log.Debugf("status: processed:%d total:%d skipped:%d clients:%d",
				proc.currentEnqueueCountLog, proc.currentQueueCountLog, proc.currentSkippedCount, len(proc.connections))
			proc.currentEnqueueCountLog = 0
			proc.currentQueueCountLog = 0
			proc.currentSkippedCount = 0
		case _ = <-ticker.C:
			proc.sendStatusFrame()
		case r := <-proc.ch:
			switch r.Op {
			case client.SubscribeOp:
				if isQueueDestination(r.Sub.Destination()) {
					queue := proc.qm.Find(r.Sub.Destination())
					// todo error handling
					queue.Subscribe(r.Sub)
				} else {
					topic := proc.tm.Find(r.Sub.Destination())
					topic.Subscribe(r.Sub)
				}

			case client.UnsubscribeOp:
				if isQueueDestination(r.Sub.Destination()) {
					queue := proc.qm.Find(r.Sub.Destination())
					// todo error handling
					queue.Unsubscribe(r.Sub)
				} else {
					topic := proc.tm.Find(r.Sub.Destination())
					topic.Unsubscribe(r.Sub)
				}

			case client.EnqueueOp:
				destination, ok := r.Frame.Header.Contains(frame.Destination)
				if !ok {
					// should not happen, already checked in lower layer
					panic("missing destination")
				}
				proc.enqueueCount++
				proc.currentEnqueueCount++

				if isQueueDestination(destination) {
					queue := proc.qm.Find(destination)
					queue.Enqueue(r.Frame)
				} else {
					topic := proc.tm.Find(destination)
					topic.Enqueue(r.Frame)
				}

			case client.RequeueOp:
				destination, ok := r.Frame.Header.Contains(frame.Destination)
				if !ok {
					// should not happen, already checked in lower layer
					panic("missing destination")
				}
				proc.requeueCount++
				proc.currentRequeueCount++

				// only requeue to queues, should never happen for topics
				if isQueueDestination(destination) {
					queue := proc.qm.Find(destination)
					queue.Requeue(r.Frame)
				}

			case client.ConnectedOp:
				//register connection
				proc.connectCount++
				proc.currentConnectCount++
				proc.connections[r.Conn.Id()] = r.Conn

			case client.DisconnectedOp:
				proc.disconnectCount++
				proc.currentDisconnectCount++
				delete(proc.connections, r.Conn.Id())
			}
		}
	}
	// this is no longer required for go 1.1
	panic("not reached")
}

func isQueueDestination(dest string) bool {
	return strings.HasPrefix(dest, QueuePrefix)
}

func (proc *requestProcessor) Listen(l net.Listener) {
	var conn_id int64 = 0

	timeout := time.Duration(0) // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				if timeout == 0 {
					timeout = 5 * time.Millisecond
				} else {
					timeout *= 2
				}
				if max := 5 * time.Second; timeout > max {
					timeout = max
				}
				log.Errorf("stomp: Accept error: %v; retrying in %v", err, timeout)
				time.Sleep(timeout)
				continue
			}
			return
		}
		timeout = 0
		// TODO: need to pass Server to connection so it has access to
		// configuration parameters.
		client.NewConn(proc.config, rw, proc.ch, conn_id)
		//conn := client.NewConn(config, rw, proc.ch, conn_id)
		//notify about new connect
		//proc.ch <- Request{Op: ConnectedOp, Conn: c}
		//proc.connections[conn_id] = conn
		conn_id++
	}
	// This is no longer required for go 1.1
	log.Panic("not reached")
}

type config struct {
	server *Server
}

func newConfig(s *Server) *config {
	return &config{server: s}
}

func (c *config) HeartBeat() time.Duration {
	duration := c.server.HeartBeatDuration()
	if duration == time.Duration(0) {
		return DefaultHeartBeat
	}
	return duration
}

func (c *config) IsDebug() bool {
	return c.server.Config.IsDebug
}

func (c *config) MaxPendingReads() int {
	if c.server.Config.MaxPendingReads <= 0 {
		return 16
	}
	return c.server.Config.MaxPendingReads
}

func (c *config) MaxPendingWrites() int {
	if c.server.Config.MaxPendingWrites <= 0 {
		return 16
	}
	return c.server.Config.MaxPendingWrites
}

func (c *config) Authenticate(login, passcode string) bool {
	if c.server.Authenticator != nil {
		return c.server.Authenticator.Authenticate(login, passcode)
	}

	// no authentication defined
	return true
}
