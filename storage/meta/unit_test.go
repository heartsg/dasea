package meta

import (
    "testing"
)


func TestConst(t *testing.T) {
    if UcUnit != 1 {
        t.Errorf("UcUnit should be 1, got %d\n", UcUnit)
    }
    if UcCount != 2 {
        t.Errorf("UcCount should be 2, got %d\n", UcCount)
    }
    if UcFictionalCurrency != 104 {
        t.Errorf("UcFictionalCurrency should be 104, got %d\n", UcFictionalCurrency)
    }
    if UcReactiveEnergy != 202 {
        t.Errorf("UcReactiveEnergy should be 202, got %d\n", UcReactiveEnergy)
    }
    if UcElectricCharge != 1016 {
        t.Errorf("UcElectricCharge should be 1016, got %d\n", UcElectricCharge)
    }
    if UcPar != 1044 {
        t.Errorf("UcPar should be 1044, got %d\n", UcPar)
    }
    
    if UUnit != 10001 {
        t.Errorf("UUnit should be 10001, got %d\n", UUnit)
    }
    if UCount != 20001 {
        t.Errorf("UCount should be 20001, got %d\n", UCount)
    }
}

// Test DB
// Must make sure the local mysql database is installed, with user dasea/password dasea and 
// db dasea
func TestUnit(t *testing.T) {
    InitEngine("mysql", []string{"dasea:dasea@tcp(127.0.0.1:3306)/dasea?charset=utf8"})
    err := CreateUnitCategoryTable()
    if err!= nil {
        t.Error(err)
    }
    
    err = CreateUnitTable()
    if err != nil {
        t.Error(err)
    }
    
    c := make([]*UnitCategory, 0)
    err = Engine.Find(&c)
    if err != nil {
        t.Error(err)
    }
    
    if len(c) != len(unitCategories) {
        t.Errorf("len of category should be %d, got: %d\n", len(unitCategories), len(c))
    }
    
    for _, i := range c {
        if i.Name != unitCategories[int(i.Id)].Name {
            t.Errorf("should be %s, got %s", unitCategories[int(i.Id)].Name, i.Name)
        }
    }
    
    u := make([]*Unit, 0)
    err = Engine.Find(&u)
    if err != nil {
        t.Error(err)
    }
    
    if len(u) != len(units) {
        t.Errorf("len of unit should be %d, got: %d\n", len(units), len(u))
    }
    
    for _, j := range u {
        if j.Name != units[int(j.Id)].Name {
            t.Errorf("should be %s, got %s", units[int(j.Id)].Name, j.Name)
        }
    }
    
    Engine.DropTables("unit_category")
    Engine.DropTables("unit")
}