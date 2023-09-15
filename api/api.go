package api

import (
	"encoding/json"
    "fmt"
    "io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patbcole117/tinyC2/node"
)

var db dbConnection = GetClient()

func Run() {
    r := chi.NewRouter()
    r.Get("/", Check)
    r.Get("/v1/l", GetAllNodes)
    r.Post("/v1/l/new", NewNode)
    r.Post("/v1/l/delete", DeleteNode)

    http.ListenAndServe("127.0.0.1:8000", r)
}

func Check (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"msg": "Good"}`))
}

func NewNode (w http.ResponseWriter, r *http.Request) {
    var msg string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err.Error())
        w.Write([]byte(msg))
        return
    }

    n := node.NewNode()
    err = json.Unmarshal(body, &n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "json.Unmarshal", "Msg": "%s"}`, err.Error())
        w.Write([]byte(msg))
        return
    }

    result, err := db.InsertNewNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "db.InsertNewNode", "Msg": "%s"}`, err.Error())
    } else {
        w.WriteHeader(http.StatusCreated)
        msg = fmt.Sprintf(`{"INSERT": "%s"}`, result.InsertedID,)
    }
    w.Write([]byte(msg))
}

func DeleteNode (w http.ResponseWriter, r *http.Request) {
    var msg string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "io.ReadAll", "Msg": "%s"}`, err.Error())
        w.Write([]byte(msg))
        return
    }

    result, err := db.DeleteNode(string(body))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "db.DeleteNode", "Msg": "%s"}`, err.Error())
    } else {
        w.WriteHeader(http.StatusCreated)
        msg = fmt.Sprintf(`{"DELETED": "%d": "%s"}`, result.DeletedCount, string(body))
    }
    w.Write([]byte(msg))
}

func GetAllNodes(w http.ResponseWriter, r *http.Request) {
    var msg string
    var nodes []node.Node
    nodes, err := db.GetAllNodes()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "InsertNewNode", "Msg": "%s"}`, err.Error())
        w.Write([]byte(msg))
        return
    }

    resp, err := json.Marshal(nodes)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        msg = fmt.Sprintf(`{"ERROR": "%s"}`, err.Error())
        w.Write([]byte(msg))
        return
    }
    w.Write(resp)
}
