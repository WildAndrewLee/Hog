package logger

import (
	"config"
	"fmt"
	"strings"
	"time"
)

func formatByteSlice(bytes []byte) string {
	repr := []string{}

	for _, b := range bytes {
		repr = append(repr, fmt.Sprintf("%#x", b))
	}

	return fmt.Sprintf("[%s]", strings.Join(repr, " "))
}

func now() string {
	return time.Now().Format(time.RFC1123)
}

func debugPrefix() {
	fmt.Printf("[%s] DEBUG: ", now())
}

func LogString(message string) {
	if !config.Debug {
		return
	}

	debugPrefix()
	fmt.Println(message)
}

func LogBytes(bytes []byte) {
	if !config.Debug {
		return
	}

	debugPrefix()
	fmt.Println(formatByteSlice(bytes))
}
