package main

import (
    "fmt"
    "github.com/bluele/gcache"
    "time"
)

func main(){
    cache := gcache.New(20000).LRU().LoaderExpireFunc(func(key interface{}) (i interface{}, duration *time.Duration, e error) {
        fmt.Println("load")
        if key == "a" {
            i = "abc"
            expire := time.Second*5
            duration = &expire
            return
        }
        return
    }).Build()


    go func() {
        ticker := time.NewTicker(time.Second*1)
        for range ticker.C {
            val,err := cache.Get("a")
            if err != nil {
                fmt.Println(err.Error())
                continue
            }
            fmt.Println(val)
        }
    }()

    time.Sleep(time.Second*20)

}
