package api

import (

	"context"
	"fmt"
	"os"

	"github.com/patbcole117/tinyC2/node"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbConnection struct {
	c *mongo.Client
}

//export MONGO=value
//$env:MONGO = "value"
func GetClient() dbConnection {
	p := os.Getenv("MONGO")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://dev:%s@homenet-asia-mongodb-de.4sgvde0.mongodb.net/?retryWrites=true&w=majority", p)
    //uri := fmt.Sprintf("mongodb://root:%s@localhost:27017/?retryWrites=true&w=majority", p)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	//fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return dbConnection{c: client}
}

func (db dbConnection) InsertNewNode(n node.Node) (*mongo.InsertOneResult, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	result, err := coll.InsertOne(context.TODO(), n)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db dbConnection) DeleteNode(id string) (*mongo.DeleteResult, error) {

	coll := db.c.Database("tinyC2").Collection("Listeners")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	result, err := coll.DeleteOne(context.TODO(), bson.D{{"_id", oid}})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (db dbConnection) GetAllNodes() ([]node.Node, error) {
	coll := db.c.Database("tinyC2").Collection("Listeners")
	filter := bson.D{}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	var nodes []node.Node
	if err = cursor.All(context.TODO(), &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}
