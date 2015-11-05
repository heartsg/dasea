package config

import (
	"testing"
	
	"github.com/fatih/structs"
)

type TestReflectStruct struct {
	A int
	B float64
	C []int
	D map[string]int
	E map[string]string
}

func TestReflect(t *testing.T) {
	test := &TestReflectStruct{}
	var err error
	for _, field := range structs.Fields(test) {
		switch field.Name() {
		case "A":
			err = fieldSet(field, "1")
			if err != nil {
				t.Error(err)
			}
			if test.A != 1 {
				t.Errorf("A: %d, wanted: %d", test.A, 1)
			}
		case "B":
			err = fieldSet(field, "3.14")
			if err != nil {
				t.Error(err)
			}
			if test.B != 3.14 {
				t.Errorf("B: %f, wanted: %f", test.B, 3.14)
			}
		case "C":
			err = fieldSet(field, "1,2,3")
			if err != nil {
				t.Error(err)
			}
			if test.C[2] != 3 {
				t.Errorf("C: %s, wanted: %s", test.C, []int{1,2,3})
			}
		case "D":
			err = fieldSet(field, "{\"abc\":1,\"def\":2}")
			if err != nil {
				t.Error(err)
			}
			if test.D["abc"] != 1 {
				t.Errorf("D: %s, wanted: %s", test.D, map[string]int{"abc":1, "def":2})
			}
		case "E":
			err = fieldSet(field, `{"abc":"123","def":"456"}`)
			if err != nil {
				t.Error(err)
			}
			if test.E["abc"] != "123" || test.E["def"] != "456" {
				t.Errorf("E: %s, wanted: %s", test.E, map[string]string{"abc":"123", "def":"456"})
			}
		}
	}
}

type TestRangeStruct struct {
	A int
	B int64
	C float64
	D string
}

func TestValidateRange(t *testing.T) {
	test := &TestRangeStruct{10, 20, 30.4, "xyz"}
	var err error
	for _, field := range structs.Fields(test) {
		switch field.Name() {
		case "A":
			err = validateRange(field, "-1", "30")
			if err != nil {
				t.Error(err)
			}
			err = validateRange(field, "", "20")
			if err != nil {
				t.Error(err)
			}
			err = validateRange(field, "11", "")
			if err == nil {
				t.Error("validate: Exceeds range but not checked")
			}
			
			err = validateRange(field, "", "9")
			if err == nil {
				t.Error("validate: Exceeds range but not checked")
			}
		case "B":
			err = validateRange(field, "-1", "30")
			if err != nil {
				t.Error(err)
			}
			err = validateRange(field, "", "20")
			if err != nil {
				t.Error(err)
			}
			err = validateRange(field, "11", "")
			if err != nil {
				t.Error(err)
			}
			
			err = validateRange(field, "", "19")
			if err == nil {
				t.Error("validate: Exceeds range but not checked")
			}
		case "C":
			err = validateRange(field, "29.0", "30.5")
			if err != nil {
				t.Error(err)
			}
			err = validateRange(field, "", "30.3")
			if err == nil {
				t.Error("validate: Exceeds range but not checked")
			}
		case "D":
			err = validateRange(field, "1", "2")
			if err != nil {
				t.Error(err)
			}
		}
	}

}


type TestLengthStruct struct {
	A string
}

func TestValidateLength(t *testing.T) {
	test := &TestLengthStruct{"xyz"}
	var err error
	for _, field := range structs.Fields(test) {
		switch field.Name() {
		case "A":
			err = validateLength(field, "1", "4")
			if err != nil {
				t.Error(err)
			}
			err = validateLength(field, "5", "6")
			if err == nil {
				t.Error("validate: Exceeds range but not checked")
			}
		}
	}

}


type TestStringStruct struct {
	A string
}

func TestValidateString(t *testing.T) {
	test := &TestStringStruct{"seafood"}
	var err error
	for _, field := range structs.Fields(test) {
		switch field.Name() {
		case "A":
			err = validateString(field, "foo.*")
			if err != nil {
				t.Error(err)
			}
			err = validateString(field, "bar.*")
			if err == nil {
				t.Error("validate: Shouldn't match but matched")
			}
		}
	}
}

type TestSelectionStruct struct {
	A int
	B string
}

func TestValidateSelection(t *testing.T) {
	test := &TestSelectionStruct{A:100, B:"fish"}
	var err error
	for _, field := range structs.Fields(test) {
		switch field.Name() {
		case "A":
			err = validateSelection(field, "50|100|150")
			if err != nil {
				t.Error(err)
			}
			err = validateSelection(field, "0|25|50")
			if err == nil {
				t.Error("validate: Shouldn't match but matched")
			}
		case "B":
			err = validateSelection(field, "fish|cat|dog")
			if err != nil {
				t.Error(err)
			}
			err = validateSelection(field, "apple|peach|orange")
			if err == nil {
				t.Error("validate: Shouldn't match but matched")
			}
		}
	}
}