package config

import (
	"reflect"
	"github.com/fatih/structs"
)

//This library requires that the tag must be in `key:"value" key:"value"` format, otherwise
// keys won't be recognized (we use structs/reflect.field.tag.Get()), may change later
// for better parsing


// TagLoader satisfies the loader interface. It parses a struct's field tags
// and populated the each field with that given tag.
type TagLoader struct {
	// DefaultTagName is the default tag name for struct fields to define
	// default values for a field. Example:
	//
	//   // Field's default value is "koding".
	//   Name string `default:"koding"`
	//
	// The default value is "default" if it's not set explicitly.
	DefaultTagName string
}

func (t *TagLoader) Load(s interface{}) error {
	if t.DefaultTagName == "" {
		t.DefaultTagName = "default"
	}

	for _, field := range structs.Fields(s) {
		if err := t.processFieldDefaultValue(field); err != nil {
			return err
		}
	}

	return nil
}

// processField gets tagName and the field, recursively checks if the field has the given
// tag, if yes, sets it otherwise ignores
func (t *TagLoader) processFieldDefaultValue(field *structs.Field) error {
	
	//first process each subfield of a struct field
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			if err := t.processFieldDefaultValue(f); err != nil {
				return err
			}
		}
	default:
		//Set default value for the field itself, including struct field
		// If there's a default value tag for struct field, it should be in Json format
		defaultVal := field.Tag(t.DefaultTagName)
		if defaultVal == "" {
			return nil
		}
	
		err := fieldSet(field, defaultVal)
		if err != nil {
			return err
		}
	}
	


	return nil
}