package config

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"time"
	"regexp"
	"encoding/json"
	"github.com/fatih/structs"
)

// fieldSet sets field value from the given string value. It converts the
// string value in a sane way and is usefulf or environment variables or flags
// which are by nature in string types.
func fieldSet(field *structs.Field, v string) error {
	switch field.Kind() {
	case reflect.Bool:
		val, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}

		if err := field.Set(val); err != nil {
			return err
		}
	case reflect.Int:
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}

		if err := field.Set(i); err != nil {
			return err
		}
	case reflect.String:
		if err := field.Set(v); err != nil {
			return err
		}
	case reflect.Map:
		switch t := field.Value().(type) {
		case map[string]int:
			si := make(map[string]int)
			if err := json.Unmarshal([]byte(v), &si); err != nil {
				return err
			}
			if err := field.Set(si); err != nil {
				return err
			}
		case map[string]string:
			ss := make(map[string]string)
			if err := json.Unmarshal([]byte(v), &ss); err != nil {
				return err
			}
			if err := field.Set(ss); err != nil {
				return err
			}
		default:
			return fmt.Errorf("config: field '%s' of type map is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}
	case reflect.Slice:
		switch t := field.Value().(type) {
		case []string:
			if err := field.Set(strings.Split(v, ",")); err != nil {
				return err
			}
		case []int:
			var list []int
			for _, in := range strings.Split(v, ",") {
				i, err := strconv.Atoi(in)
				if err != nil {
					return err
				}

				list = append(list, i)
			}

			if err := field.Set(list); err != nil {
				return err
			}
		case []int64:
			var list []int64
			for _, in := range strings.Split(v, ",") {
				i, err := strconv.ParseInt(in, 10, 64)
				if err != nil {
					return err
				}

				list = append(list, i)
			}

			if err := field.Set(list); err != nil {
				return err
			}
		case []float64:
			var list []float64
			for _, in := range strings.Split(v, ",") {
				i, err := strconv.ParseFloat(in, 64)
				if err != nil {
					return err
				}

				list = append(list, i)
			}

			if err := field.Set(list); err != nil {
				return err
			}			
		default:
			return fmt.Errorf("config: field '%s' of type slice is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}
	case reflect.Float64:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}

		if err := field.Set(f); err != nil {
			return err
		}
	case reflect.Int64:
		switch t := field.Value().(type) {
		case time.Duration:
			d, err := time.ParseDuration(v)
			if err != nil {
				return err
			}

			if err := field.Set(d); err != nil {
				return err
			}
		case int64:
			p, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return err
			}

			if err := field.Set(p); err != nil {
				return err
			}
		default:
			return fmt.Errorf("config: field '%s' of type int64 is unsupported: %s (%T)",
				field.Name(), field.Kind(), t)
		}
	default:
		return fmt.Errorf("config: field '%s' has unsupported type: %s", field.Name(), field.Kind())
	}

	return nil
}

func validateRange(field *structs.Field, minStr string, maxStr string) error {
	switch field.Kind() {
	case reflect.Int:
		i := field.Value().(int)
		
		min, minErr := strconv.Atoi(minStr)
		max, maxErr := strconv.Atoi(maxStr)
		if (minErr == nil && i < min) || (maxErr == nil && i > max) {
			return fmt.Errorf("validate: Field %s value %d, want: [%s, %s]", field.Name(), i, minStr, maxStr)
		}
	case reflect.Float64:
		i := field.Value().(float64)
		
		min, minErr := strconv.ParseFloat(minStr, 64)
		max, maxErr := strconv.ParseFloat(maxStr, 64)
		if (minErr == nil && i < min) || (maxErr == nil && i > max) {
			return fmt.Errorf("validate: Field %s value %f, want: [%s, %s]", field.Name(), i, minStr, maxStr)
		}
	case reflect.Int64:
		switch field.Value().(type) {
		case int64:
			i := field.Value().(int64)
			
			min, minErr := strconv.ParseInt(minStr, 10, 64)
			max, maxErr := strconv.ParseInt(maxStr, 10, 64)
			if (minErr == nil && i < min) || (maxErr == nil && i > max) {
				return fmt.Errorf("validate: Field %s value %d, want: [%s, %s]", field.Name(), i, minStr, maxStr)
			}
		default:
			return nil
		}
	default:
		return nil
	}
	
	return nil
}



func validateLength(field *structs.Field, minlenStr string, maxlenStr string) error {
	switch field.Kind() {
	case reflect.String:
		i := field.Value().(string)
		l := len(i)
		
		min, minErr := strconv.Atoi(minlenStr)
		max, maxErr := strconv.Atoi(maxlenStr)
		if (minErr == nil && l < min) || (maxErr == nil && l > max) {
			return fmt.Errorf("validate: Field %s value %s, length limit: [%s, %s]", field.Name(), i, minlenStr, maxlenStr)
		}
	default:
		return nil
	}
	
	return nil
}

func validateString(field *structs.Field, regex string) error {
	switch field.Kind() {
	case reflect.String:
		i := field.Value().(string)
		matched, err := regexp.MatchString(regex, i)
		if err == nil && !matched {
			return fmt.Errorf("validate: Field %s value %s, regex %s not matched.", field.Name(), i, regex)
		}
	}
	
	return nil
}


func validateSelection(field *structs.Field, selection string) error {
	selectionSet := strings.Split(selection, "|")
	switch field.Kind() {
		case reflect.Int:
			i := field.Value().(int)
			s := strconv.Itoa(i)
			for _, selection := range selectionSet {
				if s == selection {
					return nil
				}
			}
			return fmt.Errorf("validate: Field %s value %s, not in selection %s.", field.Name(), s, selection)
		case reflect.String:
			s := field.Value().(string)
			for _, selection := range selectionSet {
				if s == selection {
					return nil
				}
			}
			return fmt.Errorf("validate: Field %s value %s, not in selection %s.", field.Name(), s, selection)
		default:
			return nil
	}
	return nil
}
