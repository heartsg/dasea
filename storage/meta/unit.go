package meta

import (
	"time"
)

// UnitCategory
// Unit
// Currency
// All the above three meta data are currently loaded in memory and they cannot
// be modified (may be modifiable in future) currently.
// So they are not needed to be saved in database.
// We can just use cache values (in memory) to directly use them.
//
// Example:
//  units := meta.GetUnitsCache()
//  categories := meta.GetUnitCategoriesCache()
//  category := categories[units[meta.UKilometer].CategoryId]

type UnitCategory struct {
	Id int64
	Name string `xorm:"varchar(64) notnull"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func CreateUnitCategoryTable() error {
    u := &UnitCategory{}
	_ = Engine.DropTables(u)
	err := Engine.CreateTables(u)
    if err != nil {
        return err
    }
	for _, category := range unitCategories {
		_, err = Engine.Insert(category)
        if err != nil {
            return err
        }
	}
    return nil
}
func GetUnitCategoriesCache() map[int]*UnitCategory {
    return unitCategories
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

func CreateUnitTable() error {
    u := &Unit{}
	_ = Engine.DropTables(u)
	err := Engine.CreateTables(u)
    if err != nil {
        return err
    }
	for _, unit := range units {
		_, err = Engine.Insert(unit)
        if err != nil {
            return err
        }
	}
    return nil
}

func GetUnitsCache() map[int]*Unit {
    return units
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

func CreateCurrency() error {
    c := &Currency{}
    err := Engine.DropTables(c)
    if err != nil {
        return err
    }
    err = Engine.CreateTables(c)
    if err != nil {
        return err
    }
    return nil
    //currently not support currencies
}


// Id for unit categories
const (
    UcUnit = iota + 1
    UcCount
    UcPercentage
)
const (
    UcCirculatingCurrency = iota + 101
    UcDigitalCurrency
    UcHistoricalCurrency
    UcFictionalCurrency
)
const (
    UcReactivePower = iota + 201
    UcReactiveEnergy
)
const (
    UcLength = iota + 1001
    UcMass
    UcTime
    UcElectricCurrent
    UcThermodynamicTemperature
    UcAmountOfSubstance
    UcLuminousIntensity
    UcPlaneAngle
    UcSolidAngle
    UcPressure
    UcEnergy
    UcPower
    UcForce
    UcMagneticField
    UcInductance
    UcElectricCharge
    UcVoltage
    UcElectricCapacitance
    UcElectricalConductance
    UcMagneticFlux
    UcElectricResistance
    UcIlluminance
    UcLuminousFlux
    UcRadioactivity
    UcAbsorbedDose
    UcEquivalentDose
    UcFrequency
    UcCatalyticActivity
    UcArea
    UcVolume
    UcLineDensity
    UcAreaDensity
    UcVolumeDensity
    UcFuelConsumption
    UcVelocity
    UcAcceleration
    UcAngularVelocity
    UcFlowRate
    UcTorque
    UcIrradiance
    UcRelativeHumidity
    UcRssi
    UcStrain
    UcPar
)	

var unitCategories = map[int]*UnitCategory {
	//Id 1 - 100 reserved for no unit
	UcUnit: {Id: UcUnit, Name: "Unit"},
	UcCount: {Id: UcCount, Name: "Count"},
	UcPercentage: {Id: UcPercentage, Name: "Percentage"},
	
	//Id 101-200 reserved for currency/financial related
	UcCirculatingCurrency: {Id: UcCirculatingCurrency, Name: "Circulating Currency"},
	UcDigitalCurrency: {Id: UcDigitalCurrency, Name: "Digital Currency"},
	UcHistoricalCurrency: {Id: UcHistoricalCurrency, Name: "Historical Currency"},
	UcFictionalCurrency: {Id: UcFictionalCurrency, Name: "Fictional Currency"},
	
	//Id 201-300 reserved for reactive calculations
	UcReactivePower: {Id: UcReactivePower, Name: "Reactive Power"},
	UcReactiveEnergy: {Id: UcReactiveEnergy, Name: "Reactive Energy"},
	
	//Id 1001- for measurement units
	UcLength: {Id: UcLength, Name: "Length"},
	UcMass: {Id: UcMass, Name: "Mass"},
	UcTime: {Id: UcTime, Name: "Time"},
	UcElectricCurrent: {Id: UcElectricCurrent, Name: "Electric Current"},
	UcThermodynamicTemperature: {Id: UcThermodynamicTemperature, Name: "Thermodynamic Temperature"},
	UcAmountOfSubstance: {Id: UcAmountOfSubstance, Name: "Amount of Substance"},
	UcLuminousIntensity: {Id: UcLuminousIntensity, Name: "Luminous Intensity"},
	UcPlaneAngle: {Id: UcPlaneAngle, Name: "Plane Angle"},
	UcSolidAngle: {Id: UcSolidAngle, Name: "Solid Angle"},
	UcPressure: {Id: UcPressure, Name: "Pressure"},
	UcEnergy: {Id: UcEnergy, Name: "Energy"},
	UcPower: {Id: UcPower, Name: "Power"},
	UcForce: {Id: UcForce, Name: "Force"},
	UcMagneticField: {Id: UcMagneticField, Name: "Magnetic Field"},
	UcInductance: {Id: UcInductance, Name: "Inductance"},
	UcElectricCharge: {Id: UcElectricCharge, Name: "Electric Charge"},
	UcVoltage: {Id: UcVoltage, Name: "Voltage"},
	UcElectricCapacitance: {Id: UcElectricCapacitance, Name: "Electric Capacitance"},
	UcElectricalConductance: {Id: UcElectricalConductance, Name: "Electrical Conductance"},
	UcMagneticFlux: {Id: UcMagneticFlux, Name: "Magnetic Flux"},
	UcElectricResistance: {Id: UcElectricResistance, Name: "Electric Resistance"},
	UcIlluminance: {Id: UcIlluminance, Name: "Illuminance"},
    UcLuminousFlux: {Id: UcLuminousFlux, Name: "Luminous Flux"},
	UcRadioactivity: {Id: UcRadioactivity, Name: "Radioactivity"},
	UcAbsorbedDose: {Id: UcAbsorbedDose, Name: "Absorbed Dose"},
	UcEquivalentDose: {Id: UcEquivalentDose, Name: "Equivalent Dose"},
	UcFrequency: {Id: UcFrequency, Name: "Frequency"},
	UcCatalyticActivity: {Id: UcCatalyticActivity, Name: "Catalytic Activity"},
	UcArea: {Id: UcArea, Name: "Area"},
	UcVolume: {Id: UcVolume, Name: "Volume"},
	UcLineDensity: {Id: UcLineDensity, Name: "Line Density"},
	UcAreaDensity: {Id: UcAreaDensity, Name: "Area Density"},
	UcVolumeDensity: {Id: UcVolumeDensity, Name: "Volume Density"},
	UcFuelConsumption: {Id: UcFuelConsumption, Name: "Fuel Consumption"},
	UcVelocity: {Id: UcVelocity, Name: "Velocity"},
	UcAcceleration: {Id: UcAcceleration, Name: "Acceleration"},
	UcAngularVelocity: {Id: UcAngularVelocity, Name: "Angular Velocity"},
	UcFlowRate: {Id: UcFlowRate, Name: "Flow Rate"},
	UcTorque: {Id: UcTorque, Name: "Torque"},
	UcIrradiance: {Id: UcIrradiance, Name: "Irradiance"},
	UcRelativeHumidity: {Id: UcRelativeHumidity, Name: "Relative Humidity"},
	UcRssi: {Id: UcRssi, Name: "RSSI"},
	UcStrain: {Id: UcStrain, Name: "Strain"},
	UcPar: {Id: UcPar, Name: "PAR"},

	//stock?
	//derivatives?
}

// const for Units
const (
    UUnit = iota + UcUnit * 10000 + 1
)
const (
    UCount = iota + UcCount * 10000 + 1
)
const (
    UPercentage = iota + UcPercentage * 10000 + 1
)
const (
    UVoltAmpereReactive = iota + UcReactivePower * 10000 + 1
)
const (
    UVoltAmpereReactiveHour = iota + UcReactiveEnergy * 10000 + 1
)
const (
    UMeter = iota + UcLength * 10000 + 1
    UKilometer
    UDecimeter
    UCentimeter
    UMillimeter
    UMicrometer
    ULightYear
    UFoot
    UInch
    UYard
    UMile
    UNauticalMile
    UHundredthOfAnInch
    UThousandthOfAnInch
)
const (
    UKilogram = iota + UcMass * 10000 + 1
    UGram
    UPound
    UOunce
    UTonUK
    UTonUS
    UTonne
)
const (
    USecond = iota + UcTime * 10000 + 1
    UNanosecond
    UMicrosecond
    UMillisecond
    UMinute
    UHour
    UDay
)
const (
    UAmpere = iota + UcElectricCurrent * 10000 + 1
    UMilliampere
    UMicroampere
)
const (
    UKelvin = iota + UcThermodynamicTemperature * 10000 + 1
    UDegreeCelsius
    UDegreeFahrenheit
)
const (
    UMole = iota + UcAmountOfSubstance * 10000 + 1
)	
const (
    UCandela = iota + UcLuminousIntensity * 10000 + 1
)
const (
    URadian = iota + UcPlaneAngle * 10000 + 1
    UDegree
)
const (
    USteradian = iota + UcSolidAngle * 10000 + 1
)
const (
    UPascal = iota + UcPressure * 10000 + 1
    UHectopascal
    UKilopascal
    UMegapascal
    UStandardAtmosphere
    UBar
    UMillibar
    UCentimeterOfMercury
    UMillimeterOfMercury
    UInchOfMercury
    UCentimeterOfWater
    UMillimeterOfWater
    UMeterOfWater
    UInchOfWater
    UFootOfWater
    UNewtonPerSquareMeter
    UPoundPerSquareInch
)
const (
    UJoule = iota + UcEnergy * 10000 + 1
    UKilojoule
    UMegajoule
    UGigajoule
    UWattSecond
    UWattHour
    UKilowattHour
    UNewtonMeter
    UCalorieFood
    UCalorieInternational
    UCalorieThermochemical
    UCalorieMean
    UCalorie15
    UCalorie20
    UHorsepowerHour
)
const (
    UWatt = iota + UcPower * 10000 + 1
    UMilliwatt
    UKilowatt
    UMegawatt
    UGigawatt
    UHorsepowerElectric
    UHorsepowerMetric
)
const (
    UNewton = iota + UcForce * 10000 + 1
    UKilonewton
    UMeganewton
    UPoundForce
)
const (
    UTesla = iota + UcMagneticField * 10000 + 1
)
const (
    UHenry = iota + UcInductance * 10000 + 1
)
const (
    UCoulomb = iota + UcElectricCharge * 10000 + 1
)
const (
    UVolt = iota + UcVoltage * 10000 + 1
    UMillivolt
    UMicrovolt
)
const (
    UFarad = iota + UcElectricCapacitance * 10000 + 1
    UMillifarad
    UMicrofarad
    UNanofarad
    UPicofarad
)
const (
    USiemens = iota + UcElectricalConductance * 10000 + 1
)
const (
    UWeber = iota + UcMagneticFlux * 10000 + 1
)
const (
    UOhm = iota + UcElectricResistance * 10000 + 1
)
const (
    ULux = iota + UcIlluminance * 10000 + 1
)
const (
    ULumen = iota + UcLuminousFlux * 10000 + 1
)
const (
    UBecquerel = iota + UcRadioactivity * 10000 + 1
)
const (
    UGray = iota + UcAbsorbedDose * 10000 + 1
)
const (
    USievert = iota + UcEquivalentDose * 10000 + 1
)
const (
    UHertz = iota + UcFrequency * 10000 + 1
    UKilohertz
    UMegahertz
    UGigahertz
)
const (
    UKatal = iota + UcCatalyticActivity * 10000 + 1
)
const (
    USquareMeter = iota + UcArea * 10000 + 1
    USquareFoot
    USquareInch
    USquareKilometer
    UAcre
)
const (
    UCubicMeter = iota + UcVolume * 10000 + 1
    UCubicDecimeter
    UCubicCentimeter
    UCubicMillimeter
    UCubicFoot
    UCubicInch
    ULitre
    UGallonUK
    UGallonUSDry
    UGallonUSLiquid
    UBarrel
)
const (
    UKilogramPerMeter = iota + UcLineDensity * 10000 + 1
)
const (
    UKilogramPerSquareMeter = iota + UcAreaDensity * 10000 + 1
)
const (
    UKilogramPerCubicMeter = iota + UcVolumeDensity * 10000 + 1
    UKilogramPerLitre
)
const (
    UKilometerPerLitre = iota + UcFuelConsumption * 10000 + 1
)
const (
    UMeterPerSecond = iota + UcVelocity * 10000 + 1
    UKilometerPerSecond
    UKilometerPerHour
    UMilePerHour
    UKnot
)
const (
    UMeterPerSquareSecond = iota + UcAcceleration * 10000 + 1
)
const (
    URadianPerSecond = iota + UcAngularVelocity * 10000 + 1
)
const (
    UCubicMeterPerSecond = iota + UcFlowRate * 10000 + 1
)
const (
    UNewtonMeterTorque = iota + UcTorque * 10000 + 1
    UPoundForceFoot
    UPoundForceInch
)
const (
    UWattPerSquareMeter = iota + UcIrradiance * 10000 + 1
)
const (
    URelativeHumidity = iota + UcRelativeHumidity * 10000 + 1
)
const (
    UDbm = iota + UcRssi * 10000 + 1
)
const (
    UStrain = iota + UcStrain * 10000 + 1
)
const (
    UMicroEinstein = iota + UcPar * 10000 + 1
)
	
var units = map[int]*Unit {
	//Unit
	UUnit: {Id: UUnit, CategoryId: UcUnit, Name: "Unit", PluralName: "Units", Symbol: "units",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Count
	UCount: {Id: UCount, CategoryId: UcCount, Name: "Count", PluralName: "Counts", Symbol: "counts",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Percentage
	UPercentage: {Id: UPercentage, CategoryId: UcPercentage, Name: "Percent", PluralName: "Percent", Symbol: "%",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Reactive Power
	UVoltAmpereReactive: {Id: UVoltAmpereReactive, CategoryId: UcReactivePower, Name: "Volt-Ampere Reactive", PluralName: "Volt-Amperes Reactive", Symbol: "volt-amperes reactive", 
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Reactive Energy
	UVoltAmpereReactiveHour: {Id: UVoltAmpereReactiveHour, CategoryId: UcReactiveEnergy, Name: "Volt-Ampere Reactive Hour", PluralName: "Volt-Amperes Reactive Hour", Symbol: "volt-amperes reactive hour", 
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Length
	UMeter: {Id: UMeter, CategoryId: UcLength, Name: "Meter", PluralName: "Meters", Symbol: "m", 
		IsConversionBase: true, ConversionFactor: 1.0, IsMetricSystem: true},
	UKilometer: {Id: UKilometer, CategoryId: UcLength, Name: "Kilometer", PluralName: "Kilometers", Symbol: "km", 
		IsConversionBase: false, ConversionFactor: 1000.0, IsMetricSystem: true},
	UDecimeter: {Id: UDecimeter, CategoryId: UcLength, Name: "Decimeter", PluralName: "Decimeters", Symbol: "dm", 
		IsConversionBase: false, ConversionFactor: 0.1, IsMetricSystem: true},
	UCentimeter: {Id: UCentimeter, CategoryId: UcLength, Name: "Centimeter", PluralName: "Centimeters", Symbol: "cm", 
		IsConversionBase: false, ConversionFactor: 0.01, IsMetricSystem: true},
	UMillimeter: {Id: UMillimeter, CategoryId: UcLength, Name: "Millimeter", PluralName: "Millimeters", Symbol: "mm", 
		IsConversionBase: false, ConversionFactor: 0.001, IsMetricSystem: true},
	UMicrometer: {Id: UMicrometer, CategoryId: UcLength, Name: "Micrometer", PluralName: "Micrometers", Symbol: "\u03BCm", 
		IsConversionBase: false, ConversionFactor: 0.000001, IsMetricSystem: true},
	ULightYear: {Id: ULightYear, CategoryId: UcLength, Name: "Light Year", PluralName: "Light Years", Symbol: "light years",
		IsConversionBase: false, ConversionFactor: 9460500000000000},
	UFoot: {Id: UFoot, CategoryId: UcLength, Name: "Foot", PluralName: "Feet", Symbol: "ft",
		IsConversionBase: false, ConversionFactor: 0.3048, IsUSSystem: true, IsUKSystem: true},
	UInch: {Id: UInch, CategoryId: UcLength, Name: "Inch", PluralName: "Inches", Symbol: "inches",
		IsConversionBase: false, ConversionFactor: 0.0254, IsUSSystem: true, IsUKSystem: true},
	UYard: {Id: UYard, CategoryId: UcLength, Name: "Yard", PluralName: "Yards", Symbol: "yards",
		IsConversionBase: false, ConversionFactor: 0.9144, IsUSSystem: true, IsUKSystem: true},
	UMile: {Id: UMile, CategoryId: UcLength, Name: "Mile", PluralName: "Miles", Symbol: "miles",
		IsConversionBase: false, ConversionFactor: 1609.344, IsUSSystem: true, IsUKSystem: true},
	UNauticalMile: {Id: UNauticalMile, CategoryId: UcLength, Name: "Nautical Mile", PluralName: "Nautical Miles", Symbol: "nautical miles",
		IsConversionBase: false, ConversionFactor: 1852},
	UHundredthOfAnInch: {Id: UHundredthOfAnInch, CategoryId: UcLength, Name: "Hundredth of an Inch", PluralName: "Hundredths of an Inch", Symbol: "hundredths of an inch",
		IsConversionBase: false, ConversionFactor: 0.000254, IsUSSystem: true, IsUKSystem: true},
	UThousandthOfAnInch: {Id: UThousandthOfAnInch, CategoryId: UcLength, Name: "Thousandth of an Inch", PluralName: "Thousandths of an Inch", Symbol: "thousandths of an inches",
		IsConversionBase: false, ConversionFactor: 0.0000254, IsUSSystem: true, IsUKSystem: true},
	
	//Mass
	UKilogram: {Id: UKilogram, CategoryId: UcMass, Name: "Kilogram", PluralName: "Kilograms", Symbol: "kg",
		IsConversionBase: true, ConversionFactor: 1.0, IsMetricSystem: true},
	UGram: {Id: UGram, CategoryId: UcMass, Name: "Gram", PluralName: "grams", Symbol: "g",
		IsConversionBase: false, ConversionFactor: 0.001, IsMetricSystem: true},
	UPound: {Id: UPound, CategoryId: UcMass, Name: "Pound", PluralName: "Pounds", Symbol: "pounds",
		IsConversionBase: false, ConversionFactor: 0.45359237, IsUSSystem: true, IsUKSystem: true},
	UOunce: {Id: UOunce, CategoryId: UcMass, Name: "Ounce", PluralName: "Ounces", Symbol: "ounces",
		IsConversionBase: false, ConversionFactor: 0.0311034768, IsUSSystem: true, IsUKSystem: true},
	UTonUK: {Id: UTonUK, CategoryId: UcMass, Name: "Ton", PluralName: "Tons", Symbol: "tons",
		IsConversionBase: false, ConversionFactor: 1016.0469088, IsUKSystem: true},
	UTonUS: {Id: UTonUS, CategoryId: UcMass, Name: "Ton", PluralName: "Tons", Symbol: "tons",
		IsConversionBase: false, ConversionFactor: 907.18474, IsUSSystem: true},
	UTonne: {Id: UTonne, CategoryId: UcMass, Name: "Tonne", PluralName: "Tonnes", Symbol: "tonnes",
		IsConversionBase: false, ConversionFactor: 1000, IsMetricSystem: true},
		
	//Time
	USecond: {Id: USecond, CategoryId: UcTime, Name: "Second", PluralName: "Seconds", Symbol: "s",
		IsConversionBase: true, ConversionFactor: 1.0},
	UNanosecond: {Id: UNanosecond, CategoryId: UcTime, Name: "Nanosecond", PluralName: "Nanoseconds", Symbol: "ns",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	UMicrosecond: {Id: UMicrosecond, CategoryId: UcTime, Name: "Microsecond", PluralName: "Microseconds", Symbol: "\u03BCs",
		IsConversionBase: false, ConversionFactor: 0.000001},
	UMillisecond: {Id: UMillisecond, CategoryId: UcTime, Name: "Millisecond", PluralName: "Milliseconds", Symbol: "ms",
		IsConversionBase: false, ConversionFactor: 0.001},
	UMinute: {Id: UMinute, CategoryId: UcTime, Name: "Minute", PluralName: "Minutes", Symbol: "minutes",
		IsConversionBase: false, ConversionFactor: 60},
	UHour: {Id: UHour, CategoryId: UcTime, Name: "Hour", PluralName: "Hours", Symbol: "hours",
		IsConversionBase: false, ConversionFactor: 3600},
	UDay: {Id: UDay, CategoryId: UcTime, Name: "Day", PluralName: "Days", Symbol: "days",
		IsConversionBase: false, ConversionFactor: 3600*24},
	
	//Electric Current
	UAmpere: {Id: UAmpere, CategoryId: UcElectricCurrent, Name: "Ampere", PluralName: "Amperes", Symbol: "A",
		IsConversionBase: true, ConversionFactor: 1.0},
	UMilliampere: {Id: UMilliampere, CategoryId: UcElectricCurrent, Name: "Milliampere", PluralName: "Milliamperes", Symbol: "mA",
		IsConversionBase: false, ConversionFactor: 0.001},
	UMicroampere: {Id: UMicroampere, CategoryId: UcElectricCurrent, Name: "Microampere", PluralName: "Microamperes", Symbol: "\u03BCA",
		IsConversionBase: false, ConversionFactor: 0.000001},
	
	//Thermodynamic Temperature
	UKelvin: {Id: UKelvin, CategoryId: UcThermodynamicTemperature, Name: "Kelvin", PluralName: "Kelvins", Symbol: "K",
		IsConversionBase: true, ConversionFactor: 1.0},
	UDegreeCelsius: {Id: UDegreeCelsius, CategoryId: UcThermodynamicTemperature, Name: "Degree Celsius", PluralName: "Degrees Celsius", Symbol: "\u2103",
		IsConversionBase: false, ConversionFactor: 0},
	UDegreeFahrenheit: {Id: UDegreeFahrenheit, CategoryId: UcThermodynamicTemperature, Name: "Degree Fahrenheit", PluralName: "Degrees Fahrenheit", Symbol: "\u2109",
		IsConversionBase: false, ConversionFactor: 0},
	
	//Amount of Substance
	UMole: {Id: UMole, CategoryId: UcAmountOfSubstance, Name: "Mole", PluralName: "Moles", Symbol: "mol",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Luminous Intensity
	UCandela: {Id: UCandela, CategoryId: UcLuminousIntensity, Name: "Candela", PluralName: "Candelas", Symbol: "cd",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Plane Angle
	URadian: {Id: URadian, CategoryId: UcPlaneAngle, Name: "Radian", PluralName: "Radians", Symbol: "rad",
		IsConversionBase: true, ConversionFactor: 1.0},
	UDegree: {Id: UDegree, CategoryId: UcPlaneAngle, Name: "Degree", PluralName: "Degrees", Symbol: "\u00B0",
		IsConversionBase: false, ConversionFactor: 360/3.1415926535},
		
	//Solid Angle
	USteradian: {Id: USteradian, CategoryId: UcSolidAngle, Name: "Steradian", PluralName: "Steradian", Symbol: "sr",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Pressure or Stress
	UPascal: {Id: UPascal, CategoryId: UcPressure, Name: "Pascal", PluralName: "Pascals", Symbol: "Pa",
		IsConversionBase: true, ConversionFactor: 1.0},
	UHectopascal: {Id: UHectopascal, CategoryId: UcPressure, Name: "Hectopascal", PluralName: "Hectopascals", Symbol:"hPa",
		IsConversionBase: false, ConversionFactor: 100},
	UKilopascal: {Id: UKilopascal, CategoryId: UcPressure, Name: "Kilopascal", PluralName: "Kilopascals", Symbol:"kPa",
		IsConversionBase: false, ConversionFactor: 1000},
	UMegapascal: {Id: UMegapascal, CategoryId: UcPressure, Name: "Megapascal", PluralName: "Megapascals", Symbol:"MPa",
		IsConversionBase: false, ConversionFactor: 1000000},
	UStandardAtmosphere: {Id: UStandardAtmosphere, CategoryId: UcPressure, Name: "Standard Atmosphere", PluralName: "Standard Atmospheres", Symbol:"standard atmospheres",
		IsConversionBase: false, ConversionFactor: 101325},
	UBar: {Id: UBar, CategoryId: UcPressure, Name: "Bar", PluralName: "Bars", Symbol:"bars",
		IsConversionBase: false, ConversionFactor: 100000},
	UMillibar: {Id: UMillibar, CategoryId: UcPressure, Name: "Millibar", PluralName: "Millibars", Symbol:"millibars",
		IsConversionBase: false, ConversionFactor: 100},
	UCentimeterOfMercury: {Id: UCentimeterOfMercury, CategoryId: UcPressure, Name: "Centimeter of Mercury", PluralName: "Centimeters of Mercury", Symbol:"centimeters of mercury",
		IsConversionBase: false, ConversionFactor: 1333.22},
	UMillimeterOfMercury: {Id: UMillimeterOfMercury, CategoryId: UcPressure, Name: "Millimeter of Mercury", PluralName: "Millimeters of Mercury", Symbol:"millimeters of mercury",
		IsConversionBase: false, ConversionFactor: 133.322},
	UInchOfMercury: {Id: UInchOfMercury, CategoryId: UcPressure, Name: "Inch of Mercury", PluralName: "Inches of Mercury", Symbol:"inches of mercury",
		IsConversionBase: false, ConversionFactor: 3386.388},
	UCentimeterOfWater: {Id: UCentimeterOfWater, CategoryId: UcPressure, Name: "Centimeter of Water", PluralName: "Centimeters of Water", Symbol:"centimeters of water",
		IsConversionBase: false, ConversionFactor: 98.0665},
	UMillimeterOfWater: {Id: UMillimeterOfWater, CategoryId: UcPressure, Name: "Millimeter of Water", PluralName: "Millimeters of Water", Symbol:"millimeters of water",
		IsConversionBase: false, ConversionFactor: 9.80665},
	UMeterOfWater: {Id: UMeterOfWater, CategoryId: UcPressure, Name: "Meter of Water", PluralName: "Meters of Water", Symbol:"meters of water",
		IsConversionBase: false, ConversionFactor: 9806.65},
	UInchOfWater: {Id: UInchOfWater, CategoryId: UcPressure, Name: "Inch of Water", PluralName: "Inches of Water", Symbol:"inches of water",
		IsConversionBase: false, ConversionFactor: 249.08891},
	UFootOfWater: {Id: UFootOfWater, CategoryId: UcPressure, Name: "Foot of Water", PluralName: "Feet of Water", Symbol:"feet of water",
		IsConversionBase: false, ConversionFactor: 2989.06692},
	UNewtonPerSquareMeter: {Id: UNewtonPerSquareMeter, CategoryId: UcPressure, Name: "Newton Per Square Meter", PluralName: "Newtons Per Square Meter", Symbol:"newtons/square meter",
		IsConversionBase: false, ConversionFactor: 1.0},
	UPoundPerSquareInch: {Id: UPoundPerSquareInch, CategoryId: UcPressure, Name: "Pound Per Square Inch", PluralName: "Pounds Per Square Inch", Symbol:"pounds/square inch",
		IsConversionBase: false, ConversionFactor: 6894.757},

	//Energy, Work, Heat
	UJoule: {Id: UJoule, CategoryId: UcEnergy, Name: "Joule", PluralName: "Joules", Symbol: "J",
		IsConversionBase: true, ConversionFactor: 1.0},
	UKilojoule: {Id: UKilojoule, CategoryId: UcEnergy, Name: "Kilojoule", PluralName: "Kilojoules", Symbol: "kJ",
		IsConversionBase: false, ConversionFactor: 1000},
	UMegajoule: {Id: UMegajoule, CategoryId: UcEnergy, Name: "Megajoule", PluralName: "Megajoules", Symbol: "MJ",
		IsConversionBase: false, ConversionFactor: 1000000},
	UGigajoule: {Id: UGigajoule, CategoryId: UcEnergy, Name: "Gigajoule", PluralName: "Gigajoules", Symbol: "GJ",
		IsConversionBase: false, ConversionFactor: 1000000000},
	UWattSecond: {Id: UWattSecond, CategoryId: UcEnergy, Name: "Watt Second", PluralName: "Watt Seconds", Symbol: "Ws",
		IsConversionBase: false, ConversionFactor: 1},
	UWattHour: {Id: UWattHour, CategoryId: UcEnergy, Name: "Watt Hour", PluralName: "Watt Hours", Symbol: "Wh",
		IsConversionBase: false, ConversionFactor: 3600},
	UKilowattHour: {Id: UKilowattHour, CategoryId: UcEnergy, Name: "Kilowatt Hour", PluralName: "Kilowatt Hours", Symbol: "kWh",
		IsConversionBase: false, ConversionFactor: 3600000},
	UNewtonMeter: {Id: UNewtonMeter, CategoryId: UcEnergy, Name: "Newton Meter", PluralName: "Newton Meters", Symbol: "Nm",
		IsConversionBase: false, ConversionFactor: 1},
	UCalorieFood: {Id: UCalorieFood, CategoryId: UcEnergy, Name: "Calorie (Food)", PluralName: "Calories (Food)", Symbol: "calories",
		IsConversionBase: false, ConversionFactor: 4186},
	UCalorieInternational: {Id: UCalorieInternational, CategoryId: UcEnergy, Name: "Calorie (International Table)", PluralName: "Calories (International Table)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.1868},
	UCalorieThermochemical: {Id: UCalorieThermochemical, CategoryId: UcEnergy, Name: "Calorie (Thermochemical)", PluralName: "Calories (Thermochemical)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.184},
	UCalorieMean: {Id: UCalorieMean, CategoryId: UcEnergy, Name: "Calorie (Mean)", PluralName: "Calories (Mean)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.19002},
	UCalorie15: {Id: UCalorie15, CategoryId: UcEnergy, Name: "Calorie (15 \u2103)", PluralName: "Calories (15 \u2103)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.18580},
	UCalorie20: {Id: UCalorie20, CategoryId: UcEnergy, Name: "Calorie (20 \u2103)", PluralName: "Calories (20 \u2103)", Symbol: "cal",
		IsConversionBase: false, ConversionFactor: 4.18190},
	UHorsepowerHour: {Id: UHorsepowerHour, CategoryId: UcEnergy, Name: "Horsepower Hour", PluralName: "Horsepower Hours", Symbol: "horsepower hours",
		IsConversionBase: false, ConversionFactor: 2684520},

	//Power
	UWatt: {Id: UWatt, CategoryId: UcPower, Name: "Watt", PluralName: "Watts", Symbol: "W",
		IsConversionBase: true, ConversionFactor: 1.0},
	UMilliwatt: {Id: UMilliwatt, CategoryId: UcPower, Name: "Milliwatt", PluralName: "Milliwatts", Symbol: "mW",
		IsConversionBase: false, ConversionFactor: 0.001},
	UKilowatt: {Id: UKilowatt, CategoryId: UcPower, Name: "Kilowatt", PluralName: "Kilowatts", Symbol: "kW",
		IsConversionBase: false, ConversionFactor: 1000},
	UMegawatt: {Id: UMegawatt, CategoryId: UcPower, Name: "Megawatt", PluralName: "Megawatts", Symbol: "MW",
		IsConversionBase: false, ConversionFactor: 1000000},
	UGigawatt: {Id: UGigawatt, CategoryId: UcPower, Name: "Gigawatt", PluralName: "Gigawatts", Symbol: "GW",
		IsConversionBase: false, ConversionFactor: 100000000},
	UHorsepowerElectric: {Id: UHorsepowerElectric, CategoryId: UcPower, Name: "Horsepower (Electric)", PluralName: "Horsepower (Electric)", Symbol: "horsepower",
		IsConversionBase: false, ConversionFactor: 746},
	UHorsepowerMetric: {Id: UHorsepowerMetric, CategoryId: UcPower, Name: "Horsepower (Metric)", PluralName: "Horsepower (Metric)", Symbol: "horsepower",
		IsConversionBase: false, ConversionFactor: 735.499},	

	//Force, Weight
	UNewton: {Id: UNewton, CategoryId: UcForce, Name: "Newton", PluralName: "Newtons", Symbol: "N",
		IsConversionBase: true, ConversionFactor: 1.0},
	UKilonewton: {Id: UKilonewton, CategoryId: UcForce, Name: "Kilonewton", PluralName: "Kilonewtons", Symbol: "kN",
		IsConversionBase: false, ConversionFactor: 1000},
	UMeganewton: {Id: UMeganewton, CategoryId: UcForce, Name: "Meganewton", PluralName: "Meganewtons", Symbol: "MN",
		IsConversionBase: false, ConversionFactor: 1000000},
	UPoundForce: {Id: UPoundForce, CategoryId: UcForce, Name: "Pound Force", PluralName: "Pounds Force", Symbol: "pounds force",
		IsConversionBase: false, ConversionFactor: 4.448222},
	
	//Magnetic Field
	UTesla: {Id: UTesla, CategoryId: UcMagneticField, Name: "Tesla", PluralName: "Teslas", Symbol: "T",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Inductance
	UHenry: {Id: UHenry, CategoryId: UcInductance, Name: "Henry", PluralName: "Henries", Symbol: "H",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Electric Charge
	UCoulomb: {Id: UCoulomb, CategoryId: UcElectricCharge, Name: "Coulomb", PluralName: "Coulombs", Symbol: "C",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Voltage
	UVolt: {Id: UVolt, CategoryId: UcVoltage, Name: "Volt", PluralName: "Volts", Symbol: "V",
		IsConversionBase: true, ConversionFactor: 1.0},
	UMillivolt: {Id: UMillivolt, CategoryId: UcVoltage, Name: "Millivolt", PluralName: "Millivolts", Symbol: "mV",
		IsConversionBase: false, ConversionFactor: 0.001},
	UMicrovolt: {Id: UMicrovolt, CategoryId: UcVoltage, Name: "Microvolt", PluralName: "Microvolts", Symbol: "\u03BCV",
		IsConversionBase: false, ConversionFactor: 0.000001},
		
	//Electric Capacitance
	UFarad: {Id: UFarad, CategoryId: UcElectricCapacitance, Name: "Farad", PluralName: "Farads", Symbol: "F",
		IsConversionBase: true, ConversionFactor: 1.0},
	UMillifarad: {Id: UMillifarad, CategoryId: UcElectricCapacitance, Name: "Millifarad", PluralName: "Millifarads", Symbol: "mF",
		IsConversionBase: false, ConversionFactor: 0.001},
	UMicrofarad: {Id: UMicrofarad, CategoryId: UcElectricCapacitance, Name: "Microfarad", PluralName: "Microfarads", Symbol: "\u03BCF",
		IsConversionBase: false, ConversionFactor: 0.000001},
	UNanofarad: {Id: UNanofarad, CategoryId: UcElectricCapacitance, Name: "Nanofarad", PluralName: "Nanofarads", Symbol: "nF",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	UPicofarad: {Id: UPicofarad, CategoryId: UcElectricCapacitance, Name: "Picofarad", PluralName: "Picofarads", Symbol: "pF",
		IsConversionBase: false, ConversionFactor: 0.000000000001},

	//Electrical Conductance
	USiemens: {Id: USiemens, CategoryId: UcElectricalConductance, Name: "Siemens", PluralName: "Siemens", Symbol: "S",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Magnetic Flux
	UWeber: {Id: UWeber, CategoryId: UcMagneticFlux, Name: "Weber", PluralName: "Webers", Symbol: "Wb",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Electric Resistance
	UOhm: {Id: UOhm, CategoryId: UcElectricResistance, Name: "Ohm", PluralName: "Ohms", Symbol: "\u03A9",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Illuminance
	ULux: {Id: ULux, CategoryId: UcIlluminance, Name: "Lux", PluralName: "Lux", Symbol: "lx",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Luminous Flux
	ULumen: {Id: ULumen, CategoryId: UcLuminousFlux, Name: "Lumen", PluralName: "Lumens", Symbol: "lm",
		IsConversionBase: true, ConversionFactor: 1.0},

	//Radioactivity
	UBecquerel: {Id: UBecquerel, CategoryId: UcRadioactivity, Name: "Becquerel", PluralName: "Becquerels", Symbol: "Bq",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Absorbed Dose
	UGray: {Id: UGray, CategoryId: UcAbsorbedDose, Name: "Gray", PluralName: "Grays", Symbol: "Gy",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Equivalent Dose
	USievert: {Id: USievert, CategoryId: UcEquivalentDose, Name: "Sievert", PluralName: "Sieverts", Symbol: "Sv",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Frequency
	UHertz: {Id: UHertz, CategoryId: UcFrequency, Name: "Hertz", PluralName: "Hertz", Symbol: "Hz",
		IsConversionBase: true, ConversionFactor: 1.0},
	UKilohertz: {Id: UKilohertz, CategoryId: UcFrequency, Name: "Kilohertz", PluralName: "Kilohertz", Symbol: "kHz",
		IsConversionBase: false, ConversionFactor: 1000},
	UMegahertz: {Id: UMegahertz, CategoryId: UcFrequency, Name: "Megahertz", PluralName: "Megahertz", Symbol: "MHz",
		IsConversionBase: false, ConversionFactor: 1000000},
	UGigahertz: {Id: UGigahertz, CategoryId: UcFrequency, Name: "Gigahertz", PluralName: "Gigahertz", Symbol: "GHz",
		IsConversionBase: false, ConversionFactor: 1000000000},
	
	//Catalytic Activity
	UKatal: {Id: UKatal, CategoryId: UcCatalyticActivity, Name: "Katal", PluralName: "Katals", Symbol: "kat",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Area
	USquareMeter: {Id: USquareMeter, CategoryId: UcArea, Name: "Square Meter", PluralName: "Square Meters", Symbol: "square meters",
		IsConversionBase: true, ConversionFactor: 1.0},
	USquareFoot: {Id: USquareFoot, CategoryId: UcArea, Name: "Square Foot", PluralName: "Square Feet", Symbol: "square feet",
		IsConversionBase: false, ConversionFactor: 0.09290304},
	USquareInch: {Id: USquareInch, CategoryId: UcArea, Name: "Square Inch", PluralName: "Square Inches", Symbol: "square inches",
		IsConversionBase: false, ConversionFactor: 0.00064516},
	USquareKilometer: {Id: USquareKilometer, CategoryId: UcArea, Name: "Square Kilometer", PluralName: "Square Kilometers", Symbol: "square kilometers",
		IsConversionBase: false, ConversionFactor: 1000000},
	UAcre: {Id: UAcre, CategoryId: UcArea, Name: "Acre", PluralName: "Acres", Symbol: "acres",
		IsConversionBase: false, ConversionFactor: 4046.8564224},
		
	//Volume, Capacity
	UCubicMeter: {Id: UCubicMeter, CategoryId: UcVolume, Name: "Cubic Meter", PluralName: "Cubic Meters", Symbol: "cubic meters",
		IsConversionBase: true, ConversionFactor: 1.0},
	UCubicDecimeter: {Id: UCubicDecimeter, CategoryId: UcVolume, Name: "Cubic Decimeter", PluralName: "Cubic Decimeters", Symbol: "cubic decimeters",
		IsConversionBase: false, ConversionFactor: 0.001},
	UCubicCentimeter: {Id: UCubicCentimeter, CategoryId: UcVolume, Name: "Cubic Centimeter", PluralName: "Cubic Centimeters", Symbol: "cubic centimeters",
		IsConversionBase: false, ConversionFactor: 0.000001},
	UCubicMillimeter: {Id: UCubicMillimeter, CategoryId: UcVolume, Name: "Cubic Millimeter", PluralName: "Cubic Millimeters", Symbol: "cubic millimeters",
		IsConversionBase: false, ConversionFactor: 0.000000001},
	UCubicFoot: {Id: UCubicFoot, CategoryId: UcVolume, Name: "Cubic Foot", PluralName: "Cubic Feet", Symbol: "cubic feet",
		IsConversionBase: false, ConversionFactor: 0.028316846592},
	UCubicInch: {Id: UCubicInch, CategoryId: UcVolume, Name: "Cubic Inch", PluralName: "Cubic Inches", Symbol: "cubic inches",
		IsConversionBase: false, ConversionFactor: 0.000016387064},
	ULitre: {Id: ULitre, CategoryId: UcVolume, Name: "Litre", PluralName: "Liters", Symbol: "L",
		IsConversionBase: false, ConversionFactor: 0.001},
	UGallonUK: {Id: UGallonUK, CategoryId: UcVolume, Name: "Gallon (UK)", PluralName: "Gallons (UK)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.00454609},
	UGallonUSDry: {Id: UGallonUSDry, CategoryId: UcVolume, Name: "Gallon (US, Dry)", PluralName: "Gallons (US, Dry)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.00440488377086},
	UGallonUSLiquid: {Id: UGallonUSLiquid, CategoryId: UcVolume, Name: "Gallon (US, Liquid)", PluralName: "Gallons (US, Liquid)", Symbol: "gallons",
		IsConversionBase: false, ConversionFactor: 0.003785411784},
	UBarrel: {Id: UBarrel, CategoryId: UcVolume, Name: "Barrel", PluralName: "Barrels", Symbol: "barrels",
		IsConversionBase: false, ConversionFactor: 0.158987294928},
		
	//Line Density
	UKilogramPerMeter: {Id: UKilogramPerMeter, CategoryId: UcLineDensity, Name: "Kilogram Per Meter", PluralName: "Kilograms Per Meter", Symbol: "km/m",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Area Density
	UKilogramPerSquareMeter: {Id: UKilogramPerSquareMeter, CategoryId: UcAreaDensity, Name: "Kilogram Per Square Meter", PluralName: "Kilograms Per Square Meter", Symbol: "km/square meter",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Volume Density
	UKilogramPerCubicMeter: {Id: UKilogramPerCubicMeter, CategoryId: UcVolumeDensity, Name: "Kilogram Per Cubic Meter", PluralName: "Kilograms Per Cubic Meter", Symbol: "km/cubic meter",
		IsConversionBase: true, ConversionFactor: 1.0},
	UKilogramPerLitre: {Id: UKilogramPerLitre, CategoryId: UcVolumeDensity, Name: "Kilogram Per Litre", PluralName: "Kilograms Per Litre", Symbol: "km/L",
		IsConversionBase: true, ConversionFactor: 0.001},
	
	//Fuel Consumption
	UKilometerPerLitre: {Id: UKilometerPerLitre, CategoryId: UcFuelConsumption, Name: "Kilometer Per Litre", PluralName: "Kilometers Per Litre", Symbol: "km/L",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Velocity
	UMeterPerSecond: {Id: UMeterPerSecond, CategoryId: UcVelocity, Name: "Meter Per Second", PluralName: "Meters Per Second", Symbol: "m/s",
		IsConversionBase: true, ConversionFactor: 1.0},
	UKilometerPerSecond: {Id: UKilometerPerSecond, CategoryId: UcVelocity, Name: "Kilometer Per Second", PluralName: "Kilometers Per Second", Symbol: "km/s",
		IsConversionBase: false, ConversionFactor: 1000},
	UKilometerPerHour: {Id: UKilometerPerHour, CategoryId: UcVelocity, Name: "Kilometer Per Hour", PluralName: "Kilometers Per Hour", Symbol: "km/h",
		IsConversionBase: false, ConversionFactor: 1.0/3.6},
	UMilePerHour: {Id: UMilePerHour, CategoryId: UcVelocity, Name: "Mile Per Hour", PluralName: "Miles Per Hour", Symbol: "MPH",
		IsConversionBase: false, ConversionFactor: 0.44704},
	UKnot: {Id: UKnot, CategoryId: UcVelocity, Name: "Knot", PluralName: "Knots", Symbol: "knots",
		IsConversionBase: false, ConversionFactor: 0.514444},
		
	//Acceleration
	UMeterPerSquareSecond: {Id: UMeterPerSquareSecond, CategoryId: UcAcceleration, Name: "Meter Per Square Second", PluralName: "Meters Per Square Second", Symbol: "m/square second",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Angular Velocity
	URadianPerSecond: {Id: URadianPerSecond, CategoryId: UcAngularVelocity, Name: "Radian Per Second", PluralName: "Radians Per Second", Symbol: "rad/s",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Flow Rate
	UCubicMeterPerSecond: {Id: UCubicMeterPerSecond, CategoryId: UcFlowRate, Name: "Cubic Meter Per Second", PluralName: "Cubic Meters Per Second", Symbol: "cubic meters/s",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Torque		
	UNewtonMeterTorque: {Id: UNewtonMeterTorque, CategoryId: UcTorque, Name: "Newton Meter", PluralName: "Newton Meters", Symbol: "Nm",
		IsConversionBase: true, ConversionFactor: 1.0},
	UPoundForceFoot: {Id: UPoundForceFoot, CategoryId: UcTorque, Name: "Pound-Force Foot", PluralName: "Pound-Force Feet", Symbol: "pound-force feet",
		IsConversionBase: true, ConversionFactor: 1.355818},
	UPoundForceInch: {Id: UPoundForceInch, CategoryId: UcTorque, Name: "Pound-Force Inch", PluralName: "Pound-Force Inches", Symbol: "pound-force inches",
		IsConversionBase: true, ConversionFactor: 0.112984},
		
	//Irradiance
	UWattPerSquareMeter: {Id: UWattPerSquareMeter, CategoryId: UcIrradiance, Name: "Watt Per Square Meter", PluralName: "Watts Per Square Meter", Symbol: "W/square meter",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//Relative Humidity
	URelativeHumidity: {Id: URelativeHumidity, CategoryId: UcRelativeHumidity, Name: "Relative Humidity", PluralName: "Relative Humidity", Symbol: "relative humidity",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//RSSI
	UDbm: {Id: UDbm, CategoryId: UcRssi, Name: "dBm", PluralName: "dBm", Symbol: "dBm",
		IsConversionBase: true, ConversionFactor: 1.0},
	
	//Strain
	UStrain: {Id: UStrain, CategoryId: UcStrain, Name: "Strain", PluralName: "Strains", Symbol: "strains",
		IsConversionBase: true, ConversionFactor: 1.0},
		
	//PAR
	UMicroEinstein: {Id: UMicroEinstein, CategoryId: UcPar, Name: "microEinstein", PluralName: "microEinsteins", Symbol: "microEinsteins",
		IsConversionBase: true, ConversionFactor: 1.0},
}
