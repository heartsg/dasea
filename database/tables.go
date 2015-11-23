//Currently we rely on xorm
package database

import (
	"time"
)

type UnitCategory struct {
	Id int64
	Name string `xorm:"varchar(64) notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

type Unit struct {
	Id int64
	CategoryId int64 `xorm:"index notnull"`
	Name string `xorm:"varchar(64) notnull"`
	PluralName string `xorm:"varchar(64) default NULL"`
	Symbol string `xorm:"varchar(32) default NULL"`
	IsConversionBase bool `xorm:"index"`
	ConversionFactor float64 `xorm:"default 0"`
	IsMetricSystem bool `xorm:"index default false"`
	IsUSSystem bool `xorm:"index default false is_us_system"`
	IsUKSystem bool `xorm:"index default false is_uk_system"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

type Currency struct {
	Id int64
	CategoryId int64 `xorm:"index notnull"`
	Name string `xorm:"varchar(64) notnull"`
	ISOCode string `xorm:"varchar(8) notnull iso_code"`
	Symbol string `xorm:"varchar(8) notnull"`
	State string `xorm:"varchar(64) notnull"`
	ConversionFactor float64 `xom:"default 0"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

type User struct {
	Id int64
	Email string `xorm:"varchar(128) notnull unique"`
	Password []byte `xorm:"binary(32) notnull"`
	Salt []byte `xorm:"binary(16) notnull"`
	Token string `xorm:"varchar(16) notnull unique"`
	TokenExpireAt time.Time
	IsVerified bool `xorm:"default false"`
	Level string `xorm:"varchar(16) notnull default 'manager'"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}
type UserVerification struct {
	Id int64
	UserId int64 `xorm:"unique"`
	Token string `xorm:"varchar(16) notnull"`
	TokenExpireAt time.Time		
}
type UserPasswordReset struct {
	Id int64
	UserId int64 `xorm:"unique"`
	Token string `xorm:"varchar(16) notnull"`
	TokenExpireAt time.Time	
}

type AggregationDevice struct {
	Id int64
	Description string `xorm:"varchar(255) notnull"`
	Latitude float64 `xorm:"default 0"`
	Longitude float64 `xorm:"default 0"`
	Secret []byte `xorm:"binary(16) notnull"`
	Token string `xorm:"varchar(16) notnull unique"`
	TokenExpireAt time.Time
	CreatedAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}

type Device struct {
	Id int64
	AggregationDeviceId int64
	Description string `xorm:"varchar(255) notnull unique"`
	Latitude float64
	Longitude float64
	CreateAt time.Time `xorm:"created"`
	UpdateAt time.Time `xorm:"updated"`
	DeleteAt time.Time `xorm:"deleted"`
}

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
}

type DataStream struct {
	Id int64
	DeviceId int64
	DataStreamAttributeId int64
}

type AggregationDeviceUser struct {
	Id int64
	AggregationDeviceId int64
	UserId int64
}

type UserToken struct {
	UserId int64 `xorm:"int64 notnull unique"`
	Token string `xorm:"varchar(16) notnull unique"`
}

type DeviceToken struct {
	DeviceId int64 `xorm:"int64 notnull unique"`
	Token string `xorm:"varchar(16) notnull unique"`
}

