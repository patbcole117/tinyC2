package main

import (

	"github.com/patbcole117/tinyC2/ui"
	"github.com/patbcole117/tinyC2/ctrl"
)

func main() {
	go ctrl.Run()
	ui.KickOff()
}
