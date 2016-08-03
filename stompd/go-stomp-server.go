/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net"

	_ "github.com/KristinaEtc/slflog"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
)

var configFile string

var log = slf.WithContext("go-stompd-server.go")

// GlobalConf is a struct with global options,
// like server address and config auth filename
type GlobalConf struct {
	ListenAddr     string
	ConfigAuthFile string
}

// ConfFile is a file with all program options
type ConfFile struct {
	Global GlobalConf
}

var defaulfGlobalOpt = GlobalConf{
	ListenAddr:     "61614",
	ConfigAuthFile: "",
}

func main() {

	flag.Parse()

	var cf ConfFile
	getGlobalConf(&cf)

	// TODO: add Close method!!
	//defer slflog.Close()

	l, err := net.Listen("tcp", cf.Global.ListenAddr)
	if err != nil {
		log.WithCaller(slf.CallerShort).Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(cf.Global.ConfigAuthFile)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	log.Error("-----------------------------------------------")
	server.Serve(l, a)
}

func getGlobalConf(cf *ConfFile) {
	file, e := ioutil.ReadFile("global.conf")
	if e != nil {
		log.Errorf("File error: %s\n", e.Error())
		cf.Global = defaulfGlobalOpt
	}

	if err := json.Unmarshal([]byte(file), cf); err != nil {
		log.Error(err.Error())
		cf.Global = defaulfGlobalOpt
	}
	//log.Errorf("Results: %v\n", cf)
}
