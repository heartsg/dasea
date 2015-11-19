package logger

import (
	"testing"
	
	"github.com/heartsg/dasea/config"
)


func TestLoggerOpts(t *testing.T) {
	loader := config.New()
	loggerOpts := &LoggerOpts{}
	loader.Load(loggerOpts)
	if loggerOpts.StdoutEnable != true {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.StdoutLevel != INFO {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.StderrEnable != false {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.StderrLevel != ERROR {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FileEnable != true {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FileLevel != WARNING {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FilePath != "tmp/log" {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FileSizeLimit != 0 {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FileRotateEnable != false {
		t.Error("LoggerOpts load error")
	}
	if loggerOpts.FileRotateLimit != 0 {
		t.Error("LoggerOpts load error")
	}
}
