package main

import (
	"fmt"
	"github.com/lambdaxs/go-server/discover"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	pool, err := discover.GetHTTPDialHostPool("127.0.0.1:8500", "test")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	client := http.Client{Timeout: time.Second}

	go func() {
		timer := time.NewTicker(time.Millisecond * 500)

		for range timer.C {
			host := pool.Get()
			resp, err := client.Get(fmt.Sprintf("http://%s/", host))
			if err != nil {
				fmt.Println("req err:" + err.Error())
				continue
			}
			buf, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			fmt.Println(host, string(buf))
		}
	}()

	select {}
}

func req(host string) {

}
