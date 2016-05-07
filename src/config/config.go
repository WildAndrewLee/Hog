package config

import "flag"

var Debug bool

func init() {
	flag.BoolVar(&Debug, "debug", false, "Show debug information.")
}
