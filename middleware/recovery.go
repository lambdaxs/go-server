package middleware

import (
    grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
    "github.com/labstack/echo"
    echo_middleware "github.com/labstack/echo/middleware"
    "google.golang.org/grpc"
)

func HttpServerRecovery() echo.MiddlewareFunc {
    return echo_middleware.Recover()
}

func GRPCServerRecovery() grpc.UnaryServerInterceptor {
    return grpc_recovery.UnaryServerInterceptor()
}
