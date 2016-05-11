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

	logger.Info.Println(logger.FormatByteSlice(buff))

	logger.Info.Println(n, err)

	e := make(chan bool)

	// Send heartbeat for 1 minute
	go func() {
		for x := 0; x < 20; x++ {
			select {
			case <-e:
				return
			default:
				l, err := conn.Write(hog.NewMessage(opcodes.Heartbeat))

				if err != nil || l == 0 {
					logger.Error.Println("Disconnected from server.")
					return
				}
				time.Sleep(3 * time.Second)
			}
		}

		logger.Info.Println("Done with heartbeat.")
	}()

	go func() {
		for {
			l, err := conn.Read(buff)

			if err != nil || l == 0 {
				logger.Error.Println("Disconnected from server.")
				e <- true
				return
			}
			logger.Info.Println("READ:", logger.FormatByteSlice(buff))
		}
	}()

	for x := 0; x < 12; x++ {
		conn.Write(hog.NewMessage(opcodes.SendMessage, "Hello World!"))
		time.Sleep(5 * time.Second)
	}

	time.Sleep(3 * time.Minute)

	conn.Close()
}
