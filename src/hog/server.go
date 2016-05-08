package hog

import (
	"config"
	"logger"
	"network/opcodes"
	"time"
)

type rawMessage struct {
	i instance
	b []byte
}

// Channel for messages received by the server.
var messageQueue chan rawMessage
var clients []instance

func init() {
	messageQueue = make(chan rawMessage, config.MessageQueueSize)
}

func enqueueMessage(i instance, message []byte) {
	messageQueue <- rawMessage{i: i, b: message}
}

func joinMessage(name string) {
	m := NewMessage(opcodes.Join, name)

	for _, c := range clients {
		c.connection.Write(m)
	}
}

func exitMessage(name string) {
	m := NewMessage(opcodes.Leave, name)

	for _, c := range clients {
		c.connection.Write(m)
	}
}

func nameInUse(name string) bool {
	for _, c := range clients {
		if c.name == name {
			return true
		}
	}

	return false
}

func processQueue() {
	for m := range messageQueue {
		if config.Debug {
			logger.Info.Println("Processing message: ", m)
		}

		processMessage(m)
	}
}

func processMessage(r rawMessage) {
	i := r.i
	b := r.b
	m := ParseMessage(b)

	switch m.Op {
	case opcodes.SendMessage:
		if i.name == "" {
			i.connection.Write(NewMessage(opcodes.OpRefused))
		} else {
			enqueueMessage(i, NewMessage(opcodes.ReceiveMessage, i.name, m.Args[0]))
		}
	case opcodes.Heartbeat:
		select {
		case i.lastReceived <- time.Now():
		default:
		}
	case opcodes.Connect:
		i.ChangeName(m.Args[0])
	case opcodes.ChangeName:
		i.ChangeName(m.Args[0])
	default:
	}
}
