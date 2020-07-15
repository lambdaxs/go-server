package middleware

import (
    grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
    "github.com/labstack/echo"
    "google.golang.org/grpc"
)

// todo
func HttpServerTracer() echo.MiddlewareFunc {
    return nil
}

func GRPCServerTracer() grpc.UnaryServerInterceptor{
    return grpc_opentracing.UnaryServerInterceptor()
}