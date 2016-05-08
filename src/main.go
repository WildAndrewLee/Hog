package main

import (
	"flag"
	"fmt"
	"hog"
	"logger"
	"network/opcodes"
)

func main() {
	flag.Parse()

	test := []string{"Andrew"}

	testBuffer := hog.NewMessage(opcodes.SendMessage, test...)

	logger.Info.Println(logger.FormatByteSlice(testBuffer))

	msg := hog.ParseMessage(testBuffer)

	fmt.Printf("%+v", msg)
}
