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
    r.Post("/v1/l/update", UpdateNode)

    http.ListenAndServe("127.0.0.1:8000", r)
}

func Check (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"CHECK": "GOOD"}`))
}

func DeleteNode (w http.ResponseWriter, r *http.Request) {
    var wmsg string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("DeleteNode:io.ReadAll %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }

    result, err := db.DeleteNode(string(body))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("DeleteNode:db.DeleteNode %s", err.Error())
    } else if result.DeletedCount == 0 {
        w.WriteHeader(http.StatusCreated)
        wmsg = "NO MATCH"
    } else {
        w.WriteHeader(http.StatusCreated)
        wmsg = fmt.Sprintf(`%d %s`, result.DeletedCount, string(body) )
    }
    w.Write([]byte(wmsg))
}

func GetAllNodes(w http.ResponseWriter, r *http.Request) {
    var wmsg string
    var nodes []node.Node
    nodes, err := db.GetAllNodes()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("GetAllNodes:db.GetAllNodes %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }
    resp, err := json.Marshal(nodes)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("GetAllNodes:json.Marshal %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }
    w.Write(resp)
}

func NewNode (w http.ResponseWriter, r *http.Request) {
    var wmsg string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("NewNode:io.ReadAll %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }

    n := node.NewNode()
    err = json.Unmarshal(body, &n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("NewNode:json.Unmarshal %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }

    result, err := db.InsertNewNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("NewNode:db.InsertNewNode %s", err.Error())
    } else {
        w.WriteHeader(http.StatusCreated)
        wmsg = fmt.Sprintf("%s", result.InsertedID)
    }
    w.Write([]byte(wmsg))
}

func UpdateNode (w http.ResponseWriter, r *http.Request) {
    var wmsg string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("UpdateNode:io.ReadAll %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }

    n := node.NewNode()
    err = json.Unmarshal(body, &n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("UpdateNode:json.Unmarshal %s", err.Error())
        w.Write([]byte(wmsg))
        return
    }

    result, err := db.UpdateNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        wmsg = fmt.Sprintf("UpdateNode:db.UpdateNode %s", err.Error())
    } else if result.ModifiedCount == 0 {
        w.WriteHeader(http.StatusCreated)
        wmsg =  "NO MATCH"
    } else {
        w.WriteHeader(http.StatusCreated)
        wmsg = fmt.Sprintf(`%d %s`, result.ModifiedCount, n.Id)
    }
    w.Write([]byte(wmsg))
}