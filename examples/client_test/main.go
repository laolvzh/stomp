package main

import (
	"flag"
	"fmt"
	//"github.com/gmallard/stompngo"
	//"github.com/gmallard/stompngo_examples/sngecomm"
	"github.com/go-stomp/stomp"
	"os"
	"strconv"
	"time"
)

const defaultPort = ":61614"

var serverAddr = flag.String("server", "localhost:61614", "STOMP server endpoint")
var messageCount = flag.Int("count", 10, "Number of messages to send/receive")
var destination = flag.String("topic", "TOPIC", "Destination topic")
var queueName = flag.String("queue", "/queue/QueueAnswer", "Destination queue")
var helpFlag = flag.Bool("help", false, "Print help text")
var stop = make(chan bool)

// these are the default options that work with RabbitMQ
var options []func(*stomp.Conn) error = []func(*stomp.Conn) error{
	stomp.ConnOpt.Login("guest", "guest"),
	stomp.ConnOpt.Host("/"),
}

func main() {
	flag.Parse()
	if *helpFlag {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

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
		println("cannot connect to server", err.Error())
		return
	}

	//for i := 1; i <= *messageCount; i++ {

	/*s := stompngo.Headers{"destination", sngecomm.Dest(),
		"persistent", "true"} // send headers
	/*m := exampid + " message: "
	for i := 1; i <= sngecomm.Nmsgs(); i++ {
		t := m + fmt.Sprintf("%d", i)
		fmt.Println(sngecomm.ExampIdNow(exampid), "sending now:", t)
		//e := conn.Send(s, t)
	*/

	i := 2
	for {

		time.Sleep(1000 * time.Millisecond)

		//text := fmt.Sprintf()
		err = conn.Send(*destination, "text/plain",
			[]byte(*queueName+" "+strconv.Itoa(i)), nil...)
		if err != nil {
			println("failed to send to server", err)
			return
		}
		i++
	}
	println("sender finished")
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

	sub, err := conn.Subscribe(*queueName, stomp.AckAuto)
	if err != nil {
		println("cannot subscribe to", *queueName, err.Error())
		return
	}
	close(subscribed)

	//for i := 1; i <= *messageCount; i++ {
	for {
		fmt.Println("here?\n")
		msg := <-sub.C
		fmt.Println("GGG?\n")
		//expectedText := fmt.Sprintf("Message #%d", i)
		actualText := string(msg.Body)
		//if expectedText != actualText {
		//	println("Expected:", expectedText)
		println("Actual:", actualText)
		//}
	}
	println("receiver finished")

}
