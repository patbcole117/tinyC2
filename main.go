package main

import (
	//"github.com/patbcole117/tinyC2/ui"
	//"fmt"

	"github.com/patbcole117/tinyC2/ctrl"
	//"github.com/patbcole117/tinyC2/node"
)

func main() {
//	ui.KickOff()
    ctrl.Run()
/*
	dbh := ctrl.NewDbHandler()
	n := node.NewNode()
	fmt.Println(n.ToJsonPretty())
	dbh.InsertListenerDoc(n)
	dbh.Disconnect()
    */
}
