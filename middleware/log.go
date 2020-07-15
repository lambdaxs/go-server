package middleware

import (
    grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
    "github.com/labstack/echo"
    "github.com/lambdaxs/go-server/govern/log"
    "go.uber.org/zap"
    "google.golang.org/grpc"
    "time"
)

func HttpServerLogger() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            req := c.Request()
            res := c.Response()
            start := time.Now()
            if err := next(c); err != nil {
                c.Error(err)
            }
            cost := time.Since(start).Milliseconds() //ms
            log.Default().Info("http access", zap.Int("code", res.Status), zap.String("path", req.URL.Path), zap.Int64("cost", cost))
            return nil
        }
    }
}

func GRPCServerLogger() grpc.UnaryServerInterceptor {
    return grpc_zap.UnaryServerInterceptor(log.Default().Logger)
}

