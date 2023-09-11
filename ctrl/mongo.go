package ctrl

import (

	"context"
	"fmt"
	"os"

	"github.com/patbcole117/tinyC2/node"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbHandler struct {
	con	mongo.Client
}
//export MONGO=value
//$env:MONGO = "value"
func NewDBHandler() dbHandler{
	p := os.Getenv("MONGO")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://dev:%s@homenet-asia-mongodb-de.4sgvde0.mongodb.net/?retryWrites=true&w=majority", p)
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

	return dbHandler{con: *client}
}

func (h dbHandler) dbInsertListener(n node.Node) {
	coll := h.con.Database("tinyC2").Collection("Listeners")
	_, err := coll.InsertOne(context.TODO(), n)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
}

func (h dbHandler) dbUpdateListenerById(n node.Node) {

}

func (h dbHandler) dbDeleteListenerById(id string) {

}

func (h dbHandler) dbGetListenerById(id string) ([]byte, error) {
    return nil, nil
}

func (h dbHandler) dbDisconnect() {
	if err := h.con.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
