package comms

import (
    "bytes"
    "encoding/json"
	"net/http"
    "time"
)

type HTTPCommTX struct {
	C *http.Client
}

func NewHTTPCommTX() HTTPCommTX {
    return HTTPCommTX{C: &http.Client{Timeout: CONN_TIMEOUT}}
}

func (tx *HTTPCommTX) SendJSON(dst string, msg interface{})(*http.Response, error) {
    b, err := json.Marshal(msg)
    if err != nil {
        return nil, err
    } 
    body := bytes.NewReader(b)

    req, err := http.NewRequest(http.MethodPost, dst, body)
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", USER_AGENT)
    req.Header.Set("Date", time.Now().Format(time.RFC1123))
    
    res, err := tx.C.Do(req)
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (tx *HTTPCommTX) Get(dst string) (*http.Response, error) {
    req, err := http.NewRequest(http.MethodGet, dst, nil)
    if err != nil {
        return nil, err
    }

    req.Header.Set("User-Agent", USER_AGENT)
    req.Header.Set("Date", time.Now().Format(time.RFC1123))
    
    res, err := tx.C.Do(req)
    if err != nil {
        return nil, err
    }
    return res, nil
}