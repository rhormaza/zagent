package util

import (
    "log4go"
)

//type Logger log4go.Logger

func SetupLogger(filename string, level log4go.Level) (*log4go.Logger) {
    // Get a new logger instance
    log := make(log4go.Logger)

    flw := log4go.NewFileLogWriter(filename, false)
    flw.SetFormat("[%D %T] [%L] (%S) %M")
    flw.SetRotate(false)
    flw.SetRotateSize(0)
    flw.SetRotateLines(0)
    flw.SetRotateDaily(false)
    log.AddFilter("file", level, flw)

    return &log
}
