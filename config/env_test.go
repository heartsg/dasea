package config

import (
	"os"
	"strings"
	"testing"
	"time"
	
	"github.com/fatih/structs"
)


type TestEnvServer struct {
	Name       string
	Port       int
	ID         int64
	Labels     []int
	Enabled    bool
	Users      []string
	Postgres   TestEnvPostgres
	unexported string
	Interval   time.Duration
}

type TestEnvPostgres struct {
	Enabled           bool
	Port              int    
	Hosts             []string 
	DBName            string
	AvailabilityRatio float64
	unexported        string
}
type TestEnvCamelCaseServer struct {
	AccessKey         string
	Normal            string
	DBName            string
	AvailabilityRatio float64
}
type TestEnvTagServer struct {
	AccessKey         string `env:"ACC"`
	Normal            string `env:"NOR"`
	DBName            string `env:"DB"`
	AvailabilityRatio float64 `env:"AR"`
}


func getDefaultEnvServer() *TestEnvServer {
	return &TestEnvServer{
		Name:     "koding",
		Port:     6060,
		Enabled:  true,
		ID:       1234567890,
		Labels:   []int{123, 456},
		Users:    []string{"ankara", "istanbul"},
		Interval: 10 * time.Second,
		Postgres: TestEnvPostgres{
			Enabled:           true,
			Port:              5432,
			Hosts:             []string{"192.168.2.1", "192.168.2.2", "192.168.2.3"},
			DBName:            "configdb",
			AvailabilityRatio: 8.23,
		},
	}
}

func getDefaultEnvCamelCaseServer() *TestEnvCamelCaseServer {
	return &TestEnvCamelCaseServer{
		AccessKey:         "123456",
		Normal:            "normal",
		DBName:            "configdb",
		AvailabilityRatio: 8.23,
	}
}
func getDefaultEnvTagServer() *TestEnvTagServer {
	return &TestEnvTagServer{
		AccessKey:         "123456",
		Normal:            "normal",
		DBName:            "configdb",
		AvailabilityRatio: 8.23,
	}
}

func testEnvServerStruct(t *testing.T, s *TestEnvServer, d *TestEnvServer) {
	if s.Name != d.Name {
		t.Errorf("Name value is wrong: %s, want: %s", s.Name, d.Name)
	}

	if s.Port != d.Port {
		t.Errorf("Port value is wrong: %d, want: %d", s.Port, d.Port)
	}

	if s.Enabled != d.Enabled {
		t.Errorf("Enabled value is wrong: %t, want: %t", s.Enabled, d.Enabled)
	}

	if s.Interval != d.Interval {
		t.Errorf("Interval value is wrong: %v, want: %v", s.Interval, d.Interval)
	}

	if s.ID != d.ID {
		t.Errorf("ID value is wrong: %v, want: %v", s.ID, d.ID)
	}

	if len(s.Labels) != len(d.Labels) {
		t.Errorf("Labels value is wrong: %d, want: %d", len(s.Labels), len(d.Labels))
	} else {
		for i, label := range d.Labels {
			if s.Labels[i] != label {
				t.Errorf("Label is wrong for index: %d, label: %d, want: %d", i, s.Labels[i], label)
			}
		}
	}

	if len(s.Users) != len(d.Users) {
		t.Errorf("Users value is wrong: %d, want: %d", len(s.Users), len(d.Users))
	} else {
		for i, user := range d.Users {
			if s.Users[i] != user {
				t.Errorf("User is wrong for index: %d, user: %s, want: %s", i, s.Users[i], user)
			}
		}
	}

	// Explicitly state that Enabled should be true, no need to check
	// x == true infact.
	if s.Postgres.Enabled != d.Postgres.Enabled {
		t.Errorf("Postgres enabled is wrong %t, want: %t", s.Postgres.Enabled, d.Postgres.Enabled)
	}

	if s.Postgres.Port != d.Postgres.Port {
		t.Errorf("Postgres Port value is wrong: %d, want: %d", s.Postgres.Port, d.Postgres.Port)
	}

	if s.Postgres.DBName != d.Postgres.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.Postgres.DBName, d.Postgres.DBName)
	}

	if s.Postgres.AvailabilityRatio != d.Postgres.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.Postgres.AvailabilityRatio, d.Postgres.AvailabilityRatio)
	}

	if len(s.Postgres.Hosts) != len(d.Postgres.Hosts) {
		// do not continue testing if this fails, because others is depending on this test
		t.Fatalf("Hosts len is wrong: %v, want: %v", s.Postgres.Hosts, d.Postgres.Hosts)
	}

	for i, host := range d.Postgres.Hosts {
		if s.Postgres.Hosts[i] != host {
			t.Fatalf("Hosts number %d is wrong: %v, want: %v", i, s.Postgres.Hosts[i], host)
		}
	}
}

