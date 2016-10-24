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
	"os"

	conf "github.com/KristinaEtc/config"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
)

var configFile string

// These fields are populated by govvv
var (
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
	Version    string
)

var log = slf.WithContext("go-stompd-server.go")

// ConfFile is a file with all program options
type ConfFile struct {
	Global server.ServerConfig
}

var globalOpt = ConfFile{
	Global: server.ServerConfig{
		ListenAddr:       "localhost:61614",
		Id:               "go-stomp-server",
		Name:             "",
		Heartbeat:        30,
		Status:           30,
		StatusLog:        300,
		MaxPendingReads:  300, //1 message per second * 5minutes
		MaxPendingWrites: 300,
		IsDebug:          false,
	},
}

func main() {

	log.Infof("BuildDate=%s\n", BuildDate)
	log.Infof("GitCommit=%s\n", GitCommit)
	log.Infof("GitBranch=%s\n", GitBranch)
	log.Infof("GitState=%s\n", GitState)
	log.Infof("GitSummary=%s\n", GitSummary)
	log.Infof("VERSION=%s\n", Version)
	globalOpt.Global.Version = Version

	conf.ReadGlobalConfig(&globalOpt, "GlobalConf")

	a := auth.NewAuth()
	s := server.NewServer(&globalOpt.Global, a)

	log.Error("-----------------------------------------------")
	err := s.ListenAndServe()
	if err != nil {
		log.Errorf("error ListenAndServe %v", err)
		os.Exit(1)
	}
	//server.Serve(l, a)
}
