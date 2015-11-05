package config

import "testing"


type TestValidatorStruct struct {
	A int `required:"true"`
	B int `required:"1" select:"10|20"`
	C int `required:"false" min:"10"`
	D string `regex:"xyz|abc"`
	E TestValidatorInnerStruct
}
type TestValidatorInnerStruct struct {
	A int `customMin:"10" customMax:"20"`
	B string `customRequired:"1" minlen:"10" maxlen:"20"`
	C string `required:"true" regex:"^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"`
	D string `select:"tcp|udp"`
	E int `select:"80|8080"`
}

func TestRequiredValidators(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:0,
		D:"aaa",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			C:"xyz",
			D:"bbb",
			E:80,
		},
	}

	err := (&RequiredValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.A = 0
	err = (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.A should be required.")
	}
	
	s.A = 100
	s.B = 0
	err = (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.A should be required.")
	}
}

func TestRequiredValidatorsEmbededStruct(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:0,
		D:"aaa",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			D:"bbb",
			E:80,
		},
	}
	
	err := (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.C should be required.")
	}
	
	s.E.C = ""
	err = (&RequiredValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.C should be required.")
	}
}

func TestRequiredValidatorsCustomTag(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:0,
		D:"aaa",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			D:"bbb",
			E:80,
		},
	}
	
	err := (&RequiredValidator{RequiredTagName:"customRequired"}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.A = 0
	err = (&RequiredValidator{RequiredTagName:"customRequired"}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	s.E.B = ""
	err = (&RequiredValidator{RequiredTagName:"customRequired"}).Validate(s)
	if err == nil {
		t.Fatal("s.E.B should be required.")
	}
}


func TestRangeValidators(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:10,
		D:"aaa",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			C:"xyz",
			D:"bbb",
			E:80,
		},
	}
	err := (&RangeValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.C = 9
	err = (&RangeValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.C range should be checked.")
	}
	
	s.C = 0 //zero field should not validate
	err = (&RangeValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	err = (&RangeValidator{MinTagName:"customMin", MaxTagName:"customMax"}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.A = 21
	err = (&RangeValidator{MinTagName:"customMin", MaxTagName:"customMax"}).Validate(s)
	if err == nil {
		t.Fatal("s.E.A range should be checked.")
	}
}


func TestLengthValidators(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:10,
		D:"aaa",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			C:"xyz",
			D:"bbb",
			E:80,
		},
	}
	
	s.E.B = ""
	err := (&LengthValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.B = "1234567890"
	err = (&LengthValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.B = "12345678901234567890"
	err = (&LengthValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.B = "123456789"
	err = (&LengthValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.B length should be checked.")
	}
	
	s.E.B = "123456789012345678901"
	err = (&LengthValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.B length should be checked.")
	}
}

func TestStringValidators(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:10,
		D:"abc",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			C:"xyz",
			D:"bbb",
			E:80,
		},
	}
	
	s.E.C = ""
	err := (&StringValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.C = "192.168.0.1"
	err = (&StringValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.C = "255.255.255.0"
	err = (&StringValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.C = "123456789"
	err = (&StringValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.C regex should be checked.")
	}
	
	s.E.C = "192.168.0.1"
	s.D = "aaa"
	err = (&StringValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.D regex should be checked.")
	}
}


func TestSelectionValidators(t *testing.T) {
	s := &TestValidatorStruct{
		A:10,
		B:30,
		C:10,
		D:"abc",
		E: TestValidatorInnerStruct{
			A:10,
			B:"abc",
			C:"xyz",
			D:"bbb",
			E:80,
		},
	}
	
	s.B = 0
	s.E.D = ""
	err := (&SelectionValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.D = "tcp"
	err = (&SelectionValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.D = "udp"
	err = (&SelectionValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
	s.E.D = "UDP"
	err = (&SelectionValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.E.D selection should be checked.")
	}
	
	s.E.D = "tcp"
	s.B = 30
	err = (&SelectionValidator{}).Validate(s)
	if err == nil {
		t.Fatal("s.B regex should be checked.")
	}
	
	s.B=20
	err = (&SelectionValidator{}).Validate(s)
	if err != nil {
		t.Fatal(err)
	}
	
}