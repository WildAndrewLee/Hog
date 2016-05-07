package hog

import (
	"net"
	"time"
)

type Instance struct {
	ipAddress    net.IPAddr
	connection   net.Conn
	name         string
	lastReceived chan time.Time
}

const heartbeatInterval = 5 // Seconds

func (i *Instance) heartbeat() {
	hbi := heartbeatInterval * time.Second
	ticker := time.NewTicker(hbi)

	go func() {
		for range ticker.C {
			select {
			case <-i.lastReceived:
			default:
				ticker.Stop()
				i.Close()
				return
			}
		}
	}()
}

func (i *Instance) Close() {

}
