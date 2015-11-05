package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/fatih/structs"
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/goconfig"
)

var (
	// ErrSourceNotSet states that neither the path or the reader is set on the loader
	ErrSourceNotSet = errors.New("config path or reader is not set")

	// ErrFileNotFound states that given file is not exists
	ErrFileNotFound = errors.New("config file not found")
)


// TOMLLoader satisifies the loader interface. It loads the configuration from
// the given toml file or Reader
type TOMLLoader struct {
	Path   string
	Reader io.Reader
}

// Load loads the source into the config defined by struct s
// Defaults to using the Reader if provided, otherwise tries to read from the
// file
func (t *TOMLLoader) Load(s interface{}) error {
	var r io.Reader

	if t.Reader != nil {
		r = t.Reader
	} else if t.Path != "" {
		file, err := getConfig(t.Path)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	} else {
		return ErrSourceNotSet
	}

	if _, err := toml.DecodeReader(r, s); err != nil {
		return err
	}

	return nil
}

// JSONLoader satisifies the loader interface. It loads the configuration from
// the given json file or Reader
type JSONLoader struct {
	Path   string
	Reader io.Reader
}

// Load loads the source into the config defined by struct s
// Defaults to using the Reader if provided, otherwise tries to read from the
// file
func (j *JSONLoader) Load(s interface{}) error {
	var r io.Reader
	if j.Reader != nil {
		r = j.Reader
	} else if j.Path != "" {
		file, err := getConfig(j.Path)
		if err != nil {
			return err
		}
		defer file.Close()
		r = file
	} else {
		return ErrSourceNotSet
	}

	return json.NewDecoder(r).Decode(s)
}

func getConfig(path string) (*os.File, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configPath := path
	if !filepath.IsAbs(path) {
		configPath = filepath.Join(pwd, path)
	}

	// check if file with combined path is exists(relative path)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		return os.Open(configPath)
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	return f, err
}


// INILoader satisifies the loader interface. It loads the configuration from
// the given ini file or Reader
type INILoader struct {
	Path   string
	//Due to limitation of API of package goconfig, we do not provide this option
	// for INIReader
	//Reader io.Reader
	
	// Load from global? or from a section with name structname?
	// Default false, load from section with section name equal to structname
	LoadFromGlobal bool
}

// Load loads the source into the config defined by struct s
// Defaults to using the Reader if provided, otherwise tries to read from the
// file
func (i *INILoader) Load(s interface{}) (err error) {
	var c *goconfig.ConfigFile
	if i.Path != "" {
		c, err = goconfig.LoadConfigFile(i.Path)
		if err != nil {
			return err
		}
	} else {
		return ErrSourceNotSet
	}
	
	//Now parse struct fields from ini file
	//Since INI file doesn't natually corresponds to a golang struct,
	// most current INI parsers do not provide Json or Toml functions that
	// directly write into the struct fields, we have to parse by ourselves.
	
	//We define the rule as follows
	// If load from global, then [section] means sub-struct
	// else [section] means the struct itself
	
	// Substructs are represented by dot separated sections
	//   [parent.child.child] and etc.
	
	// Slices are separated by comma , such as
	//   array = 1,2,3
	// Maps are defined by Json style
	//   map = {"key":1, "key":2}
	
	var section string
	if i.LoadFromGlobal {
		section = ""
	} else {
		section = structs.New(s).Name()
	}
	
	for _, field := range structs.Fields(s) {
		if err := i.setField(field, c, section); err != nil {
			return err
		}
	}

	return nil
}

func (i *INILoader) setField(field *structs.Field, c *goconfig.ConfigFile, section string) error {	
	//first process each subfield of a struct field
	switch field.Kind() {
	case reflect.Struct:
		for _, f := range field.Fields() {
			var subsection string
			if section == "" {
				subsection = field.Name()
			} else {
				subsection = section + "." + field.Name()
			}
			if err := i.setField(f, c, subsection); err != nil {
				return err
			}
		}
	default:
		 v, err := c.GetValue(section, field.Name())
		 if err == nil && v != "" {
			err := fieldSet(field, v)
			if err != nil {
				return err
			}
		 } 
	}

	return nil
}