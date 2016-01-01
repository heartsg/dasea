package meta

// Data meta, including
//  - DataStreamAttribute
//  - DataStream
import (
    "time"
)

type DataStreamAttribute struct {
	Id int64
	Description string `xorm:"varchar(255) notnull unique"`
	NumDataPoints int16
	//User defined names for each data point (column)
	DataPointNames []string
	//support only an array of int, int8, int16, int32, uint, uint8, uint16, uint32, 
	// int64, uint64, float32, float64, bool
	// datetime
	DataPointTypes []string
	//an array of units id (refer to the unit table)
	DataPointUnits []int64
    
    ProjectId string `xorm:"index"` //keystone project id
    DomainId string `xorm:"index"` //keystone domain id
	CreatedAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}


func CreateDataStreamAttributeTable() error {
    a := &DataStreamAttribute{}
	_ = Engine.DropTables(a)
	err := Engine.CreateTables(a)
    return err
}
func GetDataStreamAttribute(id int64) (*DataStreamAttribute, error) {
    a := &DataStreamAttribute{}
    has, err := Engine.Id(id).Get(a)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, ErrNotFound
    }
    return a, nil
}
func InsertDataStreamAttribute(a *DataStreamAttribute) error {
    _, err := Engine.Insert(a)
    return err
}
func CreateDataStreamAttribute(projectId string, domainId string, desc string, numDataPoints int16, dataPointNames []string, 
        dataPointTypes []string, dataPointUnits []int64) (*DataStreamAttribute, error) {
	if numDataPoints <= 0 || 
		len(dataPointNames) != int(numDataPoints) || 
		len(dataPointTypes) != int(numDataPoints) || 
		len(dataPointUnits) != int(numDataPoints) {
		return nil, ErrInvalidDataPoints
	}
	a := &DataStreamAttribute {
		Description: desc,
		NumDataPoints: numDataPoints,
		DataPointNames: dataPointNames,
		DataPointTypes: dataPointTypes,
		DataPointUnits: dataPointUnits,
        ProjectId: projectId,
        DomainId: domainId,
	}

	_, err := Engine.Insert(a)
	if err != nil {
		return nil, err
	}
	
	return a, nil
}
func DeleteDataStreamAttribute(id int64) error {
    a := &DataStreamAttribute{}
    _, err := Engine.Id(id).Delete(a)
    return err
}



// DeviceId links to AggregationDevice, the ProjectId & DomainId in AggregationDevice
// should exactly match the ProjectId & DomainId in DataStreamAttribute
type DataStream struct {
	Id int64
	DeviceId int64
	DataStreamAttributeId int64

	CreatedAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}

func CreateDataStreamTable() error {
    s := &DataStream{}
	_ = Engine.DropTables(s)
	err := Engine.CreateTables(s)
    return err   
}
func InsertDataStream(s *DataStream) error {
    _, err := Engine.Insert(s)
    return err
}
func CreateDataStream(deviceId int64, dataStreamAttributeId int64) (*DataStream, error) {
	s := &DataStream {
		DeviceId: deviceId,
		DataStreamAttributeId: dataStreamAttributeId,
	}
	
	_, err := Engine.Insert(s)
	if err != nil {
		return nil, err
	}
	
	return s, nil
}

func GetDataStream(id int64) (*DataStream, error) {
	dataStream := new(DataStream)
	has, err := Engine.Id(id).Get(dataStream)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotFound
	}
	return dataStream, nil
}

func GetDataStreamAttributeByDataStreamId(dataStreamId int64) (*DataStreamAttribute, error) {
	s, err := GetDataStream(dataStreamId)
	if err != nil {
		return nil, err
	}
	return GetDataStreamAttribute(s.DataStreamAttributeId)
}

func DeleteDataStream(id int64) error {
    s := &DataStream{}
    _, err := Engine.Id(id).Delete(s)
    return err
}
