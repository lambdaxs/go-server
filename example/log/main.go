package main

import (
    "github.com/lambdaxs/go-server/log"
    "go.uber.org/zap"
)

func main(){
    log.Default().Info("this is new msg111", zap.String("name", "xiaos"))
}
