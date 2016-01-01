package data


import (
	"os"
    "errors"
	"github.com/gocql/gocql"
	"github.com/go-xorm/core"
	_ "github.com/go-sql-driver/mysql"
)

var Engine *xorm.Engine

/*
//CreateData is different from other database operations
//Because the data table is dynamic for each DataStream, 
// we do not use ORM. Moreover, since data are normally not
// saved in the same database as user/device, it is good
// to leave the implementation of data operation different from
// others so it will be easier to modify later to support
// big data database such as OceanBase or Hive separately.
func CreateData(dataStreamId int64) error {
	tableName := fmt.Sprintf("data_%d", dataStreamId)
	
	dataStreamAttribute, err := GetDataStreamAttributeByDataStreamId(dataStreamId)
	if err != nil {
		return err
	}
	
	//parse dataStream types, we currently use same logic as XORM
	tmp := ""
	for i, t := range dataStreamAttribute.DataPointTypes {
		sqlType, err := TypeName2SQLType(t)
		if err != nil {
			return err
		}
		tmp = tmp + fmt.Sprintf(", %s %s not null, index(%s)", 
			dataStreamAttribute.DataPointNames[i], 
			sqlType.Name,
			dataStreamAttribute.DataPointNames[i])
	}
	
	//create sql statement to execute
	statement := fmt.Sprintf("CREATE TABLE %s (id BIGINT not null primary key auto_increment%s)", 
		tableName, tmp)
	_, err = mysqlEngine.Exec(statement)
	if err != nil {
		return err
	}

	return nil
}

//Insert data into table
func PutDataPointsFromProtobuf(dataStreamId int64, buf []byte) error {
	tableName := fmt.Sprintf("data_%d", dataStreamId)
	
	dataStreamAttribute, err := GetDataStreamAttributeByDataStreamId(dataStreamId)
	if err != nil {
		return err
	}
	
	//Start parsing buf (from protobuf)
	//We cannot use Unmarshal from protobuf.proto because we do not know the structure
	// in static form, we only know the structure in memory, so we rely on dasea/common/protobuf.go
	// to decode raw data.
	dataBuffer := common.NewProtoBuffer(buf)
	for !dataBuffer.DecodeComplete() {
		err = dataBuffer.DecodeCheckKey(common.WireLengthDelimited, 1)
		if err != nil {
			return err
		}
		record, err := dataBuffer.DecodeRawBytes(false)
		if err != nil {
			return err
		}
		recordBuffer := common.NewProtoBuffer(record)
		statement := fmt.Sprintf("INSERT INTO %s (", tableName)
		for i, col := range dataStreamAttribute.DataPointNames {
			if i == int(dataStreamAttribute.NumDataPoints) - 1 {
				statement = statement + col + ")"
			} else {
				statement = statement + col + ", "
			}
		}
		statement = statement + " VALUES ("
		lastTag := uint64(0)
		for !recordBuffer.DecodeComplete() {
			wire, tag, err := recordBuffer.DecodeKey()
			if err != nil {
				return err
			}
			//check whether tag is valid
			if tag > uint64(dataStreamAttribute.NumDataPoints) {
				err = errors.New("Invalid data")
				return err
			}
			//check whether wire is valid
			wireCheck, err := TypeName2ProtobufWireType(dataStreamAttribute.DataPointTypes[tag - 1])
			if err != nil {
				return err
			}
			if wireCheck != wire {
				err = errors.New("Invalid data")
			}
			
			//tag will be in sequence, but if the value is 0 or nil
			// it will be omitted			
			//check whether tag and lastTag has gap, if so, insert defaut values
			if tag - lastTag > 1 {
				for it := lastTag + 1; it < tag; it++ {
					dv, err := TypeName2DefaultValue(dataStreamAttribute.DataPointTypes[it - 1])
					if err != nil {
						return err
					}
					statement = statement + dv + ", "
				}
			}
			
			//at last, insert value for this tag
			var tagValue string
			switch dataStreamAttribute.DataPointTypes[tag - 1] {
			case "bool", "uint", "uint8", "uint16", "uint32", "uint64":
				v1, err := recordBuffer.DecodeVarint()
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%d", v1)
			case "int", "int8", "int16", "int32", "int64":
				v1, err := recordBuffer.DecodeVarint()
				v1s := int64(v1)
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%d", v1s)
			case "sint", "sint8", "sint16", "sint32":
				v2, err := recordBuffer.DecodeZigzag32()
				if err != nil {
					return err
				}
				v2s := int32(v2)
				tagValue = fmt.Sprintf("%d", v2s)
			case "sint64":
				v3, err := recordBuffer.DecodeZigzag64()
				if err != nil {
					return err
				}
				v3s := int64(v3)
				tagValue = fmt.Sprintf("%d", v3s)
			case "fixed64":
				v4, err := recordBuffer.DecodeFixed64()
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%d", v4)
			case "sfixed64":
				v4, err := recordBuffer.DecodeFixed64()
				if err != nil {
					return err
				}
				v4s := int64(v4)
				tagValue = fmt.Sprintf("%d", v4s)
			case "double", "float64":
				v5, err := recordBuffer.DecodeFloat64()
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%f", v5)
			case "fixed32":
				v6, err := recordBuffer.DecodeFixed32()
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%d", v6)
			case "sfixed32":
				v6, err := recordBuffer.DecodeFixed32()
				if err != nil {
					return err
				}
				v6s := int64(v6)
				tagValue = fmt.Sprintf("%d", v6s)
			case "float", "float32":
				v7, err := recordBuffer.DecodeFloat32()
				if err != nil {
					return err
				}
				tagValue = fmt.Sprintf("%f", v7)
			case "timestamp", "datetime":
				v8, err := recordBuffer.DecodeTimestamp()
				if err != nil {
					return err
				}
				tagValue = "\"" + v8.Format("2006-01-02 15:04:05.000000") + "\""
			default:
				return errors.New("Invalid data")
			}
			statement = statement + tagValue
			if tag != uint64(dataStreamAttribute.NumDataPoints) {
				statement = statement + ", "
			}
			lastTag = tag
		}
		for it := lastTag + 1; it <= uint64(dataStreamAttribute.NumDataPoints); it++ {
			dv, err := TypeName2DefaultValue(dataStreamAttribute.DataPointTypes[it - 1])
			if err != nil {
				return err
			}
			if it != uint64(dataStreamAttribute.NumDataPoints) {
				statement = statement + dv + ", "
			} else {
				statement = statement + dv
			}
		}
		statement = statement + ")"

		_, err = mysqlEngine.Exec(statement)
		if err != nil {
			return err
		}
	}

	return nil
}

func PutDataPointsFromJson() error {
	return nil
}


func GetDataPointsByTime(dataStreamId int64) error {
	return nil
}
*/