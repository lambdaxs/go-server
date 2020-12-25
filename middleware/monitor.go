package middleware

import (
    grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/govern/log"
    "github.com/lambdaxs/go-server/govern/monitor"
    "github.com/prometheus/client_golang/prometheus"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "strconv"
    "time"
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

    //注册监控器
    prometheus.MustRegister(ServerMetric, ErrorMetric, SystemMetric)
}

func InitSystemMonitor() {
    go startSystemMonitor()
}

//system monitor
func startSystemMonitor() {
    ticker := time.NewTicker(time.Second * 10)
    for range ticker.C {
        info := monitor.GetSystemInfo()
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

func HttpServerMonitor() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            start := time.Now()
            path := c.Request().URL.Path
            code := strconv.Itoa(c.Response().Status)

            if err := next(c); err != nil {
                //记录错误监控
                ErrorMetric.WithLabelValues("http", path, code).Inc()
                //记录错误日志
                log.Default().Error("http req error" , zap.String("code", code), zap.String("path", path), zap.String("detail", err.Error()))
                return err
            }

            //记录tps p99
            ServerMetric.WithLabelValues("http", path, code).Observe(float64(time.Since(start).Milliseconds()))
            return nil
        }
    }
}

func GRPCServerMonitor() grpc.UnaryServerInterceptor {
    return grpc_prometheus.UnaryServerInterceptor
}
