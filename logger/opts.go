package logger

//Level values:
//	 0: CRITICAL
//	 1: ERROR
//	 2: WARNING
//	 3: NOTICE
//	 4: INFO
// 	 5: DEBUG

type Opts struct {
    StdoutEnable bool `default:"true"`
    StdoutLevel int `default:"4"`
    StderrEnable bool `default:"false"`
    StderrLevel int `default:"1"`
    FileEnable bool `default:"true"`
    FileLevel int `default:"2"`
    FilePath string `default:"tmp/log"`
    FileSizeLimit int64
    FileRotateEnable bool
    FileRotateLimit int
}