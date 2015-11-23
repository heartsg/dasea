package database

import (
	"time"
	"errors"
	"dasea/common"
)

var deviceTokenExpireDuration time.Duration

func init() {
	deviceTokenExpireDuration = time.Duration(24*7) * time.Hour
}

func GenerateDeviceToken() (string, error) {
	//retry 5 times before it fail
	for i := 0; i < 5; i++ {
		token, err := common.Token(16)
		if err != nil {
			return "", err
		}
		has, err := mysqlEngine.Get(&AggregationDevice{Token:token})
		if err != nil {
			return "", err
		}
		if !has {
			return token, nil
		}
	}
	return "", errors.New("Generate token failed")
}

func CreateAggregationDevice(desc string, password string, lat float64, long float64) (*AggregationDevice, error) {
	token, err := GenerateDeviceToken()
	if err != nil {
		return nil, err
	}
	expireAt := time.Now().Add(deviceTokenExpireDuration)
	
	
	device := AggregationDevice{
		Description: desc,
		Latitude: lat,
		Longitude: long,
		Secret: common.Md5Encrypt(password),
		Token: token,
		TokenExpireAt: expireAt,
	}
	
	_, err = mysqlEngine.Insert(&device)
	
	if err != nil {
		return nil, err
	}
	
	return &device, nil
}

func CreateDevice(desc string, aggregationDeviceId int64, lat float64, long float64) (*Device, error) {
	device := Device{
		Description: desc,
		AggregationDeviceId: aggregationDeviceId,
		Latitude: lat,
		Longitude: long,
	}
	
	_, err := mysqlEngine.Insert(&device)
	if err != nil {
		return nil, err
	}
	
	return &device, nil
}

