package config

import (
	"strings"
	"testing"
	"time"
	"github.com/fatih/structs"
)

type TestFlagServer struct {
	Name       string
	Port       int
	ID         int64
	Labels     []int
	Enabled    bool
	Users      []string
	Postgres   TestFlagPostgres
	unexported string
	Interval   time.Duration
}

type TestFlagPostgres struct {
	Enabled           bool
	Port              int      
	Hosts             []string 
	DBName            string   
	AvailabilityRatio float64
	unexported        string
}

type TestFlagFlattenedServer struct {
	Postgres TestFlagPostgres
}

type TestFlagCamelCaseServer struct {
	AccessKey         string
	Normal            string
	DBName            string 
	AvailabilityRatio float64
}

type TestFlagTagServer struct {
	AccessKey         string `flag:"acc"`
	Normal            string `flag:"nor"`
	DBName            string  `flag:"db"`
	AvailabilityRatio float64 `flag:"ar"`
}
func getDefaultFlagServer() *TestFlagServer {
	return &TestFlagServer{
		Name:     "koding",
		Port:     6060,
		Enabled:  true,
		ID:       1234567890,
		Labels:   []int{123, 456},
		Users:    []string{"ankara", "istanbul"},
		Interval: 10 * time.Second,
		Postgres: TestFlagPostgres{
			Enabled:           true,
			Port:              5432,
			Hosts:             []string{"192.168.2.1", "192.168.2.2", "192.168.2.3"},
			DBName:            "configdb",
			AvailabilityRatio: 8.23,
		},
	}
}

func getDefaultFlagCamelCaseServer() *TestFlagCamelCaseServer {
	return &TestFlagCamelCaseServer{
		AccessKey:         "123456",
		Normal:            "normal",
		DBName:            "configdb",
		AvailabilityRatio: 8.23,
	}
}

func getDefaultFlagTagServer() *TestFlagTagServer {
	return &TestFlagTagServer{
		AccessKey:         "123456",
		Normal:            "normal",
		DBName:            "configdb",
		AvailabilityRatio: 8.23,
	}
}

func testFlagStruct(t *testing.T, s *TestFlagServer, d *TestFlagServer) {
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
	// `x == true` infact.
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

func testFlagFlattenedStruct(t *testing.T, s *TestFlagFlattenedServer, d *TestFlagServer) {
	// Explicitly state that Enabled should be true, no need to check
	// `x == true` infact.
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

func testFlagCamelcaseStruct(t *testing.T, s *TestFlagCamelCaseServer, d *TestFlagCamelCaseServer) {
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

func testFlagTagStruct(t *testing.T, s *TestFlagTagServer, d *TestFlagTagServer) {
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



func TestFlag(t *testing.T) {
	m := &FlagLoader{}
	s := &TestFlagServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagStruct(t, s, getDefaultFlagServer())
}

func TestFlagWithPrefix(t *testing.T) {
	const prefix = "Prefix"

	m := FlagLoader{Prefix: prefix}
	s := &TestFlagServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, prefix)

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagStruct(t, s, getDefaultFlagServer())
}

func TestFlattenFlags(t *testing.T) {
	m := FlagLoader{
		Flatten: true,
	}
	s := &TestFlagFlattenedServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagFlattenedStruct(t, s, getDefaultFlagServer())
}

func TestCamelcaseFlags(t *testing.T) {
	m := FlagLoader{
		CamelCase: true,
	}
	s := &TestFlagCamelCaseServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagCamelcaseStruct(t, s, getDefaultFlagCamelCaseServer())
}

func TestFlattenAndCamelCaseFlags(t *testing.T) {
	m := FlagLoader{
		Flatten:   true,
		CamelCase: true,
	}
	s := &TestFlagFlattenedServer{}

	// get flags
	args := getFlags(t, "TestFlagFlattenedCamelCaseServer", "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagFlattenedStruct(t, s, getDefaultFlagServer())
}

func TestTagFlags(t *testing.T) {
	m := FlagLoader{
	}
	s := &TestFlagTagServer{}
	structName := structs.Name(s)

	// get flags
	args := getFlags(t, structName, "")

	m.Args = args[1:]

	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFlagTagStruct(t, s, getDefaultFlagTagServer())
}

// getFlags returns a slice of arguments that can be passed to flag.Parse()
func getFlags(t *testing.T, structName, prefix string) []string {
	if structName == "" {
		t.Fatal("struct name can not be empty")
	}

	var flags map[string]string
	switch structName {
	case "TestFlagServer":
		flags = map[string]string{
			"-name":                       "koding",
			"-port":                       "6060",
			"-enabled":                    "",
			"-users":                      "ankara,istanbul",
			"-interval":                   "10s",
			"-id":                         "1234567890",
			"-labels":                     "123,456",
			"-postgres-enabled":           "",
			"-postgres-port":              "5432",
			"-postgres-hosts":             "192.168.2.1,192.168.2.2,192.168.2.3",
			"-postgres-dbname":            "configdb",
			"-postgres-availabilityratio": "8.23",
		}
	case "TestFlagFlattenedServer":
		flags = map[string]string{
			"--enabled":           "",
			"--port":              "5432",
			"--hosts":             "192.168.2.1,192.168.2.2,192.168.2.3",
			"--dbname":            "configdb",
			"--availabilityratio": "8.23",
		}
	case "TestFlagFlattenedCamelCaseServer":
		flags = map[string]string{
			"--enabled":            "",
			"--port":               "5432",
			"--hosts":              "192.168.2.1,192.168.2.2,192.168.2.3",
			"--db-name":            "configdb",
			"--availability-ratio": "8.23",
		}
	case "TestFlagCamelCaseServer":
		flags = map[string]string{
			"--access-key":         "123456",
			"--normal":             "normal",
			"--db-name":            "configdb",
			"--availability-ratio": "8.23",
		}
	case "TestFlagTagServer":
		flags = map[string]string{
			"--acc":         "123456",
			"--nor":         "normal",
			"--db":          "configdb",
			"-ar":           "8.23",
		}
	}

	prefix = strings.ToLower(prefix)

	args := []string{"multiconfig-test"}
	for key, val := range flags {
		flag := key
		if prefix != "" {
			flag = "-" + prefix + key
		}

		if val == "" {
			args = append(args, flag)
		} else {
			args = append(args, flag, val)
		}
	}

	return args
}
