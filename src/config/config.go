package config

import "flag"

var Debug bool
var MessageQueueSize int
var Port int

func init() {
	flag.BoolVar(&Debug, "debug", dDefault, "Show debug information.")
	flag.IntVar(&MessageQueueSize, "queue", mQSDefault, "Sets the message queue size to be used by the server.")
	flag.IntVar(&Port, "port", pDefault, "Sets the port used for incoming connections.")
}
