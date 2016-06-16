/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

import (
	"Nominatim/lib/utils/basiclog"
	"flag"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/kardianos/osext"
	"github.com/ventu-io/slf"
	"github.com/ventu-io/slog"
	"net"
	"os"
	"path/filepath"
)

var (
	listenAddr     = flag.String("addr", ":61613", "Listen address")
	helpFlag       = flag.Bool("help", false, "Show this help text")
	configAuthFile = flag.String("auth", "auth.json", "configfile with logins and passwords")
	debugMode      = flag.Bool("debug", true, "Debug mode")
	logPath        = flag.String("logpath", "logs", "path to logfiles")
	logLevel       = flag.String("loglevel", "INFO", "IFOO, DEBUG, ERROR, WARN, PANIC, FATAL")
)

const (
	errorFilename = "error.log"
	infoFilename  = "info.log"
	debugFilename = "debug.log"
)

var (
	bhDebug, bhInfo, bhError, bhDebugConsole, bhStdError *basiclog.Handler
	logfileInfo, logfileDebug, logfileError              *os.File
	lf                                                   slog.LogFactory

	log slf.StructuredLogger
)

// Init loggers
func initLoggers() {

	var logHandlers []slog.EntryHandler

	// optionally define the format (this here is the default one)
	//bhInfo.SetTemplate("{{.Time}} [\033[{{.Color}}m{{.Level}}\033[0m] {{.Context}}{{if .Caller}} ({{.Caller}}){{end}}: {{.Message}}{{if .Error}} (\033[31merror: {{.Error}}\033[0m){{end}} {{.Fields}}")

	basiclog.ConfigWriterOutput(&logHandlers, getLogLevel(*logLevel), os.Stderr)

	err := setLogOutput(&logHandlers)
	if err != nil {
		basiclog.SafeLog("[go-stomp-server] Error init loggers: " + err.Error() + "\n")
	}

	lf = slog.New()
	lf.SetLevel(slf.LevelDebug)
	lf.SetEntryHandlers(logHandlers...)
	slf.Set(lf)

	log = slf.WithContext("go-stompd-server.go")
}

func main() {

	flag.Parse()
	initLoggers()

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(*configAuthFile)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	server.Serve(l, a)
}

func getPathForLogDir() (string, error) {

	if filepath.IsAbs(*logPath) == true {
		return *logPath, nil
	} else {
		filename, err := osext.Executable()
		if err != nil {
			return "", err
		}

		fpath := filepath.Dir(filename)
		fpath = filepath.Join(fpath, *logPath)
		return fpath, nil
	}

}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func setLogOutput(logHandlers *[]slog.EntryHandler) error {

	pathForLogs, err := getPathForLogDir()
	if err != nil {
		return err
	}
	exist, err := exists(pathForLogs)
	if err != nil {
		return err
	}
	if !exist {
		err = os.Mkdir(pathForLogs, 0755)
		if err != nil {
			return err
		}
	}

	basiclog.ConfigFileOutput(logHandlers, slf.LevelDebug, filepath.Join(pathForLogs, debugFilename))
	basiclog.ConfigFileOutput(logHandlers, slf.LevelInfo, filepath.Join(pathForLogs, infoFilename))
	basiclog.ConfigFileOutput(logHandlers, slf.LevelError, filepath.Join(pathForLogs, errorFilename))

	return nil
}

func getLogLevel(lvl string) slf.Level {

	switch lvl {
	case slf.LevelDebug.String():
		return slf.LevelDebug

	case slf.LevelInfo.String():
		return slf.LevelInfo

	case slf.LevelWarn.String():
		return slf.LevelWarn

	case slf.LevelError.String():
		return slf.LevelError

	case slf.LevelFatal.String():
		return slf.LevelFatal
	case slf.LevelPanic.String():
		return slf.LevelPanic
	default:
		return slf.LevelDebug
	}
}
