package operation

//Defines possible operations on data sent from user

import (
	
)

const {
	OperationFilter = iota
	OperationAdd = iota
	OperationMulply = iota
	OperationConstant = iota
	OperationCustom = iota
}

type DataPoint struct {
	DataStreamId	int64
	DataPointName	string
	DataPointType	string
	DataPointUnit	int64
}

type Operation struct {
	TargetDataPoint DataPoint
	DestDataPoint DataPoint
	OperationType
}