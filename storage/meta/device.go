package meta

// Device meta, including
// - Device
// - AggregationDevice
//
// Note that aggregation device is a Keystone user with username/password
import (
    "time"
)

type AggregationDevice struct {
    Id string  `xorm:"pk"` //Same as keystone Id (aggregation device must be a keystone user)
	Description string `xorm:"varchar(255) notnull"`
	Latitude float64 `xorm:"default 0"`
	Longitude float64 `xorm:"default 0"`
    ProjectId string `xorm:"index"` //keystone project id
    DomainId string `xorm:"index"` //keystone domain id
	CreatedAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}

func CreateAggregationDeviceTable() error {
    a := &AggregationDevice{}
	_ = Engine.DropTables(a)
	err := Engine.CreateTables(a)
    return err
}
func GetAggregationDevice(id string) (*AggregationDevice, error) {
    a := &AggregationDevice{}
    has, err := Engine.Id(id).Get(a)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, ErrNotFound
    }
    return a, nil
}
func InsertAggregationDevice(a *AggregationDevice) error {
    _, err := Engine.Insert(a)
    return err
}
func DeleteAggregationDevice(id string) error {
    a := &AggregationDevice{}
    _, err := Engine.Id(id).Delete(a)
    return err
}

type Device struct {
	Id int64
	AggregationDeviceId string `xorm:"index"`
	Description string `xorm:"varchar(255) notnull unique"`
	Latitude float64
	Longitude float64
	CreateAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}


func CreateDeviceTable() error {
    d := &Device{}
	_ = Engine.DropTables(d)
	err := Engine.CreateTables(d)
    return err
}
func GetDevice(id int64) (*Device, error) {
    d := &Device{}
    has, err := Engine.Id(id).Get(d)
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, ErrNotFound
    }
    return d, nil
}
func InsertDevice(d *Device) error {
    _, err := Engine.Insert(d)
    return err
}
func DeleteDevice(id int64) error {
    d := &Device{}
    _, err := Engine.Id(id).Delete(d)
    return err
}
