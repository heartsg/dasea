package storage

import (
    "github.com/heartsg/dasea/config"
)

var opts Opts

func init() {
    c := config.New()
    c.Load(&opts)
}