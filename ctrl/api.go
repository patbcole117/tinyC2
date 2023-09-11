package ctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/patbcole117/tinyC2/node"
)

const ver int = 1
var dbh dbHandler = NewDBHandler()

func Run() {
    r := chi.NewRouter()

    r.Get("/v1", v1Get)
    r.Get("/v1/l", v1GetListenersAll)
    r.Post("/v1/l/new", v1NewListener)
    r.Route("/v1/l/{id}", func(r chi.Router) {
        r.Use(v1ListenerByIdCtx)
        r.Delete("/", v1DeleteListenerById)
        r.Get("/", v1GetListenerById)
        r.Put("/", v1DeleteListenerById)
    })

    http.ListenAndServe("127.0.0.1:8000", r)
}

func v1NewListener (w http.ResponseWriter, r *http.Request) {
    n := node.NewNode()
    err := json.NewDecoder(r.Body).Decode(&n)
    if err != nil {
        panic(err)
    }
    dbh.dbInsertListener(n)
    res := fmt.Sprintf(`{"message": "New listener inserted."}`)
    w.Write([]byte(res))
}

func v1Get (w http.ResponseWriter, r *http.Request) {
    res := fmt.Sprintf(`{"version": "%d"}`, ver)
    w.Write([]byte(res))
}

func v1GetListenersAll (w http.ResponseWriter, r *http.Request) {
    res := `{"TODO": "dbGetAllListeners"}`
    w.Write([]byte(res))
}

func v1ListenerByIdCtx (next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        listener, err := dbh.dbGetListenerById(id)
        if err != nil {
            http.Error(w, http.StatusText(404), 404)
            return
        }
        ctx := context.WithValue(r.Context(), "listener", listener)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func v1DeleteListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := fmt.Sprintf(`{"TODO": "dbDeleteListenerById(%s)"}`, id)
    w.Write([]byte(res))
}

func v1GetListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := fmt.Sprintf(`{"TODO": "dbGetListenerById(%s)"}`, id)
    w.Write([]byte(res))
}
