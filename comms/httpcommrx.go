package comms

import (
    "encoding/json"
    "fmt"
    "io"
	"net/http"
    "time"
)

var HTTPRX_CHAN_SIZE = 10

type HTTPCommRX struct {
    Ip          string
    Port        string
    Srv         *http.Server
    Type        string
    ChanDown    chan Msg
    ChanUp      chan Msg
}

func NewHTTPCommRX(i, p string) (*HTTPCommRX) {
    rx := HTTPCommRX{Ip: i, Port: p, Type: "http"}
    rx.Srv = rx.ProvisionSrv()
    rx.ChanDown = make(chan Msg, HTTPRX_CHAN_SIZE)
    rx.ChanUp = make(chan Msg, HTTPRX_CHAN_SIZE)
    return &rx
}

func (rx *HTTPCommRX) StartSrv() error {
    if rx.Srv == nil {
        return ErrNilSrv
    }
    go rx.Srv.ListenAndServe()
    time.Sleep(SERVER_DELAY)
    return nil
}

func (rx *HTTPCommRX) StopSrv() error {
    if rx.Srv == nil {
        return ErrNilSrv
    }
    if err := rx.Srv.Close(); err != nil {
        return err
    }
    rx.Srv = rx.ProvisionSrv()
    return nil
}

func (rx *HTTPCommRX) GetAddy() string {
    return (rx.Ip + ":" + rx.Port)
}

func (rx *HTTPCommRX) ProvisionSrv() *http.Server {
    addy := rx.GetAddy()
    mux := http.NewServeMux()
    mux.HandleFunc("/", handle(rx))
    return &http.Server{
            Addr: addy,
            Handler: mux,
            ReadTimeout: CONN_TIMEOUT,
            WriteTimeout: CONN_TIMEOUT,     
    } 
}

func handle( rx *HTTPCommRX) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        body, err := io.ReadAll(r.Body)
        if err != nil {
            fmt.Printf("[!] getRoot %s\n", err.Error())
        }
        var m Msg
        err = json.Unmarshal(body, &m); if err != nil {
            fmt.Printf("[!] getRoot %s\n", err.Error())
        }
        rx.ChanUp <- m
        resp := <- rx.ChanDown
        bresp, err := json.Marshal(resp)
        if err != nil {
            fmt.Printf("[!] getRoot %s\n", err.Error())
        }
        w.Write(bresp)
    }
}
