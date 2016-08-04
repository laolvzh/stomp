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

	"github.com/KristinaEtc/utils"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
)

var configFile string

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

var globalOpt = ConfFile{Global: GlobalConf{ListenAddr: ":61614"}}

func main() {

	utils.GetFromGlobalConf(&globalOpt, "GlobalConf")

	// TODO: add Close method!!
	//defer slflog.Close()

	l, err := net.Listen("tcp", globalOpt.Global.ListenAddr)
	if err != nil {
		log.WithCaller(slf.CallerShort).Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(utils.GetConfigFilename())

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	log.Error("-----------------------------------------------")
	server.Serve(l, a)
}
