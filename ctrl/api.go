package ctrl

import (
    "context"
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/patbcole117/tinyC2/node"
)

const ver int = 1
var dbh dbHandler = NewDBHandler()

func Run() {
    r := chi.NewRouter()

    r.Get("/v1", getV1)
    r.Get("/v1/l", getV1ListenersAll)
    r.Route("/v1/l/{id}", func(r chi.Router) {
        r.Use(v1ListenerByIdCtx)
        r.Delete("/", deleteV1ListenerById)
        r.Get("/", getV1ListenerById)
        r.Put("/", updateV1ListenerById)
    })

    http.ListenAndServe("127.0.0.1:8000", r)
}

func getV1 (w http.ResponseWriter, r *http.Request) {
    res := fmt.Sprintf(`{"version": "%d"}`, ver)
    w.Write([]byte(res))
}

func getV1ListenersAll (w http.ResponseWriter, r *http.Request) {
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

func deleteV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := fmt.Sprintf(`{"TODO": "dbDeleteListenerById(%s)"}`, id)
    w.Write([]byte(res))
}

func getV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := fmt.Sprintf(`{"TODO": "dbGetListenerById(%s)"}`, id)
    w.Write([]byte(res))
}

func updateV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := fmt.Sprintf(`{"TODO": "dbUpdateListenerById(%s)"}`, id)
    w.Write([]byte(res))
}

