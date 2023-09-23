package api

import (
	"errors"
	"fmt"

	"github.com/patbcole117/tinyC2/node"
)

var (
	DISPATCHER_CHAN_SIZE = 10
	ErrNilNode = errors.New("node is nil")
)

type Dispatcher struct {
	Nodes 	[]node.Node
	db 		dbManager
}

func NewDispatcher() Dispatcher {
	d := Dispatcher{
		db: NewDBManager(),
	}
	d.Init()
	return d
}

func (d *Dispatcher) AddNode(n node.Node) {
	d.Nodes = append(d.Nodes, n)
}

func (d *Dispatcher) Init() error {
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
	fmt.Println("[+] Dispatcher Ready")
	return nil
}

func (d *Dispatcher) RemoveNode(id int) error {
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
	var n node.Node
	for i := range d.Nodes {
		fmt.Println("[?]", id)
        if d.Nodes[i].Id == id {
           if err :=  d.Nodes[i].StartSrv(); err != nil {return n, err}
		   n = d.Nodes[i]
		   return n, nil
        }
    }
    return n, ErrNilNode
}

func (d *Dispatcher) StopNode(id int) (node.Node, error) {
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