package hog

import (
	"config"
	"fmt"
	"io"
	"logger"
	"net"
	"network/opcodes"
	"time"
)

/*
This is purposely higher than the recommended client push rate
to account for possible network latency.
*/
const heartbeatInterval = 10 // Seconds

type instance struct {
	ipAddress    net.Addr
	connection   net.Conn
	name         string
	lastReceived chan time.Time // This channel should be pushed to every 1-5 seconds by the client.
	e            chan bool
	m            chan []byte
}

/*
Keep checking to see if the server-client connection
is still alive. If the connection is no longer alive then
close the connection on our end.
*/
func (i *instance) heartbeat() {
	hbi := heartbeatInterval * time.Second
	t := time.NewTicker(hbi)
	s := false

	for now := range t.C {
		/*
			If the lastReceived channel is empty that means
			we have not recieved a heartbeat within the last
			heartbeatInterval seconds. If there is a heartbeat,
			check to see if it is within the interval.
		*/
		select {
		case last := <-i.lastReceived:
			s = last.Add(hbi).Before(now)
		case <-i.e:
			return
		default:
			/*
			   This should not happen because
			   we always give instances an initial
			   heartbeat on creation.
			*/
			s = true
		}

		if s {
			break
		}

		/*
		   Send heartbeat to client so it knows it's still
		   connected.
		*/
		i.connection.Write(NewMessage(opcodes.Heartbeat))
	}

	if config.Debug {
		logger.Info.Println(fmt.Sprintf("Failed to receive heartbeat from \"%s\"", i.name), i.ipAddress)
	}

	i.e <- true
}

func (i *instance) listenHelper() {
	for {
		tr := make([]byte, 2)

		i.connection.Read(tr)
		l := int(tr[0]<<8 | tr[1])

		buff := make([]byte, l)
		_, err := i.connection.Read(buff)

		if err != nil {
			/*
			   If err is EOF that means the connection was
			   closed by the client.
			*/
			if err != io.EOF && config.Debug {
				logger.Error.Println(err)
			}

			i.e <- true
			return
		}

		i.m <- buff
	}
}

/*
Important Note:
When writing the receive heartbeat code, make sure to put it
in a select so there is no blocking when the heartbeat channel is
still occupied.
*/
func (i *instance) listen() {
	go i.listenHelper()

	for {
		select {
		case m := <-i.m:
			enqueueMessage(i, m)
		case e := <-i.e:
			if e {
				i.Close()
				return
			}
		}
	}
}

func (i *instance) Close() {
	logger.Info.Println(fmt.Sprintf("Closing connection for \"%s\"", i.name), i.ipAddress)

	exitMessage(i.name)
	i.connection.Close()

	if i.lastReceived != nil {
		close(i.lastReceived)
	}
}

func (i *instance) ChangeName(name string) {
	if i.name == name {
		return
	}

	if len(name) == 0 || len(name) > config.MaxNameLength {
		i.connection.Write(NewMessage(opcodes.NameTooLong))
	}

	if nameInUse(name) {
		i.connection.Write(NewMessage(opcodes.NameInUse))
		return
	}

	i.name = name

	joinMessage(i.name)
}

func NewInstance(c net.Conn) *instance {
	i := instance{}
	i.name = ""
	i.connection = c
	i.ipAddress = c.RemoteAddr()
	i.lastReceived = make(chan time.Time, 1)
	i.e = make(chan bool, 1)
	i.m = make(chan []byte, 1)

	i.e <- false
	i.lastReceived <- time.Now()

	return &i
}
