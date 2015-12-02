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

var writer *FileWriter

type Logger struct {
    log  *logging.Logger
}

func New(name string) *Logger {
    logger := &Logger{}
    logger.log = logging.MustGetLogger(name)
    
    return logger
}

func Close() {
    if writer != nil {
        writer.Close()
    }
    writer = nil
}

// Default opts
func ConsoleOpts() *LoggerOpts {
    return &LoggerOpts{}
}
func DebugOpts() *LoggerOpts {
    return &LoggerOpts{ StdoutEnable:true, StdoutLevel:DEBUG }
}
func ErrOpts() *LoggerOpts {
    return &LoggerOpts{ StderrEnable:true, StderrLevel:ERROR }
}
func FileOpts(fileNanme string) *LoggerOpts { 
    return &LoggerOpts{ FileEnable:true, FileLevel:WARNING, FilePath:fileName, FileRotateEnable:true, 
		FileSizeLimit:1024*1024, FileRotateLimit:3 }
}

func SetOpts(c *LoggerOpts) {
    Close()
    
    //Backend is an interface
    backends := make([]logging.Backend, 0, 3)
    

// NewStringFormatter returns a new Formatter which outputs the log record as a
// string based on the 'verbs' specified in the format string.
//
// The verbs:
//
// General:
//     %{id}        Sequence number for log message (uint64).
//     %{pid}       Process id (int)
//     %{time}      Time when log occurred (time.Time)
//     %{level}     Log level (Level)
//     %{module}    Module (string)
//     %{program}   Basename of os.Args[0] (string)
//     %{message}   Message (string)
//     %{longfile}  Full file name and line number: /a/b/c/d.go:23
//     %{shortfile} Final file name element and line number: d.go:23
//     %{color}     ANSI color based on log level
//
// For normal types, the output can be customized by using the 'verbs' defined
// in the fmt package, eg. '%{id:04d}' to make the id output be '%04d' as the
// format string.
//
// For time.Time, use the same layout as time.Format to change the time format
// when output, eg "2006-01-02T15:04:05.999Z-07:00".
//
// For the 'color' verb, the output can be adjusted to either use bold colors,
// i.e., '%{color:bold}' or to reset the ANSI attributes, i.e.,
// '%{color:reset}' Note that if you use the color verb explicitly, be sure to
// reset it or else the color state will persist past your log message.  e.g.,
// "%{color:bold}%{time:15:04:05} %{level:-8s}%{color:reset} %{message}" will
// just colorize the time and level, leaving the message uncolored.
//
// There's also a couple of experimental 'verbs'. These are exposed to get
// feedback and needs a bit of tinkering. Hence, they might change in the
// future.
//
// Experimental:
//     %{longpkg}   Full package path, eg. github.com/go-logging
//     %{shortpkg}  Base package path, eg. go-logging
//     %{longfunc}  Full function name, eg. littleEndian.PutUint32
//     %{shortfunc} Base function name, eg. PutUint32
    
    if c.StdoutEnable {
        stdoutformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s}%{color:reset} %{message}",
        )
        stdout := logging.NewLogBackend(os.Stdout, "", 0)
        stdoutformatter := logging.NewBackendFormatter(stdout, stdoutformat)
        stdoutlevel := logging.AddModuleLevel(stdoutformatter)
        stdoutlevel.SetLevel(logging.Level(c.StdoutLevel), "")
        backends = append(backends, stdoutlevel)
    }
    if c.StderrEnable {
        stderrformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s}%{color:reset} %{message}",
        )
        stderr := logging.NewLogBackend(os.Stderr, "", 0)
        stderrformatter := logging.NewBackendFormatter(stderr, stderrformat)
        stderrlevel := logging.AddModuleLevel(stderrformatter)
        stderrlevel.SetLevel(logging.Level(c.StderrLevel), "")
        backends = append(backends, stderrlevel)
    }
    if c.FileEnable {
        writer = NewWriter(c.FilePath, c.FileRotateEnable, c.FileSizeLimit, c.FileRotateLimit)
        if writer == nil { 
            fmt.Println("File writer creation failed.")
        } else {
            fileformat := logging.MustStringFormatter( 
                "{time:15:04:05.000} %{module} ▶ %{level:.4s} %{message}",
            )
            file := logging.NewLogBackend(writer, "", 0)
            
            fileformatter := logging.NewBackendFormatter(file, fileformat)
            filelevel := logging.AddModuleLevel(fileformatter)
            filelevel.SetLevel(logging.Level(c.FileLevel), "")
            backends = append(backends, filelevel)
        }
    }
    
    if len(backends) == 0 {
        //log disabled, setup default log to console
        consoleformat := logging.MustStringFormatter( 
            "%{color}%{time:15:04:05.000} %{module} ▶ %{level:.4s}%{color:reset} %{message}",
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