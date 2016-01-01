//Types conversion between golang and sql
//Similar to (also based on) type.go in xorm but we only consider the types
// we use for our application.
//Moreover, we also consider protobuf and json types mapping to sql.


package data

import (
	"errors"
	"github.com/go-xorm/core"
	"dasea/common"
)

//Xorm provides Type2SQL which converts from golang type to sql type
//We consider the golang type in string to sql type
//Note that xorm uses dialect to map various core.SQLType to the type that
// is really supported by that dialect, we don't do this currently, because we
// do not rely on xorm. We map directly to the type that we will support for the
// database that we use (currently mysql, and later we will move to oceanbase). 
//Such dialect support may be added later
func TypeName2SQLType(t string) (core.SQLType, error) {
	var st core.SQLType
	switch t {
	//"sint" types are specifically to indicate that the data have large number of negative values
	//used for protobuf to efficiently encode negative values.
	//If negative values are not common, use int is good enough.
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "sint", "sint8", "sint16", "sint32", "fixed32", "sfixed32":
		st = core.SQLType{core.Int, 0, 0}
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		st = core.SQLType{core.BigInt, 0, 0}
	case "float32", "float":
		st = core.SQLType{core.Float, 0, 0}
	case "float64", "double":
		st = core.SQLType{core.Double, 0, 0}
	case "complex64", "complex128":
		st = core.SQLType{core.Varchar, 64, 0}
	case "bool":
		st = core.SQLType{core.Bool, 0, 0}
	case "string":
		st = core.SQLType{core.Varchar, 255, 0}
	case "datetime", "timestamp":
		st = core.SQLType{core.DateTime, 0, 0}
	default:
		return core.SQLType{"UNKNOWN", 0, 0}, errors.New("Invalid type")
	}
	return st, nil	
}

//From type name to protobuf wire type (for encoding/decoding)
func TypeName2ProtobufWireType(t string) (uint64, error) {
	var pt uint64
	switch t {
	case "bool", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", 
		"uint32", "uint64", "sint", "sint8", "sint16", "sint32", "sint64":
		pt = common.WireVarint
	case "fixed64", "sfixed64", "double", "float64":
		pt = common.WireFixed64
	case "fixed32", "sfixed32", "float", "float32":
		pt = common.WireFixed32
	case "string", "bytes", "timestamp", "datetime":
		pt = common.WireLengthDelimited
	default:
		return 0, errors.New("Invalid type")
	}
	return pt, nil		
}

func TypeName2DefaultValue(t string) (string, error) {
	var dv string
	switch t {
	case "bool", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", 
		"uint32", "uint64", "sint", "sint8", "sint16", "sint32", "sint64",
		 "fixed32", "sfixed32", "float", "float32",
		"fixed64", "sfixed64", "double", "float64":
		dv = "0"
	case "string":
		dv = `"0"`
	case "timestamp", "datetime":
		dv = "\"2015-1-1 00:00:00\""
	default:
		return "", errors.New("Invalid type")
	}
	return dv, nil		
}