package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/patbcole117/tinyC2/comms"
	"github.com/patbcole117/tinyC2/beacon"
	"github.com/patbcole117/tinyC2/node"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNilNode = errors.New("node is nil")
	RUNNING = 1
	STOPPED = 0
)

type Dispatcher struct {
	Nodes 	[]node.Node
	db 		dbManager
	State	int
	exit 	chan bool
}

func NewDispatcher() Dispatcher {
	d := Dispatcher{
		State: STOPPED,
		db: NewDBManager(),
		exit: make(chan bool, 1),
	}
	d.Init()
	return d
}

func (d *Dispatcher) Stop() {
	fmt.Println("[D][+] Stop")
	d.exit <- true
}

func (d *Dispatcher) Run() {
	fmt.Println("[D][+] Run")
	d.State = RUNNING
	for {
		select {
		case <- d.exit:
			d.State = STOPPED
			return
		default:
			for i := range d.Nodes {
				select{
				case m := <- d.Nodes[i].Server.ChanUp:
					fmt.Printf("[D][<] %+v\n", m)
					d.Nodes[i].Hello = time.Now()
					res, err := db.UpdateNode(d.Nodes[i])
						if err != nil {
							fmt.Printf("[D][!] Run db.UpdateNode: %s\n", err.Error())
						} else if res.ModifiedCount == 0 {
							fmt.Println("[D][!] Run db.UpdateNode: no changes were made.")
						} else {
							fmt.Printf("[D][+] Node %d updated\n", d.Nodes[i].Id)
						}
					b, err := db.GetBeacon(m.From)
					if err == mongo.ErrNoDocuments {
						fmt.Printf("[D][+] Beacon %s does not exist.\n", m.From)
						b, err = beacon.NewBeacon(m.To, d.Nodes[i].Server.Type)
						if err != nil {
							fmt.Printf("[D][!] Run beacon.NewBeacon: %s\n", err.Error())
						}
						b.Name = m.From
						_, err = db.InsertBeacon(*b)
						if err != nil {
							fmt.Printf("[D][!] Run db.InsertBeacon: %s\n", err.Error())
						} else {
							fmt.Printf("[D][+] inserted %s\n", b.Name)
						}
					} else if err != nil {
						fmt.Printf("[D][!] Run db.GetBeacon: %s\n", err.Error())
					} else {
						b.Hello = time.Now()
						res, err := db.UpdateBeacon(*b)
						if err != nil {
							fmt.Printf("[D][!] Run db.UpdateBeacon: %s\n", err.Error())
						} else if res.ModifiedCount == 0 {
							fmt.Println("[D][!] Run db.UpdateBeacon: no changes were made.")
						} else {
							fmt.Printf("[D][+] Beacon %s updated\n", b.Name)
						}
					}
					d.Nodes[i].Server.ChanDown <- comms.Msg{From: m.To, To: m.From, Type: "", Ref: "", Content: ""}
				default:
					//fmt.Printf("[D][+] %s ChanUP is empty.\n", d.Nodes[i].Name)
					time.Sleep(1 * time.Second)
				}
			}
		}
	}
}

func (d *Dispatcher) AddNode(n node.Node) {
	fmt.Println("[D][+] AddNode")
	d.Nodes = append(d.Nodes, n)
}

func (d *Dispatcher) Init() error {
	fmt.Println("[D][+] Init")
	nodes, err := db.GetNodes()
	if err != nil {
		return err
	}
	d.Nodes = nodes
	for i := range d.Nodes {
		if d.Nodes[i].Status == node.LISTENING {
			if err := d.Nodes[i].StartSrv(); err != nil {return err}
		}
	}
	return nil
}

func (d *Dispatcher) RemoveNode(id int) error {
	fmt.Println("[D][+] RemoveNode")
	for i := range d.Nodes {
        if d.Nodes[i].Id == id {
           if err :=  d.Nodes[i].StopSrv(); err != nil {return err}
			d.Nodes = append(d.Nodes[:i], d.Nodes[i+1:]...)
			return nil
        } 
	}
    return ErrNilNode
}

func (d *Dispatcher) StartNode(id int) (node.Node, error) {
	fmt.Println("[D][+] StartNode")
	var n node.Node
	for i := range d.Nodes {
        if d.Nodes[i].Id == id {
           if err :=  d.Nodes[i].StartSrv(); err != nil {return n, err}
		   n = d.Nodes[i]
		   return n, nil
        }
    }
    return n, ErrNilNode
}

func (d *Dispatcher) StopNode(id int) (node.Node, error) {
	fmt.Println("[D][+] StopNode")
	var n node.Node
	for i := range d.Nodes {
        if d.Nodes[i].Id == id {
           if err :=  d.Nodes[i].StopSrv(); err != nil {return n, err}
		   n = d.Nodes[i]
		   return n, nil
        }
    }
    return n, ErrNilNode
}

func (d *Dispatcher)UpdateNode(id int, name, ip, port string) (node.Node, error) {
	fmt.Println("[D][+] UpdateNode")
	var n node.Node
	for i := range d.Nodes {
        if d.Nodes[i].Id == id {
			s := d.Nodes[i].Status
			if err := d.Nodes[i].StopSrv(); err != nil {return n, err}
			if name != ""{
				d.Nodes[i].Name = name
			}
			if ip != ""{
				d.Nodes[i].Ip = ip
			}
			if port!= ""{
				d.Nodes[i].Port = port
			}
			if s == node.LISTENING {
				if err := d.Nodes[i].StartSrv(); err != nil {return n, err}
			}
			n = d.Nodes[i]
			return n, nil
        }
    }
	return n, ErrNilNode
}