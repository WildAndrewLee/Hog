package main

import (
	"flag"
	"fmt"
	"hog"
)

func main() {
	flag.Parse()

	test := []string{"Andrew"}

	testBuffer := []byte{0x1}

	for _, ele := range test {
		if len(testBuffer) > 0 {
			testBuffer = append(testBuffer, 0xFF, 0xFF)
		}

		testBuffer = append(testBuffer, []byte(ele)...)
	}

	msg := hog.ParseMessage(testBuffer)

	fmt.Printf("%+v", msg)
}
