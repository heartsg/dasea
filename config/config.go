package config

import (
	"fmt"
	"os"
	"strings"
)

// 1. Overview
//
// dasea/config is based on multiconfig.
//
// It supports to load configuration/options from flags, environment variables,
// json/toml/ini files. Environment variables settings will overwite settings from
// files, and flags will overwrite settings from the previous two settings.
// 
// It also supports to validate the loaded variables. The validation currently supported
// includes
//  - default: will assign default values to the variable if it is not initialized
//  - required: whether the variable must be initialized with non-zero values
//  - min/max: the minimum/maximum allowed range for integer/floating points
//  - minlen/maxlen: the minimum/maximum allowed length for strings
//                 : (future) support min/max allowed length for slices/maps?
//  - regex: the allowed regular expression for the input string, this is useful
//         : for settings such as IP addresses, phone numbers, and etc.
//  - select: the allowed sets of strings for the variable. This can be done with
//          : regex too, but we decide to move it out as a special validation rule.
//
// 2. Tags for structure
// 
// dasea/config load the options directly into structure field. Tags can be defined to 
// customize how each field is loaded or validated.
//
// Available tags are:
//
// Options related 
//  - flag:"name"
//  - env:"name"
//  - 
// Validation related
//  - default:"value"
//  - required:"true/false"
//  - min/max/minlen/maxlen:"value"
//  - regex:"string"
//  - select:"option1|option2|option3"
//
// 3. How options are loaded into structure field
//
// - For flag, if flatten is true, a "--name" or "-name" flag will directly map to a inner field called 
//   "Name". If flatten is false, a "--child-grandchild-grandgrandchild" will map to a innerfiled called
//   "Grandgrandchild", whose parent is "Grandchild", whose parent is "Child", whose parent is the structure
//   object itself.
//
//   If camelcase is enabled, structure field with multiple capatal letters will be mapped to "-", e.g.,
//   "--is-enabled" will be mapped to strcture field IsEnabled.
//
//   If there is flag:"name" for a structure field named "Field", "--name" will map to "Field".
//
// - For env, a "NAME" environment variable will map to "Name" field of a structure. a "CHILD_GRANDCHILD" will
//   map to a field name called "Grandchild" whose parent is "Child" and etc.
//
//   If camelcase is enabled, structure field with multiple capatal letters will be mapped to "_", e.g., 
//   "IS_ENABLED" will be mapped to structure field "IsEnabled".
//
// - For JSON file, the outer field in Json must be array of maps, the key:value pairs match to the 
//   outer structure fields. The key:value pair where the value is another array of maps, will be 
//   mapped to the innter structure fields.
//
// - For TOML file, the standard toml format follows.
//
// - For INI file, if LoadFromGlobal is enabled, the global key:value pairs will be mapped to outer
//   structure fields, the inner structure fields will be from [Child], [Child.Grandchild] etc.
//
//   If LoadFromGlobal is disabled, the outer structure fields must be under [OutStructureName] section.
//   The inner structure fields will be under [OutStructureName.Child] [OutStructureName.Child.Grandchild].
//
// 4. Usage
//
// Firstly, get the loader. If you have config file name and type, you can just call
//   loader := NewWithPath(fileName, type)
// If you do not have config file, and want to rely on --config-file flag or CONFIG_FILE env variable, use
//   loader := New()
// The default config file will be "config.toml"
// After loader is created, you can load options into struct object,
//   object := new(Structure)
//   err := loader.Load(object)
// Lastly, you can validate
//   err := loader.Validate(object)



// Loader loads the configuration from a source. The implementer of Loader is
// responsible of setting the default values of the struct.
type Loader interface {
	// Load loads the source into the config defined by struct s
	Load(s interface{}) error
}

// DefaultLoader implements the Loader interface. It initializes the given
// pointer of struct s with configuration from the default sources. The order
// of load is TagLoader, FileLoader, EnvLoader and lastly FlagLoader. An error
// in any step stops the loading process. Each step overrides the previous
// step's config (i.e: defining a flag will override previous environment or
// file config). To customize the order use the individual load functions.
type DefaultLoader struct {
	Loader
	Validator
}

// NewWithPath returns a new instance of Loader to read from the given
// configuration file.
func NewWithPath(path string, confType string) *DefaultLoader {
	loaders := []Loader{}

	// Read default values defined via tag fields "default"
	loaders = append(loaders, &TagLoader{})

	// Choose configuration file format
	if confType == "toml" || (confType == "auto" && strings.HasSuffix(path, "toml")) {
		loaders = append(loaders, &TOMLLoader{Path: path})
	}

	if confType == "json" || (confType == "auto" && strings.HasSuffix(path, "json")) {
		loaders = append(loaders, &JSONLoader{Path: path})
	}
	
	if confType == "ini" || (confType == "auto" && strings.HasSuffix(path, "ini")) {
		loaders = append(loaders, &INILoader{Path: path})
	}

	e := &EnvironmentLoader{}
	f := &FlagLoader{}

	loaders = append(loaders, e, f)
	loader := MultiLoader(loaders...)

	d := &DefaultLoader{}
	d.Loader = loader
	d.Validator = MultiValidator(&RequiredValidator{}, &LengthValidator{}, &RangeValidator{}, &StringValidator{}, &SelectionValidator{})
	return d
}

// New returns a new instance of DefaultLoader without any file loaders.
func New() *DefaultLoader {
	configOptsLoader := MultiLoader(
		&TagLoader{},
		&EnvironmentLoader{},
		&FlagLoader{},
	)
	configOpts := new(ConfigOpts)
	if err := configOptsLoader.Load(configOpts); err == nil {
		if _, err := os.Stat(configOpts.File); err == nil {
			return NewWithPath(configOpts.File, "auto")
		}
	}
	
	//some errors for config files, ignore it, and build
	//a loader without config file
	loader := MultiLoader(
		&TagLoader{},
		&EnvironmentLoader{},
		&FlagLoader{},
	)

	d := &DefaultLoader{}
	d.Loader = loader
	d.Validator = MultiValidator(&RequiredValidator{}, &LengthValidator{}, &RangeValidator{}, &StringValidator{}, &SelectionValidator{})
	return d
}

// MustLoadWithPath loads with the DefaultLoader settings and from the given
// Path. It exits if the config cannot be parsed.
func MustLoadWithPath(path string, confType string, conf interface{}) {
	d := NewWithPath(path, confType)
	d.MustLoad(conf)
}

// MustLoad loads with the DefaultLoader settings. It exits if the config
// cannot be parsed.
func MustLoad(conf interface{}) {
	d := New()
	d.MustLoad(conf)
}

// MustLoad is like Load but panics if the config cannot be parsed.
func (d *DefaultLoader) MustLoad(conf interface{}) {
	if err := d.Load(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// we at koding, believe having sane defaults in our system, this is the
	// reason why we have default validators in DefaultLoader. But do not cause
	// nil pointer panics if one uses DefaultLoader directly.
	if d.Validator != nil {
		d.MustValidate(conf)
	}
}

// MustValidate validates the struct. It exits with status 1 if it can't
// validate.
func (d *DefaultLoader) MustValidate(conf interface{}) {
	if err := d.Validate(conf); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}