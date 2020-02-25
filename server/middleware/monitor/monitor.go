package monitor

import (
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServerMetric *prometheus.HistogramVec
	ErrorMetric  *prometheus.CounterVec
	SystemMetric *prometheus.HistogramVec
)

func Init(serviceName string){
	serviceName = strings.ReplaceAll(serviceName,"-","_")
	serviceName = strings.ReplaceAll(serviceName,".","_")
	serviceName = strings.ReplaceAll(serviceName," ","_")

	ServerMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "go_server",
		Subsystem: serviceName,
		Name:      "server_handle_duration_ms",
		Help:      "业务请求吞吐量tps p99",
	}, []string{"type", "path", "code"})

	ErrorMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "go_server",
		Subsystem: serviceName,
		Name:      "server_handle_error_total",
		Help:      "业务错误数",
	}, []string{"type", "path", "code"})

	SystemMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "go_server",
		Subsystem:   serviceName,
		Name:        "system_info",
		Help:        "系统数值监控",
	},[]string{"type","name"})

	go startSystemMonitor()
	//注册监控器
	prometheus.MustRegister(ServerMetric, ErrorMetric, SystemMetric)
}

//system monitor
func startSystemMonitor(){
	ticker := time.NewTicker(time.Second*10)
	for range ticker.C {
		info := GetSystemInfo()
		SystemMetric.WithLabelValues("cpu","user").Observe(info.CPU.User)
		SystemMetric.WithLabelValues("cpu","system").Observe(info.CPU.System)
		SystemMetric.WithLabelValues("cpu","idle").Observe(info.CPU.Idle)
		SystemMetric.WithLabelValues("cpu","iowait").Observe(info.CPU.IOWait)
		SystemMetric.WithLabelValues("cpu","nice").Observe(info.CPU.Nice)
		SystemMetric.WithLabelValues("cpu","steal").Observe(info.CPU.Steal)

		SystemMetric.WithLabelValues("mem","usedpercent").Observe(info.Mem.UsedPercent)
		SystemMetric.WithLabelValues("mem","total").Observe(float64(info.Mem.Total))
		SystemMetric.WithLabelValues("mem","used").Observe(float64(info.Mem.Used))
		SystemMetric.WithLabelValues("mem","free").Observe(float64(info.Mem.Free))

		SystemMetric.WithLabelValues("net","input").Observe(float64(info.NetIO.BytesRecv))
		SystemMetric.WithLabelValues("net","output").Observe(float64(info.NetIO.BytesSent))

		SystemMetric.WithLabelValues("process","fd").Observe(float64(info.Process.Fd))

		SystemMetric.WithLabelValues("load","load1").Observe(info.Load.Load1)
		SystemMetric.WithLabelValues("load","load5").Observe(info.Load.Load5)
		SystemMetric.WithLabelValues("load","load15").Observe(info.Load.Load15)
	}
}


func HTTPMonitor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		path := c.Request().URL.String()
		code := strconv.Itoa(c.Response().Status)
		if err := next(c); err != nil {
			return err
		}
		//记录tps p99
		ServerMetric.WithLabelValues("http", path, code).Observe(float64(time.Since(start).Milliseconds()))
		return nil
	}
}
