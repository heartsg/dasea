//Currently uses go-xorm (and various sql libraries)
//Will support more types (including google amazon cloud services,
//hdfs/hive based, oceandb based etc.)

package database

import (
	"fmt"
	"os"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
)

var mysqlEngine *xorm.Engine

func init() {
	var err error
	mysqlEngine, err = xorm.NewEngine("mysql", "dasea:dasea@/dasea?charset=utf8")
	logger := xorm.NewSimpleLogger(os.Stdout)
	logger.SetLevel(core.LOG_OFF)
	mysqlEngine.SetLogger(logger)
	if err != nil {
		panic(err)
	}
}

func InitUnits() {
	mysqlEngine.DropTables("unit_category")
	mysqlEngine.Sync2(new(UnitCategory))
	for _, category := range unitCategories {
		mysqlEngine.Insert(&category)
	}

	mysqlEngine.DropTables("unit")
	mysqlEngine.Sync2(new(Unit))
	for _, unit := range units {
		mysqlEngine.Insert(&unit)
	}
}

func InitUsers() {
	mysqlEngine.DropTables("user")
	mysqlEngine.DropTables("user_verification")
	mysqlEngine.DropTables("user_password_reset")
	mysqlEngine.Sync2(new(User))
	mysqlEngine.Sync2(new(UserVerification))
	mysqlEngine.Sync2(new(UserPasswordReset))
	
	//Admin is the default user that must be created during initialization
	//User name (email): admin
	//Password: admin
	//Level: admin
	admin, err := CreateUser("admin", "admin", "admin")
	if err != nil {
		panic(err)
	}
	//Auto verified
	admin.Verified()
}

func UpdateUsers() {
	mysqlEngine.Sync2(new(User))
	mysqlEngine.Sync2(new(UserVerification))
	mysqlEngine.Sync2(new(UserPasswordReset))
}

func InitDevices() {
	mysqlEngine.DropTables("aggregation_device")
	mysqlEngine.DropTables("device")
	mysqlEngine.Sync2(new(AggregationDevice))
	mysqlEngine.Sync2(new(Device))
}

func InitData() {
	dataStreams := make([]*DataStream, 0)
	mysqlEngine.Find(&dataStreams)
	for _, dataStream := range dataStreams {
		dataTableName := fmt.Sprintf("data_%d", dataStream.Id)
		mysqlEngine.DropTables(dataTableName)
	}
	mysqlEngine.DropTables("data_stream_attribute")
	mysqlEngine.DropTables("data_stream")
	mysqlEngine.Sync2(new(DataStreamAttribute))
	mysqlEngine.Sync2(new(DataStream))
}