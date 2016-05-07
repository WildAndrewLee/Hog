package hog

import (
	"fmt"
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
	ticker := time.NewTicker(hbi)

	go func() {
		for now := range ticker.C {
			stop := false

			/*
				If the lastReceived channel is empty that means
				we have not recieved a heartbeat within the last
				heartbeatInterval seconds. If there is a heartbeat,
				check to see if it is within the interval.
			*/
			select {
			case lastSeen := <-i.lastReceived:
				stop = lastSeen.Add(hbi).Before(now)
			default:
				stop = true
			}

			if stop {
				logger.LogString(fmt.Sprintf("Failed to receive heartbeat from %s.", i.name))
				ticker.Stop()
				i.Close()
			}
		}
	}()
}

/*
Important Note:
When writing the receive heartbeat code, make sure to put it
in a select so there is no blocking when the heartbeat channel is
still occupied. Rather, instead of discarding the new heartbeat time
replace the stored time with the new time.
*/
func (i *Instance) listen() {

}

func (i *Instance) Close() {
	ExitMessage(i.name)
	i.connection.Close()
	close(i.lastReceived)
}
