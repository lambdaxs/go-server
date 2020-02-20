package main

import (
    "fmt"
    "github.com/lambdaxs/go-server/discover"
    "time"
)

func main() {
    pool,err := discover.GetHTTPDialHostPool("127.0.0.1:8500","test")
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    go func() {
        timer := time.NewTicker(time.Second*1)

        for range timer.C{
            host := pool.Get()
            fmt.Println(host)
        }
    }()

    select {

    }
}