func testEnvCamelcaseServerStruct(t *testing.T, s *TestEnvCamelCaseServer, d *TestEnvCamelCaseServer) {
	if s.AccessKey != d.AccessKey {
		t.Errorf("AccessKey is wrong: %s, want: %s", s.AccessKey, d.AccessKey)
	}

	if s.Normal != d.Normal {
		t.Errorf("Normal is wrong: %s, want: %s", s.Normal, d.Normal)
	}

	if s.DBName != d.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.DBName, d.DBName)
	}

	if s.AvailabilityRatio != d.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.AvailabilityRatio, d.AvailabilityRatio)
	}

}


func testEnvTagServerStruct(t *testing.T, s *TestEnvTagServer, d *TestEnvTagServer) {
	if s.AccessKey != d.AccessKey {
		t.Errorf("AccessKey is wrong: %s, want: %s", s.AccessKey, d.AccessKey)
	}

	if s.Normal != d.Normal {
		t.Errorf("Normal is wrong: %s, want: %s", s.Normal, d.Normal)
	}

	if s.DBName != d.DBName {
		t.Errorf("DBName is wrong: %s, want: %s", s.DBName, d.DBName)
	}

	if s.AvailabilityRatio != d.AvailabilityRatio {
		t.Errorf("AvailabilityRatio is wrong: %f, want: %f", s.AvailabilityRatio, d.AvailabilityRatio)
	}

}


func TestENV(t *testing.T) {
	m := EnvironmentLoader{}
	s := &TestEnvServer{}
	structName := structs.Name(s)

	// set env variables
	setEnvVars(t, structName, "")

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEnvServerStruct(t, s, getDefaultEnvServer())
}

func TestCamelCaseEnv(t *testing.T) {
	m := EnvironmentLoader{
		CamelCase: true,
	}
	s := &TestEnvCamelCaseServer{}
	structName := structs.Name(s)

	// set env variables
	setEnvVars(t, structName, "")

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEnvCamelcaseServerStruct(t, s, getDefaultEnvCamelCaseServer())
}

func TestTagEnv(t *testing.T) {
	m := EnvironmentLoader{
		EnvTagName: "env",
	}
	s := &TestEnvTagServer{}
	structName := structs.Name(s)

	// set env variables
	setEnvVars(t, structName, "")

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEnvTagServerStruct(t, s, getDefaultEnvTagServer())
}

func TestENVWithPrefix(t *testing.T) {
	const prefix = "Prefix"

	m := EnvironmentLoader{Prefix: prefix}
	s := &TestEnvServer{}
	structName := structs.New(s).Name()

	// set env variables
	setEnvVars(t, structName, prefix)

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testEnvServerStruct(t, s, getDefaultEnvServer())
}

func setEnvVars(t *testing.T, structName, prefix string) {
	if structName == "" {
		t.Fatal("struct name can not be empty")
	}

	var env map[string]string
	switch structName {
	case "TestEnvServer":
		env = map[string]string{
			"NAME":                       "koding",
			"PORT":                       "6060",
			"ENABLED":                    "true",
			"USERS":                      "ankara,istanbul",
			"INTERVAL":                   "10s",
			"ID":                         "1234567890",
			"LABELS":                     "123,456",
			"POSTGRES_ENABLED":           "true",
			"POSTGRES_PORT":              "5432",
			"POSTGRES_HOSTS":             "192.168.2.1,192.168.2.2,192.168.2.3",
			"POSTGRES_DBNAME":            "configdb",
			"POSTGRES_AVAILABILITYRATIO": "8.23",
			"POSTGRES_FOO":               "8.23,9.12,11,90",
		}
	case "TestEnvCamelCaseServer":
		env = map[string]string{
			"ACCESS_KEY":         "123456",
			"NORMAL":             "normal",
			"DB_NAME":            "configdb",
			"AVAILABILITY_RATIO": "8.23",
		}
	case "TestEnvTagServer":
		env = map[string]string{
			"ACC":         "123456",
			"NOR":         "normal",
			"DB":          "configdb",
			"AR":          "8.23",
		}
	}

	if prefix == "" {
		prefix = structName
	}

	prefix = strings.ToUpper(prefix)

	for key, val := range env {
		var env string
		if structName != "TestEnvTagServer" {
			env = prefix + "_" + key
		} else {
			env = key
		}
		if err := os.Setenv(env, val); err != nil {
			t.Fatal(err)
		}
	}
}

func TestENVgetPrefix(t *testing.T) {
	e := &EnvironmentLoader{}
	s := &TestEnvServer{}

	st := structs.New(s)

	prefix := st.Name()

	if p := e.getPrefix(st); p != prefix {
		t.Errorf("Prefix is wrong: %s, want: %s", p, prefix)
	}

	prefix = "Test"
	e = &EnvironmentLoader{Prefix: prefix}
	if p := e.getPrefix(st); p != prefix {
		t.Errorf("Prefix is wrong: %s, want: %s", p, prefix)
	}
}
