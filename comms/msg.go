package comms

import (
	"math/rand"
)

type Msg struct {
	From            string  `bson:"from"	    json:"from"`
    To              string  `bson:"to"	        json:"to"`
	Type            string  `bson:"type"	    json:"type"`
	Ref             string  `bson:"ref"	        json:"ref"`
	Content         string  `bson:"content"	    json:"content"`
}

func NewJob(t, c string) Msg {
	return Msg{
		To: t,
		Type: "job",
		Ref: GetRef(12),
		Content: "queued",
	}
}

func GetRef(sz int) string {
	abc := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	n := make([]rune, sz)
	for i := range n {
		n[i] = abc[rand.Intn(len(abc))]
	}
	return string(n)
}