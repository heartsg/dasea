package meta

import (
	"os"
    "errors"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
)

var Engine *xorm.Engine

var (
    ErrNotFound = errors.New("Item not found in database.")
    ErrInvalidDataPoints = errors.New("Invalid data points.")
)

// We cannot initialize xorm.Engine in init() function because Opts are
// initialized in init function, and we cannot guarentee that Opts must be
// initialized before any other init functions. So we intialize engine in main.
// We must call meta.InitEngine in main before using it. 
func InitEngine(t string, hosts []string) {
	var err error
    //currently we only use one host, we will switch to using
    //2 alternatively in the future
	Engine, err = xorm.NewEngine(t, hosts[0])
	logger := xorm.NewSimpleLogger(os.Stdout)
	logger.SetLevel(core.LOG_OFF)
	Engine.SetLogger(logger)
	if err != nil {
		panic(err)
	}
}
