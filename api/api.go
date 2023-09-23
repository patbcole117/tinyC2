package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/patbcole117/tinyC2/node"
)

var (
    db	            dbManager   = NewDBManager()
    url             string      = "127.0.0.1:8000"
    d               Dispatcher = NewDispatcher()
)

func Run() {
    r := chi.NewRouter()
    r.Get("/",                   Check)
    r.Get("/v1/l/",              GetNodes)
    r.Get("/v1/l/new/",          NewNode)
    r.Get("/v1/l/{id}/",         GetNode)
    r.Post("/v1/l/{id}/",        UpdateNode)
    r.Get("/v1/l/{id}/start/",   StartNode)
    r.Get("/v1/l/{id}/stop/",    StopNode)
    r.Get("/v1/l/{id}/x/",       DeleteNode)
    fmt.Println("[+] API Ready")
    http.ListenAndServe(url, r)
}

func Check (w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    jnodes, err := json.Marshal(d.Nodes)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("Check->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jnodes)
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("DeleteNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    _, err = db.DeleteNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("DeleteNode->db.DeleteNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if err := d.RemoveNode(id); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("DeleteNode->d.RemoveNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    w.WriteHeader(http.StatusOK)
    _, bmsg := FgoodMsg(fmt.Sprintf("deleted %d", id))
    w.Write(bmsg)
}

func NewNode (w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    n := node.NewNode()
    id, err := db.GetNextNodeID()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("NewNode->db.GetNextNodeID->"+err.Error())
        w.Write(bmsg)
        return
    }
    n.Id = id
    
    result, err := db.InsertNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("NewNode->db.InsertNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    d.AddNode(n)
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("inserted %s", result.InsertedID))
    w.Write(bmsg)
}

func GetNode(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, err := db.GetNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNode->db.GetNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    jnode, err := json.Marshal(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNode->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jnode)
}

func GetNodes(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    nodes, err := db.GetNodes()
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodes->db.GetNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    jnodes, err := json.Marshal(nodes)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("GetNodes->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(jnodes)
}

func StartNode(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StartNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, err := d.StartNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StartNode->d.StartNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    res, err := db.UpdateNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StartNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusOK)
    _, bmsg := FgoodMsg(fmt.Sprintf("started %d", n.Id))
    w.Write(bmsg)
}

func StopNode(w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, err := d.StopNode(id)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->d.StopNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    res, err := db.UpdateNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("StopNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusOK)
    _, bmsg := FgoodMsg(fmt.Sprintf("stopped %d", n.Id))
    w.Write(bmsg)

}

func UpdateNode (w http.ResponseWriter, r *http.Request) {
    LogRequest(r)
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->strconv.Atoi->"+err.Error())
        w.Write(bmsg)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->io.ReadAll->"+err.Error())
        w.Write(bmsg)
        return
    }
    
    var b map[string]string
    if err = json.Unmarshal(body, &b); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->json.Marshal->"+err.Error())
        w.Write(bmsg)
        return
    }

    n, err := d.UpdateNode(id, b["name"], b["ip"], b["port"]);
    if  err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->d.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    res, err := db.UpdateNode(n)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        _, bmsg := FbadMsg("UpdateNode->db.UpdateNode->"+err.Error())
        w.Write(bmsg)
        return
    }

    if res.ModifiedCount == 0 {
        _, bmsg := FgoodMsg("no changes were made.")
        w.Write(bmsg)
        return
    }
    w.WriteHeader(http.StatusCreated)
    _, bmsg := FgoodMsg(fmt.Sprintf("updated %d", n.Id))
    w.Write(bmsg)
}

//Helpers
func LogRequest(r *http.Request) {
    fmt.Println("[>]", time.Now().Format(time.RFC1123Z), r.RemoteAddr, r.Method, r.Host + r.URL.Path)
}
func FbadMsg (msg string) (string, []byte) {
    fmt.Println("[!]", msg)
	s := fmt.Sprintf(`{"type": "bad", "msg": "%s"}`, msg)
	return s, []byte(s)
}
func FgoodMsg (msg string) (string, []byte) {
    fmt.Println("[+]", msg)
    msg = strings.ReplaceAll(msg, `"`, `\"`)
	s := fmt.Sprintf(`{"type": "good", "msg": "%s"}`, msg)
	return s, []byte(s)
}