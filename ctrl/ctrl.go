package ctrl

import (
    "fmt"
    "net/http"

    "github.com/go-chi/chi/v5"
)

func Run() {
    r := chi.NewRouter()

    r.Get("/v1", v1)
    r.Get("/v1/l", v1Lis)
    r.Get("/v1/l/new", v1LisNew)
    r.Get("/v1/l/:id", v1LisId)
    r.Get("/v1/l/:id/delete", v1LisIdDelete)
    r.Get("/v1/l/:id/update", v1LisIdUpdate)
    r.Get("/v1/l/:id/c", v1LisIdC)

    http.ListenAndServe("127.0.0.1:8000", r)
}

func v1 (w http.ResponseWriter, r *http.Request) error {
    w.Write([]byte())
}

func v1Lis (w http.ResponseWriter, r *http.Request) error {
    w.Write()
}

func v1LisNew (w http.ResponseWriter, r *http.Request) error {
    w.Write( c.JSON(&fiber.Map{"data": "v1LisNew"})
}

func v1LisId (w http.ResponseWriter, r *http.Request) error {
    data := fmt.Sprintf(c.Params("id"))
    w.Write( c.JSON(&fiber.Map{"data": data})
}

func v1LisIdDelete (w http.ResponseWriter, r *http.Request) error {
    w.Write( c.JSON(&fiber.Map{"data": "v1LisIdDelete"})
}

func v1LisIdUpdate (w http.ResponseWriter, r *http.Request) error {
    w.Write( c.JSON(&fiber.Map{"data": "v1LisIdUpdate"})
}

func v1LisIdC (w http.ResponseWriter, r *http.Request) error {
    w.Write( c.JSON(&fiber.Map{"data": "v1LisIdC"})
}
