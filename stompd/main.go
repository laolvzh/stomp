/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

import (
	"flag"
	"github.com/KristinaEtc/slflog"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
	"net"
)

var (
	listenAddr     = flag.String("addr", ":61613", "Listen address")
	helpFlag       = flag.Bool("help", false, "Show this help text")
	configAuthFile = flag.String("auth", "auth.json", "configfile with logins and passwords")
	logPath        = flag.String("logpath", "logs", "path to logfiles")
	logLevel       = flag.String("loglevel", "INFO", "IFOO, DEBUG, ERROR, WARN, PANIC, FATAL - loglevel for stderr")
)

func main() {

	flag.Parse()
	slflog.InitLoggers(*logPath, *logLevel)
	// TODO: add Close method!!
	//defer slflog.Close()
	log := slf.WithContext("go-stompd-server.go")

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(*configAuthFile)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	server.Serve(l, a)
}
