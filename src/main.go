package main

import (
	"flag"
	"fmt"
	"hog"
)

func main() {
	flag.Parse()

	test := []string{"hello", "there", "こんにちは"}

	testBuffer := []byte{0x1}

	for _, ele := range test {
		testBuffer = append(testBuffer, []byte(ele)...)
		testBuffer = append(testBuffer, 0xFF, 0xFF)
	}

	msg := hog.ParseMessage(testBuffer)

	fmt.Printf("%+v", msg)
}
