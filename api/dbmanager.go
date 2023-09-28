package api

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/patbcole117/tinyC2/beacon"
	"github.com/patbcole117/tinyC2/comms"
	"github.com/patbcole117/tinyC2/node"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pEnv 		string 						= os.Getenv("MONGO")
//var uri		string						= fmt.Sprintf("mongodb://root:%s@localhost:27017/?retryWrites=true&w=majority", p)
var uri			string 						= fmt.Sprintf("mongodb+srv://dev:%s@homenet-asia-mongodb-de.4sgvde0.mongodb.net/?retryWrites=true&w=majority", pEnv)
var serverAPI 	*options.ServerAPIOptions 	= options.ServerAPI(options.ServerAPIVersion1)
var opts 		*options.ClientOptions		= options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

type dbManager struct {
	c		*mongo.Client
}

//export MONGO=value
//$env:MONGO = "value"
func NewDBManager() dbManager {
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", 
		Value: 1}}).Err(); err != nil {
		fmt.Println("[!] Ping failed.")
		panic(err)
	}
	return dbManager{c: client}
}

func (db dbManager) DeleteNode(id int) (*mongo.DeleteResult, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	result, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		return nil, err
	}
	if  result.DeletedCount == 0 {
		return nil, errors.New(fmt.Sprintf("no document with id %d", id))
	}
	return result, nil
}

func (db dbManager) InsertMsg(msg comms.Msg) (*mongo.InsertOneResult, error) {
	coll := db.c.Database("tinyC2").Collection("Messages")
	res, err := coll.InsertOne(context.TODO(), msg)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db dbManager) InsertBeacon(b beacon.Beacon) (*mongo.InsertOneResult, error) {
	coll := db.c.Database("tinyC2").Collection("Beacons")
	res, err := coll.InsertOne(context.TODO(), b)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db dbManager) InsertNode(n node.Node) (*mongo.InsertOneResult, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	res, err := coll.InsertOne(context.TODO(), n)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (db dbManager) GetBeacon(n string) (*beacon.Beacon, error) {
	coll := db.c.Database("tinyC2").Collection("Beacons")

	var b beacon.Beacon
	err := coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: n}}).Decode(&b)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (db dbManager) GetNode(id int) (*node.Node, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")

	var n node.Node
	err := coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&n)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (db dbManager) GetMsgs() ([]comms.Msg, error) {
	coll := db.c.Database("tinyC2").Collection("Messages")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var msgs []comms.Msg
	if err = cursor.All(context.TODO(), &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (db dbManager) GetNodes() ([]node.Node, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	cursor, err := coll.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var nodes []node.Node
	if err = cursor.All(context.TODO(), &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (db dbManager) UpdateNode(n node.Node) (*mongo.UpdateResult, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")

	filter := bson.D{{Key: "id", Value: n.Id}}
	update := bson.D{{"$set", bson.D{{"name", n.Name}, {"ip", n.Ip}, {"port", n.Port},
		{"status", n.Status}, {"hello", n.Hello}}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db dbManager) UpdateBeacon(b beacon.Beacon) (*mongo.UpdateResult, error) {
	coll := db.c.Database("tinyC2").Collection("Beacons")

	filter := bson.D{{Key: "name", Value: b.Name}}
	update := bson.D{{"$set", bson.D{{"name", b.Name}, {"home", b.Home}, {"hello", b.Hello}}}}
	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db dbManager) GetNextNodeID() (int, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	opts := options.FindOne().SetSort(bson.D{{Key: "id", Value: -1}}) 

	var n node.Node
	if err := coll.FindOne(context.TODO(), bson.D{}, opts).Decode(&n); err != nil {
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}
	return n.Id+1, nil
}