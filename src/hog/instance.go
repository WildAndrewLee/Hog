package hog

import (
	"net"
	"time"
)

type Instance struct {
	ipAddress    net.IPAddr
	connection   net.Conn
	name         string
	lastReceived time.Time
}

const heartbeatInterval = 5 // Seconds

func (i *Instance) heartbeat() {
	i.lastReceived = time.Now()

	hbi := heartbeatInterval * time.Second
	ticker := time.NewTicker(hbi)

	go func() {
		for t := range ticker.C {
			if i.lastReceived.Add(hbi).Before(t) {
				ticker.Stop()
				i.Close()
				return
			} else {
				i.lastReceived = t
			}
		}
	}()
}

func (i *Instance) Close() {

}
