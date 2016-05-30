package main

import (
	"flag"

	l4g "github.com/alecthomas/log4go"
	"github.com/go-stomp/stomp"
	"github.com/go-stomp/stomp/server/utils"
	"os"
	"time"
)

var log l4g.Logger = l4g.NewLogger()

const (
	defaultPort = ":61614"
	clientID    = "clientID"
)

var (
	testFile = "test.csv"
	LOGFILE  = "client.log"
)

var (
	serverAddr  = flag.String("server", "localhost:61614", "STOMP server endpoint")
	destination = flag.String("topic", "mainTopic", "Destination topic")
	queueFormat = flag.String("queue", "/queue/", "Queue format")
	stop        = make(chan bool)
)

// these are the default options that work with RabbitMQ
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func init() {
	log.AddFilter("stdout", l4g.INFO, l4g.NewConsoleLogWriter())
	log.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter(LOGFILE, false))
	//
}

func main() {

	// logger configuration
	defer log.Close()

	flag.Parsed()
	flag.Parse()

	subscribed := make(chan bool)

	go recvMessages(subscribed)
	// wait until we know the receiver has subscribed
	<-subscribed

	go sendMessages()

	<-stop
	<-stop
}

func sendMessages() {
	defer func() {
		stop <- true
	}()

	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil {
		log.Error("cannot connect to server", err.Error())
		return
	}

	fs, err := utils.NewFileScanner(testFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer fs.Close()

	//пустой слайс или массив из одного нила
	fs.Scanner = fs.GetScanner()

	for fs.Scanner.Scan() {
		locs := fs.Scanner.Text()
		//log.Info("locs: %s", locs)
		time.Sleep(1000 * time.Millisecond)
		reqInJSON, err := utils.MakeReq(locs, clientID, log)
		if err != nil {
			log.Error("Could not get coordinates in JSON: wrong format")
			continue
		}
		//log.Info("reqInJSON: %s", *reqInJSON)

		time.Sleep(1000 * time.Millisecond)

		err = conn.Send(*destination, "text/json", []byte(*queueFormat+clientID+" "+*reqInJSON), nil...)
		if err != nil {
			println("failed to send to server", err)
			return
		}
	}
}

func recvMessages(subscribed chan bool) {
	defer func() {
		stop <- true
	}()
	conn, err := stomp.Dial("tcp", *serverAddr, options...)
	if err != nil {
		println("cannot connect to server", err.Error())
		return
	}

	sub, err := conn.Subscribe(*queueFormat+clientID, stomp.AckAuto)
	if err != nil {
		println("cannot subscribe to", *queueFormat+clientID, err.Error())
		return
	}
	close(subscribed)

	for {
		msg := <-sub.C
		if msg.Body == nil {
			log.Warn("Got empty message; ignore")
			continue
		}
		actualText := string(msg.Body)
		println("Actual:", actualText)
	}
}
