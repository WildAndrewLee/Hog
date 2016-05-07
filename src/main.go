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

	testBuffer := []byte{opcodes.SendMessage}

	for _, ele := range test {
		if len(testBuffer) > 1 {
			testBuffer = append(testBuffer, 0xFF, 0xFF)
		}

		testBuffer = append(testBuffer, []byte(ele)...)
	}

	logger.Info.Println(logger.FormatByteSlice(testBuffer))

	msg := hog.ParseMessage(testBuffer)

	fmt.Printf("%+v", msg)
}
