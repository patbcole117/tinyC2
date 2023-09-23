package node

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/patbcole117/tinyC2/beacon"
	"github.com/patbcole117/tinyC2/comms"
)

const (
	ERROR    	= "-1"
	STOPPED     = "0"
	LISTENING   = "1"

	SERVER_DEFAULT_IP      		string        	= "127.0.0.1"
	SERVER_DEFAULT_PORT    		string          = "80"
	SERVER_DEFAULT_TIMEOUT 		time.Duration 	= 3 * time.Second
	SERVER_DEFAULT_NAME_SIZE 	int 			= 10
	NODE_DEFAULT_CHAN_SIZE 		int 			= 10
)

type Node struct {
	Id 			int					`bson:"id" 		json:id`
	Name   		string				`bson:"name"	json:"name"`
	Ip     		string   			`bson:"ip"		json:"ip"`
	Port   		string      		`bson:"port"	json:"port"`
	Status 		string      		`bson:"status"	json:"status"`
    Dob     	time.Time			`bson:"dob"		json:"dob"`
	Hello 		time.Time			`bson:"hello"	json:"hello"`
	ChanUp		chan beacon.Msg		`bson:"-" 		json:"-`
	ChanDown	chan beacon.Msg		`bson:"-" 		json:"-`
	Server 		*comms.HTTPCommRX   `bson:"-" 		json:"-"`
}

func NewNode() Node {
	n := Node {
	Ip:			SERVER_DEFAULT_IP,
	Port:		SERVER_DEFAULT_PORT,
	Status:		STOPPED,
	ChanUp: 	make(chan beacon.Msg, NODE_DEFAULT_CHAN_SIZE),
	ChanDown: 	make(chan beacon.Msg, NODE_DEFAULT_CHAN_SIZE),
	}
    n.Dob = time.Now()
	n.Hello = time.Now()
	n.Server = comms.NewHTTPCommRX(n.Ip, n.Port)
	n.initName(SERVER_DEFAULT_NAME_SIZE)
	return n
}

func (n *Node) initName(sz int) {
	abc := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	b := make([]rune, sz)
	for i := range b {
		b[i] = abc[rand.Intn(len(abc))]
	}
	n.Name = string(b)
}

func (n *Node) StartSrv() error {
	n.Server = comms.NewHTTPCommRX(n.Ip, n.Port)
	if err := n.Server.StartSrv(); err != nil {
		return err
	}
	n.Status = LISTENING
	return nil
}

func (n *Node) StopSrv() error {
	if err := n.Server.StopSrv(); err != nil {
		return err
	}
	n.Status = STOPPED
	return nil
}

func (n *Node) UnmarshalJSON(j []byte) error {
    type Alias Node
    aux := &struct {
        Dob 	string	`json:"dob"`
		Hello 	string	`json:"hello"`	
        *Alias
    }{
        Alias:  (*Alias)(n),
    }

    if err := json.Unmarshal(j, &aux); err != nil {
        return err
    }
   
    t, err := time.Parse(time.RFC3339, aux.Dob)
    if err != nil {
        return err
    }
	n.Dob = t
    
	t, err = time.Parse(time.RFC3339, aux.Hello)
    if err != nil {
        return err
    }
	n.Hello 	= t

	n.ChanUp 	= make(chan beacon.Msg, NODE_DEFAULT_CHAN_SIZE)
	n.ChanDown 	= make(chan beacon.Msg, NODE_DEFAULT_CHAN_SIZE)
	n.Server = comms.NewHTTPCommRX(n.Ip, n.Port)
    return nil
}

func (n *Node) MarshalJSON() ([]byte, error) {
    type Alias Node
    return json.Marshal(&struct {
        Dob 		string 	`bson:"dob"		json:"dob"`
		Hello 		string	`bson:"hello"	json:"hello"`
		ChanDown 	string	`bson:"-" 		json:"-`
		ChanUp 		string	`bson:"-" 		json:"-`
        *Alias
    }{
        Dob: 		n.Dob.Format(time.RFC3339),
		Hello: 		n.Hello.Format(time.RFC3339),
		ChanDown: 	"-",
		ChanUp: 	"-",
        Alias:  (*Alias)(n),
    })
}

func (n *Node) ToJsonPretty() (string) {
	b, err := json.MarshalIndent(n, "", "    ")
	if err != nil {
		return string(make([]byte, 0))
	}
	return string(b)
}

func (n *Node) urlRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Date", time.Now().Format(time.UnixDate))
	msg := fmt.Sprintf("[+] Hello from %s.\n", n.Name)
	io.WriteString(w, msg)
}