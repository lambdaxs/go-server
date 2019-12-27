package driver

import (
    "io/ioutil"
    "net/http"
)

type HttpClient struct {
    *http.Client
}

func (h *HttpClient)Do(req *http.Request) (data []byte, err error){
    resp,err := h.Client.Do(req)
    if err != nil {
        return
    }
    buf,err := ioutil.ReadAll(resp.Body)
    defer resp.Body.Close()
    if err != nil {
        return
    }
    return buf,nil
}


