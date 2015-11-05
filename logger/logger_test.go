package logger

import (
	"testing"
	"os"
)

func TestConsoleLogger(t *testing.T) {
	c := &LoggerConfig{} //all false, defaulted to console
	l := NewLogger(c)
	l.Debug("debug")
	l.Info("info")
	l.Notice("notice")
	l.Warning("warning")
	l.Error("error")
	l.Critical("critical")
}

func TestStdoutLogger(t *testing.T) {
	c := &LoggerConfig{ StdoutEnable:true, StdoutLevel:DEBUG } 
	l := NewLogger(c)
	l.Debug("debug")
	l.Info("info")
	l.Notice("notice")
	l.Warning("warning")
	l.Error("error")
	l.Critical("critical")
}

func TestStderrLogger(t *testing.T) {
	c := &LoggerConfig{ StderrEnable:true, StderrLevel:ERROR } 
	l := NewLogger(c)
	l.Debug("debug")
	l.Info("info")
	l.Notice("notice")
	l.Warning("warning")
	l.Error("error")
	l.Critical("critical")
}

func TestFileLogger(t *testing.T) {
	c := &LoggerConfig{ FileEnable:true, FileLevel:WARNING, FilePath:"test.log", FileRotateEnable:true, 
		FileSizeLimit:1024*1024, FileRotateLimit:3 } 
	l := NewLogger(c)
	l.Debug("debug")
	l.Info("info")
	l.Notice("notice")
	l.Warning("warning")
	l.Error("error")
	l.Critical("critical")
	l.Close()
	os.Remove("test.log.0")	
}

func TestAllLogger(t *testing.T) {
	c := &LoggerConfig{ StdoutEnable:true, StdoutLevel:DEBUG,
		StderrEnable:true, StderrLevel:ERROR,
		FileEnable:true, FileLevel:WARNING, FilePath:"test_together.log", FileRotateEnable:true, 
		FileSizeLimit:1024*1024, FileRotateLimit:3 } 
	l := NewLogger(c)
	l.Debug("debug together")
	l.Info("info together")
	l.Notice("notice together")
	l.Warning("warning together")
	l.Error("error together")
	l.Critical("critical together")
	l.Close()
	os.Remove("test_together.log.0")
}