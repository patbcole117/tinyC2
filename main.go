package main

import (
	"github.com/patbcole117/tinyC2/api"
	"github.com/patbcole117/tinyC2/ui"
)

func main() {
	go api.Run()
	ui.KickOff()
}
