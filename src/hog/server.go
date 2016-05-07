package hog

func JoinMessage(name string) {
	// Use Join opcode.
	// [0x3, <name in bytes>...]
}

func ExitMessage(name string) {
	// Use Leave opcode.
	// [0x4, <name in bytes>...]
}
