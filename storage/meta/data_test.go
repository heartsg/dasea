package meta

import (
    "testing"
)

// test DB
func TestData(t *testing.T) {
    InitEngine("mysql", []string{"dasea:dasea@tcp(127.0.0.1:3306)/dasea?charset=utf8"})
    err := CreateDataStreamAttributeTable()
    if err!= nil {
        t.Error(err)
    }
    
    err = CreateDataStreamTable()
    if err != nil {
        t.Error(err)
    }

    err = InsertDataStreamAttribute(&DataStreamAttribute{
        Description: "test 1",
        NumDataPoints: 2,
        DataPointNames: []string {"data1", "data2"},
        DataPointTypes: []string {"uint8", "uint8"},
        DataPointUnits: []int64 {10001, 10001},
        ProjectId: "test",
        DomainId: "test",
    })
    _, err = CreateDataStreamAttribute("test", "test", "test 2", 1, []string{"data"}, []string{"uint16"}, []int64{10001})
    if err != nil {
        t.Error(err)
    }
    
    err = InsertDataStream(&DataStream{
        DeviceId: 1,
        DataStreamAttributeId: 1,
    })
    if err != nil {
        t.Error(err)
    }
    _, err = CreateDataStream(1, 2)
    if err != nil {
        t.Error(err)
    }
    
    a1, err := GetDataStreamAttribute(1)
    if err != nil {
        t.Error(err)
    }
    a2, err := GetDataStreamAttribute(2)
    if err != nil {
        t.Error(err)
    }
    s1, err := GetDataStream(1)
    if err != nil {
        t.Error(err)
    }
    s2, err := GetDataStream(2)
    if err != nil {
        t.Error(err)
    }
   
    if a1.Id != 1 || a2.Id != 2 || s1.Id != 1 || s2.Id != 2 {
        t.Error("Something wrong with insert and get")
    }
    
    s3, err := GetDataStreamAttributeByDataStreamId(1)
    if err != nil {
        t.Error(err)
    }
    s4, err := GetDataStreamAttributeByDataStreamId(2)
    if err != nil {
        t.Error(err)
    }
    if s3.Id != 1 || s4.Id != 2 {
        t.Error("Something wrong with Get Attribute by stream id")
    }

    err = DeleteDataStreamAttribute(1)
    if err != nil {
        t.Error(err)
    }
    err = DeleteDataStreamAttribute(2)
    if err != nil {
        t.Error(err)
    }
    err = DeleteDataStream(1)
    if err != nil {
        t.Error(err)
    }
    err = DeleteDataStream(2)
    if err != nil {
        t.Error(err)
    }
    
    a3, err := GetDataStreamAttribute(1)
    if a3 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    a4, err := GetDataStreamAttribute(2)
    if a4 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    
    s5, err := GetDataStream(1)
    if s5 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    s6, err := GetDataStream(2)
    if s6 != nil || err != ErrNotFound {
        t.Error("Still get something after delete")
    }
    
    
    Engine.Id(1).Unscoped().Delete(s1)
    Engine.Id(2).Unscoped().Delete(s1)
    Engine.Id(1).Unscoped().Delete(a1)
    Engine.Id(2).Unscoped().Delete(a1)
    Engine.DropTables("data_stream_attribute")
    Engine.DropTables("data_stream")
}