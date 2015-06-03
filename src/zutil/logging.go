package zutil

import (
    "log4go"
)

// Logging stuff
//type Logger log4go.Logger
var log log4go.Logger = nil

var logLevelMap =  map[string]log4go.Level{
    "FINEST"  : log4go.FINEST,
    "FINE"    : log4go.FINE,
    "DEBUG"   : log4go.DEBUG,
    "TRACE"   : log4go.TRACE,
    "INFO"    : log4go.INFO,
    "WARNING" : log4go.WARNING,
    "ERROR"   : log4go.ERROR,
    "CRITICAL": log4go.CRITICAL,
}

func SetupLogger(filename string, level string) (*log4go.Logger) {
    // Get a new logger instance
    if log == nil {
        log = make(log4go.Logger)
        setupFileLogger(&log, filename, logLevelMap[level])
        setupConsoleLogger(&log, logLevelMap[level])
        //setupFileLogger(&log, filename, log4go.INFO)
        //setupConsoleLogger(&log, log4go.INFO)
        log.Info("======== Initializing application ========")
        log.Info("log file at:%s, logging level: %s", filename, level)
    }
    return &log

}

func setupFileLogger(log *log4go.Logger, filename string, level log4go.Level) {

    flw := log4go.NewFileLogWriter(filename, false)
    flw.SetFormat("[%D %T] [%L] (%S) %M")
    flw.SetRotate(false)
    flw.SetRotateSize(0)
    flw.SetRotateLines(0)
    flw.SetRotateDaily(false)
    // Logging to file
    log.AddFilter("file", level, flw)
}

func setupConsoleLogger(log *log4go.Logger, level log4go.Level) {

    flw1 := log4go.NewConsoleLogWriter()
    flw1.SetFormat("[%D %T] [%L] (%S) %M")
    // Logging to console
    log.AddFilter("stdout", level, flw1)
}

func GetLogger() (*log4go.Logger) {
    return &log 
}