var unitCategories = []UnitCategory {
	//Id 1 - 100 reserved for no unit
	{Id: 1, Name: "Unit"},
	{Id: 2, Name: "Count"},
	{Id: 3, Name: "Percentage"},
	
	//Id 101-200 reserved for currency/financial related
	{Id: 101, Name: "Circulating Currency"},
	{Id: 102, Name: "Digital Currency"},
	{Id: 103, Name: "Historical Currency"},
	{Id: 104, Name: "Fictional Currency"},
	
	//Id 201-300 reserved for reactive calculations
	{Id: 201, Name: "Reactive Power"},
	{Id: 202, Name: "Reactive Energy"},
	
	//Id 1001- for measurement units
	{Id: 1001, Name: "Length"},
	{Id: 1002, Name: "Mass"},
	{Id: 1003, Name: "Time"},
	{Id: 1004, Name: "Electric Current"},
	{Id: 1005, Name: "Thermodynamic Temperature"},
	{Id: 1006, Name: "Amount of Substance"},
	{Id: 1007, Name: "Luminous Intensity"},
	{Id: 1008, Name: "Plane Angle"},
	{Id: 1009, Name: "Solid Angle"},
	{Id: 1010, Name: "Pressure"},
	{Id: 1011, Name: "Energy"},
	{Id: 1012, Name: "Power"},
	{Id: 1013, Name: "Force"},
	{Id: 1014, Name: "Magnetic Field"},
	{Id: 1015, Name: "Inductance"},
	{Id: 1016, Name: "Electric Charge"},
	{Id: 1017, Name: "Voltage"},
	{Id: 1018, Name: "Electric Capacitance"},
	{Id: 1019, Name: "Electrical Conductance"},
	{Id: 1020, Name: "Magnetic Flux"},
	{Id: 1021, Name: "Electric Resistance"},
	{Id: 1022, Name: "Illuminance"},
	{Id: 1023, Name: "Luminous Flux"},
	{Id: 1024, Name: "Radioactivity"},
	{Id: 1025, Name: "Absorbed Dose"},
	{Id: 1026, Name: "Equivalent Dose"},
	{Id: 1027, Name: "Frequency"},
	{Id: 1028, Name: "Catalytic Activity"},
	{Id: 1029, Name: "Area"},
	{Id: 1030, Name: "Volume"},
	{Id: 1031, Name: "Line Density"},
	{Id: 1032, Name: "Area Density"},
	{Id: 1033, Name: "Volume Density"},
	{Id: 1034, Name: "Fuel Consumption"},
	{Id: 1035, Name: "Velocity"},
	{Id: 1036, Name: "Acceleration"},
	{Id: 1037, Name: "Angular Velocity"},
	{Id: 1038, Name: "Flow Rate"},
	{Id: 1039, Name: "Torque"},
	{Id: 1040, Name: "Irradiance"},
	{Id: 1041, Name: "Relative Humidity"},
	{Id: 1042, Name: "RSSI"},
	{Id: 1043, Name: "Strain"},
	{Id: 1044, Name: "PAR"},
	
	//stock?
	//derivatives?
}

