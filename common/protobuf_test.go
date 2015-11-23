package common

import (
	"testing"
	"time"
	"dasea/database/testdata"
	"github.com/golang/protobuf/proto"
)

//Test TestData_1
func TestProtobuf1(t *testing.T) {
	now := time.Now()
	test := &testdata.TestData_1 {
		[]*testdata.TestData_1_Data_1{
			&testdata.TestData_1_Data_1{
				Radar: 100,
				Dummy: 0,
				Time: &testdata.Timestamp{Seconds:now.Unix(), Nanos:int32(now.Nanosecond())},
			},
			&testdata.TestData_1_Data_1{
				Radar: 200,
				Dummy: 1,
				Time: &testdata.Timestamp{Seconds:now.Unix(), Nanos:int32(now.Nanosecond())},			
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_1{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	radar, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if radar != 100 {
		t.Fatal("radar should be 100")
	}
	
	//dummy is 0, and it will not be encoded, escape it
	
	err = testData.DecodeCheckKey(WireLengthDelimited, 3)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != now.Unix() || ti.Nanosecond() != now.Nanosecond() {
		t.Fatal("time not match")
	}
	
	
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	rawBytes, err = protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData = NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	radar, err = testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if radar != 200 {
		t.Fatal("radar should be 200")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	dummy, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if dummy != 1 {
		t.Fatal("dummy should be 1")
	}
	
	err = testData.DecodeCheckKey(WireLengthDelimited, 3)
	if err != nil {
		t.Fatal(err)
	}
	ti, err = testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != now.Unix() || ti.Nanosecond() != now.Nanosecond() {
		t.Fatal("time not match")
	}	
}

//Test TestData2
func TestProtobuf2(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 100.0,
				F: 100.0,
				Time: &testdata.Timestamp{Seconds:100, Nanos:100},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireFixed32, 5)
	if err != nil {
		t.Fatal(err)
	}
	
	e, err := testData.DecodeFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if e != 100.0 {
		t.Fatal("e should be 100.0")
	}
	
	err = testData.DecodeCheckKey(WireFixed64, 6)
	if err != nil {
		t.Fatal(err)
	}
	
	f, err := testData.DecodeFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 100.0 {
		t.Fatal("f should be 100.0")
	}

	err = testData.DecodeCheckKey(WireLengthDelimited, 7)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != 100 || ti.Nanosecond() != 100 {
		t.Fatal("time not match")
	}
}

//Test when float field is 0
func TestProtobufFloat0(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 0.0,
				F: 0.0,
				Time: &testdata.Timestamp{Seconds:100, Nanos:100},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	//tag 5 & 6 escaped because they are 0.0
	
	err = testData.DecodeCheckKey(WireLengthDelimited, 7)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != 100 || ti.Nanosecond() != 100 {
		t.Fatal("time not match")
	}
}

//Test when timestamp field is 0
func TestProtobufTimestamp0(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 100.0,
				F: 100.0,
				Time: &testdata.Timestamp{Seconds:0, Nanos:0},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireFixed32, 5)
	if err != nil {
		t.Fatal(err)
	}
	
	e, err := testData.DecodeFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if e != 100.0 {
		t.Fatal("e should be 100.0")
	}
	
	err = testData.DecodeCheckKey(WireFixed64, 6)
	if err != nil {
		t.Fatal(err)
	}
	
	f, err := testData.DecodeFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 100.0 {
		t.Fatal("f should be 100.0")
	}

	err = testData.DecodeCheckKey(WireLengthDelimited, 7)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != 0 || ti.Nanosecond() != 0 {
		t.Fatal("time not match")
	}
}

//Test when timestamp field is 0
func TestProtobufTimestampNanos0(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 100.0,
				F: 100.0,
				Time: &testdata.Timestamp{Seconds:100, Nanos:0},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireFixed32, 5)
	if err != nil {
		t.Fatal(err)
	}
	
	e, err := testData.DecodeFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if e != 100.0 {
		t.Fatal("e should be 100.0")
	}
	
	err = testData.DecodeCheckKey(WireFixed64, 6)
	if err != nil {
		t.Fatal(err)
	}
	
	f, err := testData.DecodeFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 100.0 {
		t.Fatal("f should be 100.0")
	}

	err = testData.DecodeCheckKey(WireLengthDelimited, 7)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != 100 || ti.Nanosecond() != 0 {
		t.Fatal("time not match")
	}
}


//Test when timestamp field is 0
func TestProtobufTimestampSeconds0(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 100.0,
				F: 100.0,
				Time: &testdata.Timestamp{Seconds:0, Nanos:100},
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireFixed32, 5)
	if err != nil {
		t.Fatal(err)
	}
	
	e, err := testData.DecodeFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if e != 100.0 {
		t.Fatal("e should be 100.0")
	}
	
	err = testData.DecodeCheckKey(WireFixed64, 6)
	if err != nil {
		t.Fatal(err)
	}
	
	f, err := testData.DecodeFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 100.0 {
		t.Fatal("f should be 100.0")
	}

	err = testData.DecodeCheckKey(WireLengthDelimited, 7)
	if err != nil {
		t.Fatal(err)
	}
	ti, err := testData.DecodeTimestamp()
	if err != nil {
		t.Fatal(err)
	}
	if ti.Unix() != 0 || ti.Nanosecond() != 100 {
		t.Fatal("time not match")
	}
}


//Test when timestamp field is 0
func TestProtobufTimestampNil(t *testing.T) {
	test := &testdata.TestData_2 {
		[]*testdata.TestData_2_Data_2{
			&testdata.TestData_2_Data_2{
				A: 100,
				B: 100,
				C: 100,
				D: 100,
				E: 100.0,
				F: 100.0,
				Time: nil,
			},
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)
	err = testData.DecodeCheckKey(WireVarint, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	a, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 {
		t.Fatal("a should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireVarint, 2)
	if err != nil {
		t.Fatal(err)
	}
	
	b, err := testData.DecodeVarint()
	if err != nil {
		t.Fatal(err)
	}
	if b != 100 {
		t.Fatal("b should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 3)
	if err != nil {
		t.Fatal(err)
	}
	
	c, err := testData.DecodeZigzag32()
	if err != nil {
		t.Fatal(err)
	}
	if c != 100 {
		t.Fatal("c should be 100")
	}
	
	err = testData.DecodeCheckKey(WireVarint, 4)
	if err != nil {
		t.Fatal(err)
	}
	
	d, err := testData.DecodeZigzag64()
	if err != nil {
		t.Fatal(err)
	}
	if d != 100 {
		t.Fatal("d should be 100")
	}
	
	
	err = testData.DecodeCheckKey(WireFixed32, 5)
	if err != nil {
		t.Fatal(err)
	}
	
	e, err := testData.DecodeFloat32()
	if err != nil {
		t.Fatal(err)
	}
	if e != 100.0 {
		t.Fatal("e should be 100.0")
	}
	
	err = testData.DecodeCheckKey(WireFixed64, 6)
	if err != nil {
		t.Fatal(err)
	}
	
	f, err := testData.DecodeFloat64()
	if err != nil {
		t.Fatal(err)
	}
	if f != 100.0 {
		t.Fatal("f should be 100.0")
	}

	if !testData.DecodeComplete() {
		t.Fatal("timestamp is nil, there should not be any more data")
	}
}


//Test when all fields are 0
func TestProtobufAll0(t *testing.T) {
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
		},
	}
	data, err := proto.Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	
	testUnmarshal := &testdata.TestData_2{}
	err = proto.Unmarshal(data, testUnmarshal)
	if err != nil {
		t.Fatal(err)
	}

	protoBuffer := NewProtoBuffer(data)
	err = protoBuffer.DecodeCheckKey(WireLengthDelimited, 1)
	if err != nil {
		t.Fatal(err)
	}
	
	rawBytes, err := protoBuffer.DecodeRawBytes(false)
	if err != nil {
		t.Fatal(err)
	}

	testData := NewProtoBuffer(rawBytes)

	if !testData.DecodeComplete() {
		t.Fatal("all fields are 0, there should not be any more data")
	} 
}