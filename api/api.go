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
    r.Post("/v1/l/new", NewNode)

    http.ListenAndServe("127.0.0.1:8000", r)
}

func Check (w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"msg": "Good"}`))
}

func NewNode (w http.ResponseWriter, r *http.Request) {
    var resp string
    body, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    n := node.NewNode()
    err = json.Unmarshal(body, &n)
    if err != nil {
        panic(err)
    }
    result, err := db.InsertNewNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        resp = fmt.Sprintf(`{"ERROR": "%s"}`, err.Error())
    } else {
        w.WriteHeader(http.StatusCreated)
        resp = fmt.Sprintf(`{"INSERT": "%s"}`, result.InsertedID,)
    }
    w.Write([]byte(resp))
}
