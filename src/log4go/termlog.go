// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"io"
	"os"
	"time"
)

var stdout io.Writer = os.Stdout
var format string = FORMAT_SHORT 

// This is the standard writer that prints to standard output.
type ConsoleLogWriter chan *LogRecord

// This creates a new ConsoleLogWriter
func NewConsoleLogWriter() ConsoleLogWriter {
	records := make(ConsoleLogWriter, LogBufferLength)
	go records.run(stdout)
	return records
}

func (w ConsoleLogWriter) run(out io.Writer) {
    // Using a new method to be able to set format
	for rec := range w {
        fmt.Fprint(out, FormatLogRecord(format, rec))
	}
}

// This is the ConsoleLogWriter's output method.  This will block if the output
// buffer is full.
func (w ConsoleLogWriter) LogWrite(rec *LogRecord) {
	w <- rec
}

// Close stops the logger from sending messages to standard output.  Attempts to
// send log messages to this logger after a Close have undefined behavior.
func (w ConsoleLogWriter) Close() {
	close(w)
	time.Sleep(50 * time.Millisecond) // Try to give console I/O time to complete
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w ConsoleLogWriter) SetFormat(logFormat string) {
	format = logFormat
}
