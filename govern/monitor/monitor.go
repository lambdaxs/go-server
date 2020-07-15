package monitor

import (
    "github.com/labstack/echo"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartMonitorServer(c echo.Context) error {
    promhttp.Handler().ServeHTTP(c.Response(), c.Request())
    return nil
}