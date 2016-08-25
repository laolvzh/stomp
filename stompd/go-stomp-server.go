/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

import _ "github.com/KristinaEtc/slflog"

import (
	"net"

	conf "github.com/KristinaEtc/config"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
)

var configFile string

var (
	// These fields are populated by govvv
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	Version    string
)

var log = slf.WithContext("go-stompd-server.go")

// GlobalConf is a struct with global options,
// like server address and config auth filename
type GlobalConf struct {
	ListenAddr string
}

// ConfFile is a file with all program options
type ConfFile struct {
	Global GlobalConf
}

var globalOpt = ConfFile{Global: GlobalConf{ListenAddr: "localhost:61614"}}

func main() {

	conf.ReadGlobalConfig(&globalOpt, "GlobalConf")

	// TODO: add Close method!!
	//defer slflog.Close()

	l, err := net.Listen("tcp", globalOpt.Global.ListenAddr)
	if err != nil {
		log.WithCaller(slf.CallerShort).Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth()

	log.Infof("BuildDate=%s\n", BuildDate)
	log.Infof("GitCommit=%s\n", GitCommit)
	log.Infof("GitBranch=%s\n", GitBranch)
	log.Infof("GitState=%s\n", GitState)
	log.Infof("GitSummary=%s\n", GitSummary)
	log.Infof("VERSION=%s\n", Version)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	log.Error("-----------------------------------------------")
	server.Serve(l, a)
}
