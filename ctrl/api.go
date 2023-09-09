package ctrl
/*
import (
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/patbcole117/tinyC2/node"
)

const ver int = 1

func Run() {
    r := chi.NewRouter()

    r.Get("/v1", getV1)
    r.Get("/v1/l", getV1Listeners)
    r.Route("/v1/l/{id}", func(r chi.Router) {
        r.Use(getV1ListenerByIdCtx)
        r.Delete("/", deleteV1ListenerById)
        r.Get("/", getV1ListenerById)
        r.Put("/", putV1ListenerById)
    })

    http.ListenAndServe("127.0.0.1:8000", r)
}

func getV1 (w http.ResponseWriter, r *http.Request) {
    res := fmt.Sprintf(`{"version": "%d"}`, ver)
    w.Write([]byte(res))
}

func getV1Listeners (w http.ResponseWriter, r *http.Request) {
    res := `{"TODO": "SELECT * FROM LISTENERS;"}`
    w.Write([]byte(res))
}

func getV1ListenerByIdCtx (next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        listener, err := dbGetListenerById(id)
        if err != nil {
            http.Error(w, http.StatusText(404), 404)
            return
        }
        ctx := context.WithValue(r.Context(), "listener", listener)
        next.ServeHTTP(w, t.WithContext(ctx))
    })
}

func deleteV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := `{"TODO": "SELECT * FROM LISTENERS WHERE ID = x;"}`
    w.Write([]byte(res))
}

func getV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := `{"TODO": "SELECT * FROM LISTENERS WHERE ID = x;"}`
    w.Write([]byte(res))
}

func putV1ListenerById (w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    res := `{"TODO": "SELECT * FROM LISTENERS WHERE ID = x;"}`
    w.Write([]byte(res))
}

func dbGetListenerById(id int) []byte {
    n := node.NewNode()
    b, err :=  json.Marshal(*n)
    if err != nil {
        
    }
    return b
}
*/