package hog

import (
	"config"
	"logger"
	"network/opcodes"
)

// Channel for messages to the entire room ONLY.
var messageQueue chan []byte
var clients []Instance

func init() {
	messageQueue = make(chan []byte, config.MessageQueueSize)
}

func enqueueMessage(message []byte) {
	messageQueue <- message
}

func joinMessage(name string) {
	enqueueMessage(NewMessage(opcodes.Join, name))
}

func exitMessage(name string) {
	enqueueMessage(NewMessage(opcodes.Leave, name))
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
		logger.Info.Println("Processing message: ", m)
	}
}
