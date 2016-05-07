package logger

import (
	"fmt"
	"strings"
)

func FormatByteSlice(bytes []byte) string {
	repr := []string{}

	for _, b := range bytes {
		repr = append(repr, fmt.Sprintf("%#x", b))
	}

	return fmt.Sprintf("[%s]", strings.Join(repr, " "))
}
