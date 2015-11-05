package config

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

// Validator validates the config against any predefined rules, those predefined
// rules should be given to this package. The implementer will be responsible
// about the logic
type Validator interface {
	// Validate validates the config struct
	Validate(s interface{}) error
}

// RequiredValidator validates the struct against zero values
type RequiredValidator struct {
	// RequiredTagName is the default tag name for struct fields to define a field that has to be
	// initialized (non zero). The default value is "required" if it's not set explicitly.
	// Example:
	//  `required:"true"`
	//  `required:"1"`
	//  `required:"false"`
	//  values other than true and 1 will be considered as false
	RequiredTagName string
}

// RangeValidator validates allowed range of struct fields
type RangeValidator struct {
	// MinTagName is the default tag name for struct fields to define minimum values
	// for an integer/float field. The default value is "min" if it's not set explicitly.
	MinTagName string
	
	// MaxTagName is the default tag name for struct fields to define maximum values
	// for an integer/float field. The default value is "max" if it's not set explicitly.
	MaxTagName string
}

// LengthValidator validates allowed length of string/slice
type LengthValidator struct {
	// MinlenTagName is the default tag name for struct fields to define minimum length
	// for a string field. The default value is "minlen" if it's not set explicitly.
	MinlenTagName string
	
	// MaxlenTagName is the default tag name for struct fields to define maximum length
	// for a string field. The default value is "maxlen" if it's not set explicitly.
	MaxlenTagName string
}

// StringValidator validates allowed string fields
type StringValidator struct {	
	// RegexTagName is the default tag name for struct fields to define regular expression
	// for a string field. The default value is "regex" if it's not set explicitly.
	RegexTagName string
}


// SelectionValidator validates fields with allowed values
type SelectionValidator struct {
	// SelectTagName is the default tag name for struct fields to define a set of allowed values.
	// The default value is "select" if it's not set explicitly.
	SelectTagName string
}



// Validate validates the given struct agaist field's zero values. If
// intentionaly, the value of a field is `zero-valued`(e.g false, 0, "")
// required tag should not be set for that field.
func (e *RequiredValidator) Validate(s interface{}) error {
	if e.RequiredTagName == "" {
		e.RequiredTagName = "required"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *RequiredValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		// this is used for error messages below, when we have an error at the
		// child properties add parent properties into the error message as well
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		val := field.Tag(e.RequiredTagName)
		if val != "true" && val != "1" {
			return nil
		}
		if field.IsZero() {
			return fmt.Errorf("Field '%s' is required", fieldName)
		}
	}

	return nil
}




func (e *RangeValidator) Validate(s interface{}) error {
	if e.MinTagName == "" {
		e.MinTagName = "min"
	}
	
	if e.MaxTagName == "" {
		e.MaxTagName = "max"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *RangeValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		// this is used for error messages below, when we have an error at the
		// child properties add parent properties into the error message as well
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		if field.IsZero() { //if not initialized, we won't validate it
			return nil
		}
		minStr := field.Tag(e.MinTagName)
		maxStr := field.Tag(e.MaxTagName)
		if minStr == "" && maxStr == "" {
			return nil
		}
		
		
		err := validateRange(field, minStr, maxStr)
		if err != nil {
			return err
		}
	}

	return nil
}



func (e *LengthValidator) Validate(s interface{}) error {
	if e.MaxlenTagName == "" {
		e.MaxlenTagName = "maxlen"
	}
	
	if e.MinlenTagName == "" {
		e.MinlenTagName = "minlen"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *LengthValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		// this is used for error messages below, when we have an error at the
		// child properties add parent properties into the error message as well
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		if field.IsZero() {
			return nil
		}
		minlenStr := field.Tag(e.MinlenTagName)
		maxlenStr := field.Tag(e.MaxlenTagName)
		if minlenStr == "" && maxlenStr == "" {
			return nil
		}

		err := validateLength(field, minlenStr, maxlenStr)
		if err != nil {
			return err
		}
	}

	return nil
}


func (e *StringValidator) Validate(s interface{}) error {
	if e.RegexTagName == "" {
		e.RegexTagName = "regex"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *StringValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		// this is used for error messages below, when we have an error at the
		// child properties add parent properties into the error message as well
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		if field.IsZero() {
			return nil
		}
		regexStr := field.Tag(e.RegexTagName)
		if regexStr == "" {
			return nil
		}
		
		err := validateString(field, regexStr)
		if err != nil {
			return err
		}
	}

	return nil
}


func (e *SelectionValidator) Validate(s interface{}) error {
	if e.SelectTagName == "" {
		e.SelectTagName = "select"
	}

	for _, field := range structs.Fields(s) {
		if err := e.processField("", field); err != nil {
			return err
		}
	}

	return nil
}

func (e *SelectionValidator) processField(fieldName string, field *structs.Field) error {
	fieldName += field.Name()
	switch field.Kind() {
	case reflect.Struct:
		// this is used for error messages below, when we have an error at the
		// child properties add parent properties into the error message as well
		fieldName += "."

		for _, f := range field.Fields() {
			if err := e.processField(fieldName, f); err != nil {
				return err
			}
		}
	default:
		if field.IsZero() {
			return nil
		}
		selectionStr := field.Tag(e.SelectTagName)
		if selectionStr == "" {
			return nil
		}
		
		err := validateSelection(field, selectionStr)
		if err != nil {
			return err
		}
	}

	return nil
}



