package hog

import (
	"bytes"
	"network/opcodes"
	"unicode/utf8"
)

type Message struct {
	Op   byte
	Args []string
}

func ParseMessage(buff []byte) Message {
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

	return Message{Op: o, Args: args}
}

func NewMessage(o byte, m ...string) []byte {
	b := []byte{o}

	for _, ele := range m {
		if len(b) > 1 {
			b = append(b, opcodes.Separator...)
		}

		b = append(b, []byte(ele)...)
	}

	return b
}
