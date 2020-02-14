package local

import (
    "net"
    "os"
    "sync"
)

var (
    hostName string
    ip string
    hostOnce sync.Once
    ipOnce sync.Once
)

//本机名称
func HostName() string {
    hostOnce.Do(func() {
        hostName,_ = os.Hostname()
    })
    return hostName
}

//本地内网IP
func LocalIP() string {
    ipOnce.Do(func() {
        ip = getInnerIPV4()
    })
    return ip
}

func getInnerIPV4() string {
    addrs, _ := net.InterfaceAddrs()
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}
