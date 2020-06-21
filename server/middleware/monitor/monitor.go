package monitor

import (
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ServerMetric *prometheus.HistogramVec
	ErrorMetric  *prometheus.CounterVec
	SystemMetric *prometheus.GaugeVec
)

func init() {
	ServerMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "go_server_handle_duration_ms",
		Help:      "业务请求吞吐量tps p99",
	}, []string{"type", "path", "code"})

	ErrorMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "go_server_handle_error_total",
		Help:      "业务错误数",
	}, []string{"type", "path", "code"})

	SystemMetric = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "",
		Subsystem: "",
		Name:      "go_system_info",
		Help:      "系统数值监控",
	}, []string{"type", "name"})

	go startSystemMonitor()
	//注册监控器
	prometheus.MustRegister(ServerMetric, ErrorMetric, SystemMetric)
}

//system monitor
func startSystemMonitor() {
	ticker := time.NewTicker(time.Second * 10)
	for range ticker.C {
		info := GetSystemInfo()
		SystemMetric.WithLabelValues("cpu", "user").Set(info.CPU.User)
		SystemMetric.WithLabelValues("cpu", "system").Set(info.CPU.System)
		SystemMetric.WithLabelValues("cpu", "idle").Set(info.CPU.Idle)
		SystemMetric.WithLabelValues("cpu", "iowait").Set(info.CPU.IOWait)
		SystemMetric.WithLabelValues("cpu", "nice").Set(info.CPU.Nice)
		SystemMetric.WithLabelValues("cpu", "steal").Set(info.CPU.Steal)

		SystemMetric.WithLabelValues("mem", "usedpercent").Set(info.Mem.UsedPercent)
		SystemMetric.WithLabelValues("mem", "total").Set(float64(info.Mem.Total))
		SystemMetric.WithLabelValues("mem", "used").Set(float64(info.Mem.Used))
		SystemMetric.WithLabelValues("mem", "free").Set(float64(info.Mem.Free))

		SystemMetric.WithLabelValues("net", "input").Set(float64(info.NetIO.BytesRecv))
		SystemMetric.WithLabelValues("net", "output").Set(float64(info.NetIO.BytesSent))

		SystemMetric.WithLabelValues("process", "fd").Set(float64(info.Process.Fd))

		SystemMetric.WithLabelValues("load", "load1").Set(info.Load.Load1)
		SystemMetric.WithLabelValues("load", "load5").Set(info.Load.Load5)
		SystemMetric.WithLabelValues("load", "load15").Set(info.Load.Load15)
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
