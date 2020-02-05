package log

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "os"
    "path/filepath"
    "strings"
)

var log *zap.Logger

type Config struct {
    Development bool
    ServiceName string
    FilePath    string
    MaxSize     int
    MaxBackups  int
    MaxAge int
}

func Default() *zap.Logger {
    if log == nil {
        log = NewLogger(Config{})
    }
    return log
}

func NewLogger(cfg Config) *zap.Logger {
    hook := lumberjack.Logger{
        Filename:   "", // 日志文件路径
        MaxSize:    128,      // 每个日志文件保存的最大尺寸 单位：M
        MaxBackups: 30,       // 日志文件最多保存多少个备份
        MaxAge:     7,        // 文件最多保存多少天
        Compress:   true,     // 是否压缩
    }
    if cfg.FilePath != "" {
        hook.Filename = cfg.FilePath
    }else {
        hook.Filename = getFilePath()
    }
    if cfg.MaxSize != 0 {
        hook.MaxSize = cfg.MaxSize
    }
    if cfg.MaxBackups != 0 {
        hook.MaxBackups = cfg.MaxBackups
    }
    if cfg.MaxAge != 0 {
        hook.MaxAge = cfg.MaxAge
    }
    logLevel := zap.NewAtomicLevel()
    config := zap.NewProductionEncoderConfig()
    config.EncodeTime = zapcore.ISO8601TimeEncoder

    writers := []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
    //开启控制台输出
    if cfg.Development {
        writers = append(writers, zapcore.AddSync(os.Stdout))
    }
    core := zapcore.NewCore(
        zapcore.NewJSONEncoder(config),
        zapcore.NewMultiWriteSyncer(writers...),
        logLevel,
    )
    options := []zap.Option{}
    if cfg.ServiceName != "" {
        options = append(options, zap.Fields(zap.String("serviceName", "serviceName")),)
    }
    if cfg.Development {
        options = append(options, zap.AddCaller(), zap.Development())
    }
    // 构造日志
    logger := zap.New(core, options...)
    return logger
}

func getCurrentDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
        panic(err)
    }
    return strings.Replace(dir, "\\", "/", -1)
}

func getFilePath() string {
    logfile := getCurrentDirectory() + "/" + getAppname() + ".log"
    return logfile
}

func getAppname() string {
    full := os.Args[0]
    full = strings.Replace(full, "\\", "/", -1)
    splits := strings.Split(full, "/")
    if len(splits) >= 1 {
        name := splits[len(splits)-1]
        name = strings.TrimSuffix(name, ".exe")
        return name
    }
    return ""
}
