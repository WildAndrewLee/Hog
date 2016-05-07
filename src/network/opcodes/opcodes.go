package opcodes

// Go y u no hav slice constants.
var Separator = []byte{0xFF, 0xFF}

const (
	SendMessage    byte = 0x01
	ReceiveMessage byte = 0x02
	Join           byte = 0x03
	Leave          byte = 0x04
	Heartbeat      byte = 0x05
)
