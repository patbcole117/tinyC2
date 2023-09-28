package comms

import (
    "errors"
    "net/http"
    "time"
)

var (
    CONN_TIMEOUT time.Duration = 10 * time.Second
    SERVER_DELAY = 2 * time.Second
    USER_AGENT string = "Mozilla/5.0 (Android 4.4; Mobile; rv:41.0) Geko/41.0 Firefox/41.0"
    ErrNilSrv error = errors.New("server is nil")
    ErrCommDNE error = errors.New("Comms package does not exist")
)

type CommsPackageTX interface {
    SendJSON(dst string, msg interface{}) (*http.Response, error)
    Get(dst string) (*http.Response, error)
}

func NewCommsPackageTX(c string) (CommsPackageTX, error) {
    switch c {
    case "http":
        tx := NewHTTPCommTX()
        return &tx, nil
    default:
        return nil, ErrCommDNE
    }
}