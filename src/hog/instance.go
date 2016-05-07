package hog

import (
	"config"
	"logger"
	"net"
	"time"
)

type Instance struct {
	ipAddress  net.IPAddr
	connection net.Conn
	name       string
	// This channel should be pushed to every 1-5 seconds by the client.
	lastReceived chan time.Time
}

/*
This is purposely higher than the client push rate
to account for possible network latency.
*/
const heartbeatInterval = 10 // Seconds

/*
Keep checking to see if the server-client connection
is still alive. If the connection is no longer alive then
*/
func (i *Instance) heartbeat() {
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
still occupied. Rather, instead of discarding the new heartbeat time
replace the stored time with the new time.
*/
func (i *Instance) listen() {
	for true {
		buff := make([]byte, 256)
		_, err := i.connection.Read(buff)

		if err != nil && config.Debug {
			logger.Error.Println(err)
		}
	}
}

func (i *Instance) Close() {
	ExitMessage(i.name)
	i.connection.Close()
	close(i.lastReceived)
}
