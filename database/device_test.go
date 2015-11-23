package database

import (
	"testing"
)

func TestDevice(t *testing.T) {
	InitDevices()
	ad, err := CreateAggregationDevice("root #03-15", "abc", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ad.Id != 1 {
		t.Fatal("Aggregation device Id should be 1")
	}
	
	d, err := CreateDevice("child #01-11", ad.Id, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	
	if d.Id != 1 {
		t.Fatal("Device Id should be 1")
	}
}