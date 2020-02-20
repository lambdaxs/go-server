package discover

import (
    "fmt"
    "math/rand"
    "sync"
    "log"
    "time"

    consulapi "github.com/hashicorp/consul/api"
)

type HttpHostPool struct {
    address string
    serviceName string
    client *consulapi.Client

    sync.RWMutex
    addrs  []string
    length int
    lastIndex uint64
    t *time.Ticker
}

func (h *HttpHostPool) Get() string {
    h.RLock()
    defer h.RUnlock()

    rand.Shuffle(h.length, func(i, j int) {
        h.addrs[i], h.addrs[j] = h.addrs[j], h.addrs[i]
    })
    if h.length > 0 {
        return h.addrs[0]
    }
    return ""
}

func (h *HttpHostPool)watch() {
    for {
        select {
        case <-h.t.C:
            addrs, updated,err := h.resolver()
            if err != nil {
                return
            }
            if updated {
                h.Lock()
                h.addrs = addrs
                h.length = len(addrs)
                h.Unlock()
            }
        }
    }
}

func (h *HttpHostPool)resolver() (list []string, updated bool,err error){
    serviceEntries, metadata, err := h.client.Health().Service(fmt.Sprintf("http:%s",h.serviceName), "", true, &consulapi.QueryOptions{})
    if err != nil {
        return
    }
    if metadata.LastIndex != h.lastIndex {
        updated = true
    }else {
        updated = false
    }
    h.lastIndex = metadata.LastIndex

    list = make([]string, 0)
    for _, serviceEntry := range serviceEntries {
        list = append(list, fmt.Sprintf("%s:%d", serviceEntry.Service.Address, serviceEntry.Service.Port))
    }
    return
}

func GetHTTPDialHostPool(address string, serviceName string) (pool *HttpHostPool, err error) {
    rand.Seed(time.Now().Unix())

    config := consulapi.DefaultConfig()
    config.Address = address
    client, err := consulapi.NewClient(config)
    if err != nil {
        log.Fatal("LearnHTTP: create consul client error", err.Error())
        return
    }

    pool = &HttpHostPool{
        address:  address,
        serviceName: serviceName,
        client: client,
        t: time.NewTicker(time.Second*3),
    }
    list,_,err := pool.resolver()
    if err != nil {
        return
    }

    pool.Lock()
    pool.addrs = list
    pool.length = len(list)
    pool.Unlock()

    go pool.watch()

    return
}
