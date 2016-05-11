package hog

import (
	"network/opcodes"
	"time"
)

func sendMessage(r rawMessage, m message) {
	if r.i.name == "" {
		r.i.connection.Write(NewMessage(opcodes.OpRefused))
	} else if len(m.Args) != 1 {
		invalidOp(r.i, r.b)
	} else {
		broadcastMessage(NewMessage(opcodes.ReceiveMessage, r.i.name, m.Args[0]))
	}
}

func heartbeat(r rawMessage, m message) {
	select {
	case r.i.lastReceived <- time.Now():
	default:
		<-r.i.lastReceived
		r.i.lastReceived <- time.Now()
	}
}

func tryChangeName(r rawMessage, m message) {
	if len(m.Args) != 1 {
		invalidOp(r.i, r.b)
	}
	r.i.ChangeName(m.Args[0])
}

func connect(r rawMessage, m message) {
	tryChangeName(r, m)
	joinMessage(r.i.name)
}
