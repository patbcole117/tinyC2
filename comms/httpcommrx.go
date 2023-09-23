package comms

import (
    "encoding/json"
    "fmt"
    "io"
	"net/http"
    "time"
)

type HTTPCommRX struct {
    Ip      string
    Port    string
    Srv       *http.Server
}
func NewHTTPCommRX(i, p string) *HTTPCommRX {
    rx := HTTPCommRX{Ip: i, Port: p}
    rx.Srv = rx.ProvisionSrv()
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
    mux.HandleFunc("/", getRoot)
    return &http.Server{
            Addr: addy,
            Handler: mux,
            ReadTimeout: CONN_TIMEOUT,
            WriteTimeout: CONN_TIMEOUT,     
    } 
}

func getRoot (w http.ResponseWriter, r *http.Request) {
    msg, _ := io.ReadAll(r.Body)
    fmt.Println("Recieved:")
    fmt.Printf("%s\n", msg)
    fmt.Println("Reply:")
    reply := map[string]string{
        "From": "Home",
        "Type": "job",
        "Ref": "1",
        "Content":"Do the thing"}    
    breply, _ := json.Marshal(reply) 
    fmt.Println(string(breply))
    w.Write(breply)
}

func urlEcho(w http.ResponseWriter, r *http.Request) {
    var msg []byte
    if r.Body != http.NoBody {
        msg, _ = io.ReadAll(r.Body)
    } else {
        msg = []byte("No body in request.")
    }
    h, _ := json.MarshalIndent(r.Header, "", "  ") 
    reply := fmt.Sprintf("%s\n%s", msg, h)
    w.Write([]byte(reply))
}
