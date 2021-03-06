package hog

import (
	"config"
	"fmt"
	"logger"
	"net"
	"network/opcodes"
	"strconv"
)

type rawMessage struct {
	i *instance
	b []byte
}

// Channel for messages received by the server.
var messageQueue chan rawMessage
var clients []*instance
var handlers map[byte]func(r rawMessage, m message)

func init() {
	messageQueue = make(chan rawMessage, config.MessageQueueSize)
	handlers = make(map[byte]func(r rawMessage, m message))

	handlers[opcodes.SendMessage] = sendMessage
	handlers[opcodes.ChangeName] = tryChangeName
	handlers[opcodes.Connect] = connect
	handlers[opcodes.Heartbeat] = heartbeat
}

func enqueueMessage(i *instance, message []byte) {
	messageQueue <- rawMessage{i: i, b: message}
}

func broadcastMessage(m []byte) {
	for _, c := range clients {
		c.connection.Write(m)
	}
}

func joinMessage(name string) {
	broadcastMessage(NewMessage(opcodes.Join, name))
}

func exitMessage(name string) {
	broadcastMessage(NewMessage(opcodes.Leave, name))
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
			logger.Info.Println(fmt.Sprintf("Processing message: \"%s\"", m.i.name), m.b)
		}

		processMessage(m)
	}
}

func invalidOp(i *instance, b []byte) {
	if config.Debug {
		logger.Info.Println("Received invalid message:", b)
	}
	i.connection.Write(NewMessage(opcodes.OpRefused))
}

func processMessage(r rawMessage) {
	m := ParseMessage(r.b)

	if fn, ok := handlers[m.Op]; ok {
		fn(r, m)
	} else {
		invalidOp(r.i, r.b)
	}
}

func Start() {
	l, err := net.Listen("tcp", ":"+strconv.Itoa(config.Port))

	if err != nil {
		logger.Error.Println("An error occurred while attempting to setup the server.")
		logger.Error.Println(err)
	}

	defer l.Close()

	go processQueue()

	for {
		conn, err := l.Accept()

		if err != nil {
			if config.Debug {
				logger.Error.Println(err)
			}

			continue
		}

		if config.Debug {
			logger.Info.Println("Accepted connection from IP:", conn.RemoteAddr())
		}

		i := NewInstance(conn)

		clients = append(clients, i)

		i.connection.Write(NewMessage(opcodes.ConnectSuccess))
		go i.listen()
		go i.heartbeat()
	}
}
