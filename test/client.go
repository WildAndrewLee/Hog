package main

import (
	"config"
	"hog"
	"logger"
	"net"
	"network/opcodes"
	"strconv"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(config.Port))

	if err != nil {
		logger.Error.Println("Unable to connect to server.")
		return
	}

	n, err := conn.Write(hog.NewMessage(opcodes.Connect, "Andrew"))

	buff := make([]byte, 64)
	conn.Read(buff)

	logger.Info.Println(buff)

	logger.Info.Println(n, err)

	for x := 0; x < 20; x++ {
		l, err := conn.Write(hog.NewMessage(opcodes.Heartbeat))

		if err != nil || l == 0 {
			logger.Error.Println("Disconnected from server.")
			return
		}

		logger.Info.Println("Send heartbeat.")
		time.Sleep(3 * time.Second)
	}

	time.Sleep(3 * time.Minute)

	conn.Close()
}
