package main

import (
    "fmt"
    "net/url"
)

func main(){
    str := "username:password@/dbname?charset=utf8&parseTime=True&loc=Local&readTimeout=3s&writeTime=3s"
    rs,err := updateDSNQuery(str, map[string]string{
        "a":"2",
        "readTimeout":"1s",
    })
    if err != nil {
        return
    }
    fmt.Println(rs)
}

func updateDSNQuery(dsn string, kv map[string]string) (string,error) {
    u,err := url.Parse(dsn)
    if err != nil {
        return "",err
    }
    q := u.Query()
    for k,v := range kv {
        q.Set(k, v)
    }
    u.RawQuery = q.Encode()
    return u.String(),nil
}
