package stomp

import (
	"fmt"
	//"log"

	"github.com/go-stomp/stomp/frame"
)

// The Subscription type represents a client subscription to
// a destination. The subscription is created by calling Conn.Subscribe.
//
// Once a client has subscribed, it can receive messages from the C channel.
type Subscription struct {
	C              chan *Message
	id             string
	destination    string
	conn           *Conn
	ackMode        AckMode
	completed      bool
	opts           []func(*frame.Frame) error
	controlChannel chan func()
	frameCh        chan *frame.Frame
}

// BUG(jpj): If the client does not read messages from the Subscription.C
// channel quickly enough, the client will stop reading messages from the
// server.

// Identification for this subscription. Unique among
// all subscriptions for the same Client.
func (s *Subscription) Id() string {
	return s.id
}

// Destination for which the subscription applies.
func (s *Subscription) Destination() string {
	return s.destination
}

// AckMode returns the Acknowledgement mode specified when the
// subscription was created.
func (s *Subscription) AckMode() AckMode {
	return s.ackMode
}

// Active returns whether the subscription is still active.
// Returns false if the subscription has been unsubscribed.
func (s *Subscription) Active() bool {
	return !s.completed
}

// Unsubscribes and closes the channel C.
func (s *Subscription) Unsubscribe() {

	s.controlChannel <- func() {
		f := frame.New(frame.UNSUBSCRIBE, frame.Id, s.id)
		s.conn.sendFrame(f)
		s.completed = true
		close(s.C)
	}
}

// Read a message from the subscription. This is a convenience
// method: many callers will prefer to read from the channel C
// directly.
func (s *Subscription) Read() (*Message, error) {
	msg, ok := <-s.C
	if !ok {
		//continue
		return nil, ErrCompletedSubscription
	}
	if msg.Err != nil {
		//continue
		return nil, msg.Err
	}
	//continue
	return msg, nil
}

func (s *Subscription) readLoop(ch chan *frame.Frame) {
	for {

		select {
		case f := <-s.controlChannel:
			f()
			continue

		case f, ok := <-ch:
			if !ok {
				continue
			}

			if f.Command == frame.MESSAGE {
				destination := f.Header.Get(frame.Destination)
				contentType := f.Header.Get(frame.ContentType)
				msg := &Message{
					Destination:  destination,
					ContentType:  contentType,
					Conn:         s.conn,
					Subscription: s,
					Header:       f.Header,
					Body:         f.Body,
				}
				if !s.completed {
					s.C <- msg
				}
			} else if f.Command == frame.ERROR {
				//log.Warn("subs: f.Command == frame.ERROR")
				message, _ := f.Header.Contains(frame.Message)
				text := fmt.Sprintf("Subscription %s: %s: ERROR message:%s",
					s.id,
					s.destination,
					message)
				log.Debugf("subs: test=%s", text)
				contentType := f.Header.Get(frame.ContentType)
				msg := &Message{
					Err: &Error{
						Message: f.Header.Get(frame.Message),
						Frame:   f,
					},
					ContentType:  contentType,
					Conn:         s.conn,
					Subscription: s,
					Header:       f.Header,
					Body:         f.Body,
				}
				if !s.completed {
					s.C <- msg
				}
				//s.completed = true
				continue
			}
		}
	}
}
