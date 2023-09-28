package main

import (
	"time"
	"github.com/patbcole117/tinyC2/api"
	"github.com/patbcole117/tinyC2/beacon"
)

func main() {
    go api.Run()
	time.Sleep(5 * time.Second)
	b, err := beacon.NewBeacon("http://127.0.0.1:80/", "http")
	if err != nil {
		panic(err)
	}

	b.Run()
}