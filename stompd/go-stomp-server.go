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
	"github.com/kardianos/osext"
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

var defaulfGlobalOpt = GlobalConf{
	ListenAddr: ":61614",
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

	a := auth.NewAuth(getConfigFilename())

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	log.Error("-----------------------------------------------")
	server.Serve(l, a)
}

func getGlobalConf(cf *ConfFile) {

	file, e := ioutil.ReadFile(getConfigFilename())
	if e != nil {
		log.WithCaller(slf.CallerShort).Errorf("Error: %s\n", e.Error())
		cf.Global = defaulfGlobalOpt
	}

	if err := json.Unmarshal([]byte(file), cf); err != nil {
		log.WithCaller(slf.CallerShort).Errorf("Error parsing JSON: %s", err.Error())
		cf.Global = defaulfGlobalOpt
	} else {
		log.Infof("Global options will be used from [%s] file", getConfigFilename())
	}
	//log.Errorf("Results: %v\n", cf)
}

func getConfigFilename() string {
	binaryPath, err := osext.Executable()
	if err != nil {
		log.WithCaller(slf.CallerShort).Errorf("Error: could not get a path to binary file: %s\n", err.Error())
	}
	return binaryPath + ".config"
}
