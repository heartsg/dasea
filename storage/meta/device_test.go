package meta

import (
    "testing"
)

// test DB
func TestDevice(t *testing.T) {
    InitEngine("mysql", []string{"dasea:dasea@tcp(127.0.0.1:3306)/dasea?charset=utf8"})
    err := CreateAggregationDeviceTable()
    if err!= nil {
        t.Error(err)
    }
    
    err = CreateDeviceTable()
    if err != nil {
        t.Error(err)
    }

    err = InsertAggregationDevice(&AggregationDevice{
        Id: "123",
        Description: "test device (aggregation)",
        ProjectId: "test",
        DomainId: "test",
    })

    if err != nil {
        t.Error(err)
    }
    
    err = InsertDevice(&Device{
        AggregationDeviceId: "123",
    })
    
    a1, err := GetAggregationDevice("123")
    if err != nil {
        t.Error(err)
    }
    
    d1, err := GetDevice(1)
    if err != nil {
        t.Error(err)
    }
    
    
    if a1.ProjectId != "test" || d1.AggregationDeviceId != "123" {
        t.Error("Something wrong with insert and get")
    }

    err = DeleteDevice(1)
    if err != nil {
        t.Error(err)
    }
    err = DeleteAggregationDevice("123")
    if err != nil {
        t.Error(err)
    }
    
    a2, err := GetAggregationDevice("123")
    if a2 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    d2, err := GetDevice(1)
    if d2 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    
    Engine.Id(1).Unscoped().Delete(d1)
    Engine.Id("123").Unscoped().Delete(a1)
    Engine.DropTables("device")
    Engine.DropTables("aggregation_device")
}