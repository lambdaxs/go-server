package monitor

import "github.com/prometheus/client_golang/prometheus"

var (
    ServerMetric *prometheus.HistogramVec
    ErrorMetric *prometheus.CounterVec
)

func InitServerMonitor(namespace, subsystem string) {
    ServerMetric = prometheus.NewHistogramVec(prometheus.HistogramOpts{
        Namespace:   namespace,
        Subsystem:   subsystem,
        Name:        "server_handle_duration_ms",
        Help:        "业务请求吞吐量tps p99",
    }, []string{"path","code"})

    ErrorMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
        Namespace:   namespace,
        Subsystem:   subsystem,
        Name:        "server_handle_error_total",
        Help:        "业务错误数",
    }, []string{"path","code"})
    //注册监控器
    prometheus.MustRegister(ServerMetric, ErrorMetric)
}