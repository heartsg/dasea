package storage

import (
    "testing"
    "github.com/heartsg/dasea/config"
)


func TestOpts(t *testing.T) {
    var opts Opts
    c := config.NewWithPath("testdata/config.toml", "auto")
    c.Load(&opts)
    
    if opts.MetaDB.Hosts[0] != "127.0.0.1" || opts.DataDB.Hosts[0] != "127.0.0.1"  {
        t.Error("Database Host error")
    }
    if opts.KeystoneMiddleware.Client.AuthUrl != "127.0.0.1:35357" {
        t.Error("Keystone AuthUrl error")
    }
}