var units = []Unit {
	//Unit
	{Id: 10001, CategoryId: 1, Name: "Unit", PluralName: "Units", Symbol: "units",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Count
	{Id: 20001, CategoryId: 2, Name: "Count", PluralName: "Counts", Symbol: "counts",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Percentage
	{Id: 30001, CategoryId: 3, Name: "Percent", PluralName: "Percent", Symbol: "%",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Reactive Power
	{Id: 2010001, CategoryId: 201, Name: "Volt-Ampere Reactive", PluralName: "Volt-Amperes Reactive", Symbol: "volt-amperes reactive", 
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Reactive Energy
	{Id: 2020001, CategoryId: 202, Name: "Volt-Ampere Reactive Hour", PluralName: "Volt-Amperes Reactive Hour", Symbol: "volt-amperes reactive hour", 
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Length
	{Id: 10010001, CategoryId: 1001, Name: "Meter", PluralName: "Meters", Symbol: "m", 
		IsConversionBase: true, ConversionFactor: 1.0, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Kilometer", PluralName: "Kilometers", Symbol: "km", 
		IsConversionBase: false, ConversionFactor: 1000.0, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Decimeter", PluralName: "Decimeters", Symbol: "dm", 
		IsConversionBase: false, ConversionFactor: 0.1, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Centimeter", PluralName: "Centimeters", Symbol: "cm", 
		IsConversionBase: false, ConversionFactor: 0.01, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Millimeter", PluralName: "Millimeters", Symbol: "mm", 
		IsConversionBase: false, ConversionFactor: 0.001, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Micrometer", PluralName: "Micrometers", Symbol: "\u03BCm", 
		IsConversionBase: false, ConversionFactor: 0.000001, IsMetricSystem: true},
	{CategoryId: 1001, Name: "Light Year", PluralName: "Light Years", Symbol: "light years",
		IsConversionBase: false, ConversionFactor: 9460500000000000},
	{CategoryId: 1001, Name: "Foot", PluralName: "Feet", Symbol: "ft",
		IsConversionBase: false, ConversionFactor: 0.3048, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1001, Name: "Inch", PluralName: "Inches", Symbol: "inches",
		IsConversionBase: false, ConversionFactor: 0.0254, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1001, Name: "Yard", PluralName: "Yards", Symbol: "yards",
		IsConversionBase: false, ConversionFactor: 0.9144, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1001, Name: "Mile", PluralName: "Miles", Symbol: "miles",
		IsConversionBase: false, ConversionFactor: 1609.344, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1001, Name: "Nautical Mile", PluralName: "Nautical Miles", Symbol: "nautical miles",
		IsConversionBase: false, ConversionFactor: 1852},
	{CategoryId: 1001, Name: "Hundredth of an Inch", PluralName: "Hundredths of an Inch", Symbol: "hundredths of an inch",
		IsConversionBase: false, ConversionFactor: 0.000254, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1001, Name: "Thousandth of an Inch", PluralName: "Thousandths of an Inch", Symbol: "thousandths of an inches",
		IsConversionBase: false, ConversionFactor: 0.0000254, IsUSSystem: true, IsUKSystem: true},
	
	//Mass
	{Id: 10020001, CategoryId: 1002, Name: "Kilogram", PluralName: "Kilograms", Symbol: "kg",
		IsConversionBase: true, ConversionFactor: 1.0, IsMetricSystem: true},
	{CategoryId: 1002, Name: "Gram", PluralName: "grams", Symbol: "g",
		IsConversionBase: false, ConversionFactor: 0.001, IsMetricSystem: true},
	{CategoryId: 1002, Name: "Pound", PluralName: "Pounds", Symbol: "pounds",
		IsConversionBase: false, ConversionFactor: 0.45359237, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1002, Name: "Ounce", PluralName: "Ounces", Symbol: "ounces",
		IsConversionBase: false, ConversionFactor: 0.0311034768, IsUSSystem: true, IsUKSystem: true},
	{CategoryId: 1002, Name: "Ton", PluralName: "Tons", Symbol: "tons",
		IsConversionBase: false, ConversionFactor: 1016.0469088, IsUKSystem: true},
	{CategoryId: 1002, Name: "Ton", PluralName: "Tons", Symbol: "tons",
		IsConversionBase: false, ConversionFactor: 907.18474, IsUSSystem: true},
	{CategoryId: 1002, Name: "Tonne", PluralName: "Tonnes", Symbol: "tonnes",
		IsConversionBase: false, ConversionFactor: 1000, IsMetricSystem: true},
		
	//Time
	{Id: 10030001, CategoryId: 1003, Name: "Second", PluralName: "Seconds", Symbol: "s",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1003, Name: "Nanosecond", PluralName: "Nanoseconds", Symbol: "ns",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	{CategoryId: 1003, Name: "Microsecond", PluralName: "Microseconds", Symbol: "\u03BCs",
		IsConversionBase: false, ConversionFactor: 0.000001},
	{CategoryId: 1003, Name: "Millisecond", PluralName: "Milliseconds", Symbol: "ms",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1003, Name: "Minute", PluralName: "Minutes", Symbol: "minutes",
		IsConversionBase: false, ConversionFactor: 60},
	{CategoryId: 1003, Name: "Hour", PluralName: "Hours", Symbol: "hours",
		IsConversionBase: false, ConversionFactor: 3600},
	{CategoryId: 1003, Name: "Day", PluralName: "Days", Symbol: "days",
		IsConversionBase: false, ConversionFactor: 3600*24},
	
	//Electric Current
	{Id: 10040001, CategoryId: 1004, Name: "Ampere", PluralName: "Amperes", Symbol: "A",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1004, Name: "Milliampere", PluralName: "Milliamperes", Symbol: "mA",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1004, Name: "Microampere", PluralName: "Microamperes", Symbol: "\u03BCA",
		IsConversionBase: false, ConversionFactor: 0.000001},
	
	//Thermodynamic Temperature
	{Id: 10050001, CategoryId: 1005, Name: "Kelvin", PluralName: "Kelvins", Symbol: "K",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1005, Name: "Degree Celsius", PluralName: "Degrees Celsius", Symbol: "\u2103",
		IsConversionBase: false, ConversionFactor: 0},
	{CategoryId: 1005, Name: "Degree Fahrenheit", PluralName: "Degrees Fahrenheit", Symbol: "\u2109",
		IsConversionBase: false, ConversionFactor: 0},
	
	//Amount of Substance
	{Id: 10060001, CategoryId: 1006, Name: "Mole", PluralName: "Moles", Symbol: "mol",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Luminous Intensity
	{Id: 10070001, CategoryId: 1007, Name: "Candela", PluralName: "Candelas", Symbol: "cd",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Plain Angle
	{Id: 10080001, CategoryId: 1008, Name: "Radian", PluralName: "Radians", Symbol: "rad",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1008, Name: "Degree", PluralName: "Degrees", Symbol: "\u00B0",
		IsConversionBase: false, ConversionFactor: 360/3.1415926535},
		
	//Solid Angle
	{Id: 10090001, CategoryId: 1009, Name: "Steradian", PluralName: "Steradian", Symbol: "sr",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Pressure or Stress
	{Id: 10100001, CategoryId: 1010, Name: "Pascal", PluralName: "Pascals", Symbol: "Pa",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1010, Name: "Hectopascal", PluralName: "Hectopascals", Symbol:"hPa",
		IsConversionBase: false, ConversionFactor: 100},
	{CategoryId: 1010, Name: "Kilopascal", PluralName: "Kilopascals", Symbol:"kPa",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1010, Name: "Megapascal", PluralName: "Megapascals", Symbol:"MPa",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1010, Name: "Standard Atmosphere", PluralName: "Standard Atmospheres", Symbol:"standard atmospheres",
		IsConversionBase: false, ConversionFactor: 101325},
	{CategoryId: 1010, Name: "Bar", PluralName: "Bars", Symbol:"bars",
		IsConversionBase: false, ConversionFactor: 100000},
	{CategoryId: 1010, Name: "Millibar", PluralName: "Millibars", Symbol:"millibars",
		IsConversionBase: false, ConversionFactor: 100},
	{CategoryId: 1010, Name: "Centimeter of Mercury", PluralName: "Centimeters of Mercury", Symbol:"centimeters of mercury",
		IsConversionBase: false, ConversionFactor: 1333.22},
	{CategoryId: 1010, Name: "Millimeter of Mercury", PluralName: "Millimeters of Mercury", Symbol:"millimeters of mercury",
		IsConversionBase: false, ConversionFactor: 133.322},
	{CategoryId: 1010, Name: "Inch of Mercury", PluralName: "Inches of Mercury", Symbol:"inches of mercury",
		IsConversionBase: false, ConversionFactor: 3386.388},
	{CategoryId: 1010, Name: "Centimeter of Water", PluralName: "Centimeters of Water", Symbol:"centimeters of water",
		IsConversionBase: false, ConversionFactor: 98.0665},
	{CategoryId: 1010, Name: "Millimeter of Water", PluralName: "Millimeters of Water", Symbol:"millimeters of water",
		IsConversionBase: false, ConversionFactor: 9.80665},
	{CategoryId: 1010, Name: "Meter of Water", PluralName: "Meters of Water", Symbol:"meters of water",
		IsConversionBase: false, ConversionFactor: 9806.65},
	{CategoryId: 1010, Name: "Inch of Water", PluralName: "Inches of Water", Symbol:"inches of water",
		IsConversionBase: false, ConversionFactor: 249.08891},
	{CategoryId: 1010, Name: "Foot of Water", PluralName: "Feet of Water", Symbol:"feet of water",
		IsConversionBase: false, ConversionFactor: 2989.06692},
	{CategoryId: 1010, Name: "Newton Per Square Meter", PluralName: "Newtons Per Square Meter", Symbol:"newtons/square meter",
		IsConversionBase: false, ConversionFactor: 1.0},
	{CategoryId: 1010, Name: "Pound Per Square Inch", PluralName: "Pounds Per Square Inch", Symbol:"pounds/square inch",
		IsConversionBase: false, ConversionFactor: 6894.757},

	//Energy, Work, Heat
	{Id: 10110001, CategoryId: 1011, Name: "Joule", PluralName: "Joules", Symbol: "J",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1011, Name: "Kilojoule", PluralName: "Kilojoules", Symbol: "kJ",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1011, Name: "Megajoule", PluralName: "Megajoules", Symbol: "MJ",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1011, Name: "Gigajoule", PluralName: "Gigajoules", Symbol: "GJ",
		IsConversionBase: false, ConversionFactor: 1000000000},
	{CategoryId: 1011, Name: "Watt Second", PluralName: "Watt Seconds", Symbol: "Ws",
		IsConversionBase: false, ConversionFactor: 1},
	{CategoryId: 1011, Name: "Watt Hour", PluralName: "Watt Hours", Symbol: "Wh",
		IsConversionBase: false, ConversionFactor: 3600},
	{CategoryId: 1011, Name: "Kilowatt Hour", PluralName: "Kilowatt Hours", Symbol: "kWh",
		IsConversionBase: false, ConversionFactor: 3600000},
	{CategoryId: 1011, Name: "Newton Meter", PluralName: "Newton Meters", Symbol: "Nm",
		IsConversionBase: false, ConversionFactor: 1},
	{CategoryId: 1011, Name: "Calorie (Food)", PluralName: "Calories (Food)", Symbol: "calories",
		IsConversionBase: false, ConversionFactor: 4186},
	{CategoryId: 1011, Name: "Calorie (International Table)", PluralName: "Calories (International Table)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.1868},
	{CategoryId: 1011, Name: "Calorie (Thermochemical)", PluralName: "Calories (Thermochemical)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.184},
	{CategoryId: 1011, Name: "Calorie (Mean)", PluralName: "Calories (Mean)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.19002},
	{CategoryId: 1011, Name: "Calorie (15 \u2103)", PluralName: "Calories (15 \u2103)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.18580},
	{CategoryId: 1011, Name: "Calorie (20 \u2103)", PluralName: "Calories (20 \u2103)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.18190},
	{CategoryId: 1011, Name: "Horsepower Hour", PluralName: "Horsepower Hours", Symbol: "horsepower hours",
		IsConversionBase: false, ConversionFactor: 2684520},

	//Power
	{Id: 10120001, CategoryId: 1012, Name: "Watt", PluralName: "Watts", Symbol: "W",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1012, Name: "Milliwatt", PluralName: "Milliwatts", Symbol: "mW",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1012, Name: "Kilowatt", PluralName: "Kilowatts", Symbol: "kW",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1012, Name: "Megawatt", PluralName: "Megawatts", Symbol: "MW",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1012, Name: "Gigawatt", PluralName: "Gigawatts", Symbol: "GW",
		IsConversionBase: false, ConversionFactor: 100000000},
	{CategoryId: 1012, Name: "Milliwatt", PluralName: "Milliwatts", Symbol: "mW",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1012, Name: "Horsepower (Electric)", PluralName: "Horsepower (Electric)", Symbol: "horsepower",
		IsConversionBase: false, ConversionFactor: 746},
	{CategoryId: 1012, Name: "Horsepower (Metric)", PluralName: "Horsepower (Metric)", Symbol: "horsepower",
		IsConversionBase: false, ConversionFactor: 735.499},	

	//Force, Weight
	{Id: 10130001, CategoryId: 1013, Name: "Newton", PluralName: "Newtons", Symbol: "N",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1013, Name: "Kilonewton", PluralName: "Kilonewtons", Symbol: "kN",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1013, Name: "Meganewton", PluralName: "Meganewtons", Symbol: "MN",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1013, Name: "Pound Force", PluralName: "Pounds Force", Symbol: "pounds force",
		IsConversionBase: false, ConversionFactor: 4.448222},
	
	//Magnetic Field
	{Id: 10140001, CategoryId: 1014, Name: "Tesla", PluralName: "Teslas", Symbol: "T",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Inductance
	{Id: 10150001, CategoryId: 1015, Name: "Henry", PluralName: "Henries", Symbol: "H",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Electric Charge
	{Id: 10160001, CategoryId: 1016, Name: "Coulomb", PluralName: "Coulombs", Symbol: "C",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Voltage
	{Id: 10170001, CategoryId: 1017, Name: "Volt", PluralName: "Volts", Symbol: "V",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1017, Name: "Millivolt", PluralName: "Millivolts", Symbol: "mV",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1017, Name: "Microvolt", PluralName: "Microvolts", Symbol: "\u03BCV",
		IsConversionBase: false, ConversionFactor: 0.000001},
		
	//Electric Capacitance
	{Id: 10180001, CategoryId: 1018, Name: "Farad", PluralName: "Farads", Symbol: "F",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1018, Name: "Millifarad", PluralName: "Millifarads", Symbol: "mF",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1018, Name: "Microfarad", PluralName: "Microfarads", Symbol: "\u03BCF",
		IsConversionBase: false, ConversionFactor: 0.000001},
	{CategoryId: 1018, Name: "Nanofarad", PluralName: "Nanofarads", Symbol: "nF",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	{CategoryId: 1018, Name: "Picofarad", PluralName: "Picofarads", Symbol: "pF",
		IsConversionBase: false, ConversionFactor: 0.000000000001},

	//Electrical Conductance
	{Id: 10190001, CategoryId: 1019, Name: "Siemens", PluralName: "Siemens", Symbol: "S",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Magnetic Flux
	{Id: 10200001, CategoryId: 1020, Name: "Weber", PluralName: "Webers", Symbol: "Wb",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Electric Resistance
	{Id: 10210001, CategoryId: 1021, Name: "Ohm", PluralName: "Ohms", Symbol: "\u03A9",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Illuminance
	{Id: 10220001, CategoryId: 1022, Name: "Lux", PluralName: "Lux", Symbol: "lx",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Luminous Flux
	{Id: 10230001, CategoryId: 1023, Name: "Lumen", PluralName: "Lumens", Symbol: "lm",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Radioactivity
	{Id: 10240001, CategoryId: 1024, Name: "Becquerel", PluralName: "Becquerels", Symbol: "Bq",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Absorbed Dose
	{Id: 10250001, CategoryId: 1025, Name: "Gray", PluralName: "Grays", Symbol: "Gy",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Equivalent Dose
	{Id: 10260001, CategoryId: 1026, Name: "Sievert", PluralName: "Sieverts", Symbol: "Sv",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Frequency
	{Id: 10270001, CategoryId: 1027, Name: "Hertz", PluralName: "Hertz", Symbol: "Hz",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1027, Name: "Kilohertz", PluralName: "Kilohertz", Symbol: "kHz",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1027, Name: "Megahertz", PluralName: "Megahertz", Symbol: "MHz",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1027, Name: "Gigahertz", PluralName: "Gigahertz", Symbol: "GHz",
		IsConversionBase: false, ConversionFactor: 1000000000},
	
	//Catalytic Activity
	{Id: 10280001, CategoryId: 1028, Name: "Katal", PluralName: "Katals", Symbol: "kat",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Area
	{Id: 10290001, CategoryId: 1029, Name: "Square Meter", PluralName: "Square Meters", Symbol: "square meters",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1029, Name: "Square Foot", PluralName: "Square Feet", Symbol: "square feet",
		IsConversionBase: false, ConversionFactor: 0.09290304},
	{CategoryId: 1029, Name: "Square Inch", PluralName: "Square Inches", Symbol: "square inches",
		IsConversionBase: false, ConversionFactor: 0.00064516},
	{CategoryId: 1029, Name: "Square Kilometer", PluralName: "Square Kilometers", Symbol: "square kilometers",
		IsConversionBase: false, ConversionFactor: 1000000},
	{CategoryId: 1029, Name: "Acre", PluralName: "Acres", Symbol: "acres",
		IsConversionBase: false, ConversionFactor: 4046.8564224},
		
	//Volume, Capacity
	{Id: 10300001, CategoryId: 1030, Name: "Cubic Meter", PluralName: "Cubic Meters", Symbol: "cubic meters",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1030, Name: "Cubic Decimeter", PluralName: "Cubic Decimeters", Symbol: "cubic decimeters",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1030, Name: "Cubic Centimeter", PluralName: "Cubic Centimeters", Symbol: "cubic centimeters",
		IsConversionBase: false, ConversionFactor: 0.000001},
	{CategoryId: 1030, Name: "Cubic Millimeter", PluralName: "Cubic Millimeters", Symbol: "cubic millimeters",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	{CategoryId: 1030, Name: "Cubic Foot", PluralName: "Cubic Feet", Symbol: "cubic feet",
		IsConversionBase: false, ConversionFactor: 0.028316846592},
	{CategoryId: 1030, Name: "Cubic Inch", PluralName: "Cubic Inches", Symbol: "cubic inches",
		IsConversionBase: false, ConversionFactor: 0.000016387064},
	{CategoryId: 1030, Name: "Litre", PluralName: "Liters", Symbol: "L",
		IsConversionBase: false, ConversionFactor: 0.001},
	{CategoryId: 1030, Name: "Gallon (UK)", PluralName: "Gallons (UK)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.00454609},
	{CategoryId: 1030, Name: "Gallon (US, Dry)", PluralName: "Gallons (US, Dry)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.00440488377086},
	{CategoryId: 1030, Name: "Gallon (US, Liquid)", PluralName: "Gallons (US, Liquid)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.003785411784},
	{CategoryId: 1030, Name: "Barrel", PluralName: "Barrels", Symbol: "barrels",
		IsConversionBase: false, ConversionFactor: 0.158987294928},
		
	//Line Density
	{Id: 10310001, CategoryId: 1031, Name: "Kilogram Per Meter", PluralName: "Kilograms Per Meter", Symbol: "km/m",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Area Density
	{Id: 10320001, CategoryId: 1032, Name: "Kilogram Per Square Meter", PluralName: "Kilograms Per Square Meter", Symbol: "km/square meter",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Volume Density
	{Id: 10330001, CategoryId: 1033, Name: "Kilogram Per Cubic Meter", PluralName: "Kilograms Per Cubic Meter", Symbol: "km/cubic meter",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1033, Name: "Kilogram Per Litre", PluralName: "Kilograms Per Litre", Symbol: "km/L",
		IsConversionBase: true, ConversionFactor: 0.001},
	
	//Fuel Consumption
	{Id: 10340001, CategoryId: 1034, Name: "Kilometer Per Litre", PluralName: "Kilometers Per Litre", Symbol: "km/L",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Velocity
	{Id: 10350001, CategoryId: 1035, Name: "Meter Per Second", PluralName: "Meters Per Second", Symbol: "m/s",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1035, Name: "Kilometer Per Second", PluralName: "Kilometers Per Second", Symbol: "km/s",
		IsConversionBase: false, ConversionFactor: 1000},
	{CategoryId: 1035, Name: "Kilometer Per Hour", PluralName: "Kilometers Per Hour", Symbol: "km/h",
		IsConversionBase: false, ConversionFactor: 1.0/3.6},
	{CategoryId: 1035, Name: "Mile Per Hour", PluralName: "Miles Per Hour", Symbol: "MPH",
		IsConversionBase: false, ConversionFactor: 0.44704},
	{CategoryId: 1035, Name: "Knot", PluralName: "Knots", Symbol: "knots",
		IsConversionBase: false, ConversionFactor: 0.514444},
		
	//Acceleration
	{Id: 10360001, CategoryId: 1036, Name: "Meter Per Square Second", PluralName: "Meters Per Square Second", Symbol: "m/square second",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Angular Velocity
	{Id: 10370001, CategoryId: 1037, Name: "Radian Per Second", PluralName: "Radians Per Second", Symbol: "rad/s",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Flow Rate
	{Id: 10380001, CategoryId: 1038, Name: "Cubic Meter Per Second", PluralName: "Cubic Meters Per Second", Symbol: "cubic meters/s",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Torque		
	{Id: 10390001, CategoryId: 1039, Name: "Newton Meter", PluralName: "Newton Meters", Symbol: "Nm",
		IsConversionBase: true, ConversionFactor: 1.0},
	{CategoryId: 1039, Name: "Pound-Force Foot", PluralName: "Pound-Force Feet", Symbol: "pound-force feet",
		IsConversionBase: true, ConversionFactor: 1.355818},
	{CategoryId: 1039, Name: "Pound-Force Inch", PluralName: "Pound-Force Inches", Symbol: "pound-force inches",
		IsConversionBase: true, ConversionFactor: 0.112984},
		
	//Irradiance
	{Id: 10400001, CategoryId: 1040, Name: "Watt Per Square Meter", PluralName: "Watts Per Square Meter", Symbol: "W/square meter",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Relative Humidity
	{Id: 10410001, CategoryId: 1041, Name: "Relative Humidity", PluralName: "Relative Humidity", Symbol: "relative humidity",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//RSSI
	{Id: 10420001, CategoryId: 1042, Name: "dBm", PluralName: "dBm", Symbol: "dBm",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Strain
	{Id: 10430001, CategoryId: 1043, Name: "Strain", PluralName: "Strains", Symbol: "strains",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//PAR
	{Id: 10440001, CategoryId: 1044, Name: "microEinstein", PluralName: "microEinsteins", Symbol: "microEinsteins",
		IsConversionBase: true, ConversionFactor: 1.0},
}
