package hog

import (
	"bytes"
	"config"
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
	ipAddress  net.IPAddr
	connection net.Conn
	name       string
	// This channel should be pushed to every 1-5 seconds by the client.
	lastReceived chan time.Time
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
		default:
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
		logger.Info.Println("Failed to receive heartbeat from " + i.name + ".")
	}

	t.Stop()
	i.Close()
}

/*
Important Note:
When writing the receive heartbeat code, make sure to put it
in a select so there is no blocking when the heartbeat channel is
still occupied.
*/
func (i *instance) listen() {
	for {
		buff := new(bytes.Buffer)
		_, err := io.Copy(buff, i.connection)

		if err != nil {
			/*
			   If err is EOF that means the connection was
			   closed by the client.
			*/
			if err != io.EOF && config.Debug {
				logger.Error.Println(err)
			}

			i.Close()
			return
		}

		go i.connection.Write(NewMessage(opcodes.Heartbeat))
	}
}

func (i *instance) Close() {
	exitMessage(i.name)
	i.connection.Close()
	close(i.lastReceived)
}

func (i *instance) ChangeName(name string) {
	if i.name == name {
		return
	}

	if nameInUse(name) {
		i.connection.Write(NewMessage(opcodes.NameInUse))
	}

	i.name = name
	joinMessage(i.name)
}
