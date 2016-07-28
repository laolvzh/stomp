/*
A simple, stand-alone STOMP server.

TODO: graceful shutdown

TODO: UNIX daemon functionality

TODO: Windows service functionality (if possible?)

TODO: Logging options (syslog, windows event log)
*/
package main

/*go build -ldflags "-X github.com/KristinaEtc/slflog.configLogfile=/usr/share/go-stomp-server/go-stomp-server.logconfig
-X go-stomp-server.pathToConfig=/usr/share/go-stomp-server/go-stomp-server.config" go-stomp-server.go */

//important: do not move
import (
	"flag"
	"net"
	"os"

	_ "github.com/KristinaEtc/slflog"
	"github.com/go-stomp/stomp/server"
	"github.com/go-stomp/stomp/server/auth"
	"github.com/kardianos/osext"
	"github.com/ventu-io/slf"
)

var configFile string

var log = slf.WithContext("go-stompd-server.go")

var (
	listenAddr     = flag.String("addr", ":61613", "Listen address")
	helpFlag       = flag.Bool("help", false, "Show this help text")
	configAuthFile = flag.String("conf", getPathToConfig(), "configfile with logins and passwords")
)

func main() {

	flag.Parse()
	// TODO: add Close method!!
	//defer slflog.Close()

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.WithCaller(slf.CallerShort).Fatalf("Failed to listen: %s", err.Error())
	}
	defer func() { l.Close() }()

	a := auth.NewAuth(*configAuthFile)

	log.Debugf("listening on %v %s", l.Addr().Network(), l.Addr().String())
	log.Error("[go-stomp-server]--------------------------new connection---------------------")
	server.Serve(l, a)
}

func getPathToConfig() string {

	var path = configFile

	// path to config was setted by a linker value
	if path != "" {
		exist, err := exists(path)
		if err != nil {
			log.WithCaller(slf.CallerShort).Errorf("Error: wrong configure file from linker value %s: %s\n", path, err.Error())
			path = ""
		} else if exist != true {
			log.WithCaller(slf.CallerShort).Errorf("Error: Configure file from linker value %s: does not exist\n", path)
			path = ""
		}
	}

	// no path from a linker value or wrong linker value; searching where a binary is situated
	if path == "" {
		pathTemp, err := osext.Executable()
		if err != nil {
			log.WithCaller(slf.CallerShort).Errorf("Error: could not get a path to binary file for getting configfile: %s\n", err.Error())
		} else {
			path = pathTemp + ".config"
		}
	}
	log.WithCaller(slf.CallerShort).Infof("Configfile that will be used: [%s]", path)
	return path
}

// Exists returns whether the given file or directory exists or not.
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
