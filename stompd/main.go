/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

import (
	"Nominatim/lib/utils/basic"
	"flag"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/ventu-io/slf"
	"github.com/ventu-io/slog"
	"net"
	"os"
	"path/filepath"
)

// TODO: experimenting with ways to gracefully shutdown the server,
// at the moment it just dies ungracefully on SIGINT.

/*

func main() {
	// create a channel for listening for termination signals
	stopChannel := newStopChannel()

	for {
		select {
		case sig := <-stopChannel:
			log.Println("received signal:", sig)
			break
		}
	}

}
*/

var (
	listenAddr     = flag.String("addr", ":61613", "Listen address")
	helpFlag       = flag.Bool("help", false, "Show this help text")
	configAuthFile = flag.String("auth", "auth.json", "configfile with logins and passwords")
	debugMode      = flag.Bool("debug", true, "Debug mode")
)

const LogDir = "logs/"

const (
	errorFilename = "error.log"
	infoFilename  = "info.log"
	debugFilename = "debug.log"
)

var (
	bhDebug, bhInfo, bhError, bhDebugConsole *basic.Handler
	logfileInfo, logfileDebug, logfileError  *os.File
	lf                                       slog.LogFactory

	log slf.StructuredLogger
)

// Init loggers
func init() {

	bhDebug = basic.New(slf.LevelDebug)
	bhDebugConsole = basic.New(slf.LevelDebug)
	bhInfo = basic.New()
	bhError = basic.New(slf.LevelError)

	// optionally define the format (this here is the default one)
	bhInfo.SetTemplate("{{.Time}} [\033[{{.Color}}m{{.Level}}\033[0m] {{.Context}}{{if .Caller}} ({{.Caller}}){{end}}: {{.Message}}{{if .Error}} (\033[31merror: {{.Error}}\033[0m){{end}} {{.Fields}}")

	// TODO: create directory in /var/log, if in linux:
	// if runtime.GOOS == "linux" {
	os.Mkdir("."+string(filepath.Separator)+LogDir, 0766)

	// interestings with err: if not initialize err before,
	// how can i use global logfileInfo?
	var err error
	logfileInfo, err = os.OpenFile(LogDir+infoFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Could not open/create %s logfile", LogDir+infoFilename)
	}

	logfileDebug, err = os.OpenFile(LogDir+debugFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Could not open/create logfile", LogDir+debugFilename)
	}

	logfileError, err = os.OpenFile(LogDir+errorFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Panicf("Could not open/create logfile", LogDir+errorFilename)
	}

	if *debugMode == true {
		bhDebugConsole.SetWriter(os.Stdout)
	}

	bhDebug.SetWriter(logfileDebug)
	bhInfo.SetWriter(logfileInfo)
	bhError.SetWriter(logfileError)

	lf = slog.New()
	lf.SetLevel(slf.LevelDebug) //lf.SetLevel(slf.LevelDebug, "app.package1", "app.package2")
	lf.SetEntryHandlers(bhInfo, bhError, bhDebug)

	if *debugMode == true {
		lf.SetEntryHandlers(bhInfo, bhError, bhDebug, bhDebugConsole)
	} else {
		lf.SetEntryHandlers(bhInfo, bhError, bhDebug)
	}

	// make this into the one used by all the libraries
	slf.Set(lf)

	log = slf.WithContext("main-stompd.go")
}

func main() {
	flag.Parse()
	if *helpFlag {
		log.Warnf("Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parsed()

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(*configAuthFile)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	server.Serve(l, a)
}
