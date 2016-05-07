package hog

import "config"

var messageQueue chan []byte

func init() {
	messageQueue = make(chan []byte, config.MessageQueueSize)
}

func EnqueueMessage(message []byte) {
	messageQueue <- message
}

func JoinMessage(name string) {
	// Use Join opcode.
	// [0x3, <name in bytes>...]
}

func ExitMessage(name string) {
	// Use Leave opcode.
	// [0x4, <name in bytes>...]
}
