package database

import (
	"testing"
	"time"
	"github.com/golang/protobuf/proto"
	"dasea/database/testdata"
)

func init() {
	InitData()
}
func TestData1(t *testing.T) {
	dsa, err := CreateDataStreamAttribute("radar data", int16(3), []string{"radar", "dummy", "time"},
		[]string{"int8", "int8", "timestamp"}, []int64{10350001, 10001, 10001} )
	if err != nil {
		t.Fatal(err)
	}
	if dsa.Id != 1 {
		t.Fatal("Data stream attribute Id should be 1")
	}
	
	ds, err := CreateDataStream(1, dsa.Id)
	if err != nil {
		t.Fatal(err)
	}
	if ds.Id != 1 {
		t.Fatal("Data stream Id should be 1")
	}
	
	err = CreateData(ds.Id)
	if err != nil {
		t.Fatal(err)
	}
	
	//data is type of []byte
	time := time.Now()
	test := &testdata.TestData_1 {
		[]*testdata.TestData_1_Data_1{
			&testdata.TestData_1_Data_1{
				Radar: 100,
				Dummy: 0,
				Time: &testdata.Timestamp{Seconds:time.Unix(), Nanos:int32(time.Nanosecond())},
			},
			&testdata.TestData_1_Data_1{
				Radar: 200,
				Dummy: 1,
				Time: &testdata.Timestamp{Seconds:time.Unix(), Nanos:int32(time.Nanosecond())},			
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}

	err = PutDataPointsFromProtobuf(ds.Id, data)
	if err != nil {
		t.Fatal(err)
	}
}


func TestData2(t *testing.T) {
	dsa, err := CreateDataStreamAttribute("test", int16(7), []string{"a", "b", "c", "d", "e", "f", "time"},
		[]string{"int32", "int64", "sint32", "sint64", "float", "double", "timestamp"}, 
		[]int64{10001, 10001, 10001, 10001, 10001, 10001, 10001} )
	if err != nil {
		t.Fatal(err)
	}
	if dsa.Id != 2 {
		t.Fatal("Data stream attribute Id should be 2")
	}
	
	ds, err := CreateDataStream(1, dsa.Id)
	if err != nil {
		t.Fatal(err)
	}
	if ds.Id != 2 {
		t.Fatal("Data stream Id should be 2")
	}
	
	err = CreateData(ds.Id)
	if err != nil {
		t.Fatal(err)
	}
	
	//data is type of []byte
	time := time.Now()
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 0,
				B: 0,
				C: 0,
				D: 0,
				E: 0,
				F: 0,
				Time: nil,
			},
			&testdata.TestData_2_Data_2{
				A: 0,
				B: 0,
				C: 100,
				D: 0,
				E: 0,
				F: 0,
				Time: nil,
			},
			&testdata.TestData_2_Data_2{
				A: 1,
				B: 0,
				C: 100,
				D: 0,
				E: 0,
				F: 0,
				Time: nil,
			},
			&testdata.TestData_2_Data_2{
				A: 0,
				B: 0,
				C: 100,
				D: 0,
				E: 0,
				F: 0,
				Time: nil,
			},
			&testdata.TestData_2_Data_2{
				A: 0,
				B: 0,
				C: 100,
				D: -1000,
				E: 0,
				F: 0,
				Time: nil,
			},
			&testdata.TestData_2_Data_2{
				A: -1,
				B: -2,
				C: -3,
				D: -4,
				E: -13.2,
				F: 132.0,
				Time: &testdata.Timestamp{Seconds:time.Unix(), Nanos:int32(time.Nanosecond())},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}

	err = PutDataPointsFromProtobuf(ds.Id, data)
	if err != nil {
		t.Fatal(err)
	}
}