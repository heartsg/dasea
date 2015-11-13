package logger

//Efficiency is not major concern for now
//Use go-logging package
//May need to reimplement

import (
    "os"
    "fmt"
    "github.com/op/go-logging"
)

const (
	CRITICAL int = int(logging.CRITICAL)
	ERROR = int(logging.ERROR)
	WARNING = int(logging.WARNING)
	NOTICE = int(logging.NOTICE)
	INFO = int(logging.INFO)
	DEBUG = int(logging.DEBUG)
)

type LoggerConfig struct {
    StdoutEnable bool
    StdoutLevel int
    StderrEnable bool
    StderrLevel int
    FileEnable bool
    FileLevel int
    FilePath string
    FileSizeLimit int64
    FileRotateEnable bool
    FileRotateLimit int
}

type Logger struct {
    log  *logging.Logger
    writer *FileWriter
}

func NewLogger(config *LoggerConfig) *Logger {
    logger := &Logger{}
    logger.log = logging.MustGetLogger("dasea")
    
    logger.SetConfig(config)

    return logger
}

func (l *Logger) Close() {
    if l.writer != nil {
        l.writer.Close()
    }
}
func (l *Logger) SetConfig(c *LoggerConfig) {
    l.Close()
    
    //Backend is an interface
    backends := make([]logging.Backend, 0, 3)
    
    if c.StdoutEnable {
        stdoutformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
        )
        stdout := logging.NewLogBackend(os.Stdout, "", 0)
        stdoutformatter := logging.NewBackendFormatter(stdout, stdoutformat)
        stdoutlevel := logging.AddModuleLevel(stdoutformatter)
        stdoutlevel.SetLevel(logging.Level(c.StdoutLevel), "")
        backends = append(backends, stdoutlevel)
    }
    if c.StderrEnable {
        stderrformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
        )
        stderr := logging.NewLogBackend(os.Stderr, "", 0)
        stderrformatter := logging.NewBackendFormatter(stderr, stderrformat)
        stderrlevel := logging.AddModuleLevel(stderrformatter)
        stderrlevel.SetLevel(logging.Level(c.StderrLevel), "")
        backends = append(backends, stderrlevel)
    }
    if c.FileEnable {
        writer := NewWriter(c.FilePath, c.FileRotateEnable, c.FileSizeLimit, c.FileRotateLimit)
        if writer == nil { 
            fmt.Println("File writer creation failed.")
        } else {
            fileformat := logging.MustStringFormatter( 
                "{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}",
            )
            file := logging.NewLogBackend(writer, "", 0)
            
            fileformatter := logging.NewBackendFormatter(file, fileformat)
            filelevel := logging.AddModuleLevel(fileformatter)
            filelevel.SetLevel(logging.Level(c.FileLevel), "")
            backends = append(backends, filelevel)
            
            l.writer = writer
        }
    }
    
    if len(backends) == 0 {
        //log disabled, setup default log to console
        consoleformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
        )
        console := logging.NewLogBackend(os.Stdout, "", 0)
        consoleformatter := logging.NewBackendFormatter(console, consoleformat)
        consolelevel := logging.AddModuleLevel(consoleformatter)
        consolelevel.SetLevel(logging.Level(WARNING), "")
        backends = append(backends, consolelevel)
    }
    logging.SetBackend(backends...)
}


// Fatal is equivalent to l.Critical(fmt.Sprint()) followed by a call to os.Exit(1).
func (l *Logger) Fatal(args ...interface{}) {
	l.log.Fatal(args...)
}

// Fatalf is equivalent to l.Critical followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
    
}

// Panic is equivalent to l.Critical(fmt.Sprint()) followed by a call to panic().
func (l *Logger) Panic(args ...interface{}) {
	l.log.Panic(args...)
}

// Panicf is equivalent to l.Critical followed by a call to panic().
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args...)
}

// Critical logs a message using CRITICAL as log level.
func (l *Logger) Critical(format string, args ...interface{}) {
    l.log.Critical(format, args...)
}

// Error logs a message using ERROR as log level.
func (l *Logger) Error(format string, args ...interface{}) {
    l.log.Error(format, args...)
}

// Warning logs a message using WARNING as log level.
func (l *Logger) Warning(format string, args ...interface{}) {
    l.log.Warning(format, args...)
}

// Notice logs a message using NOTICE as log level.
func (l *Logger) Notice(format string, args ...interface{}) {
    l.log.Notice(format, args...)
}

// Info logs a message using INFO as log level.
func (l *Logger) Info(format string, args ...interface{}) {
    l.log.Info(format, args...)
}

// Debug logs a message using DEBUG as log level.
func (l *Logger) Debug(format string, args ...interface{}) {
    l.log.Debug(format, args...)
}