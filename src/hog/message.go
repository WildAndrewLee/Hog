package hog

import (
	"bytes"
	"network/opcodes"
	"unicode/utf8"
)

type message struct {
	Op   byte
	Args []string
}

func ParseMessage(buff []byte) message {
	o := buff[0]
	args := []string{}

	buff = buff[1:]

	arg := []rune{}

	for len(buff) > 0 {
		/*
			0xFFFF is an unused codepoint in unicode.
			We will use it to delimit operation arguments,
			[0x11, 0x10, 0xFF, 0xFF, 0x3C, 0x49] => [[0x11, 0x10], [0x3C, 0x49]]
		*/
		if len(buff) >= 2 {
			if bytes.Equal(buff[0:2], opcodes.Separator) {
				args = append(args, string(arg))
				arg = []rune{}

				buff = buff[2:]

				if len(buff) == 0 {
					break
				}
			}
		}

		token, size := utf8.DecodeRune(buff)

		arg = append(arg, token)
		buff = buff[size:]
	}

	if len(arg) > 0 {
		args = append(args, string(arg))
	}

	return message{Op: o, Args: args}
}

/*
Constructs a byte slice given a message op code and
a variable number of string arguments. Truncates the byte
array in the event that a message is too long.
*/
func NewMessage(o byte, m ...string) []byte {
	b := []byte{o}

	for _, ele := range m {
		if len(b) > 1 {
			b = append(b, opcodes.Separator...)
		}

		b = append(b, []byte(ele)...)
	}

	l := len(b)

	if l > 0xFFFF+2 {
		b = b[:0xFFFF]
		l = 0xFFFF
	}

	b = append([]byte{byte(l >> 16), byte(l & 0xFFFF)}, b...)

	return b
}
