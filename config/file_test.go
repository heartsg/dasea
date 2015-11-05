package config

import (
	"time"
	"testing"
)


type TestFileServer struct {
	Name       string 
	Port       int    
	ID         int64
	Labels     []int
	Enabled    bool
	Users      []string
	Postgres   TestFilePostgres
	unexported string
	Interval   time.Duration
}


type TestFilePostgres struct {
	Enabled           bool
	Port              int      
	Hosts             []string 
	DBName            string   
	AvailabilityRatio float64
	unexported        string
}

var (
	testTOML = "testdata/config.toml"
	testJSON = "testdata/config.json"
	testINI = "testdata/config.ini"
)

func getDefaultFileServer() *TestFileServer {
	return &TestFileServer{
		Name:     "koding",
		Port:     6060,
		Enabled:  true,
		ID:       1234567890,
		Labels:   []int{123, 456},
		Users:    []string{"ankara", "istanbul"},
		Interval: 10 * time.Second,
		Postgres: TestFilePostgres{
			Enabled:           true,
			Port:              5432,
			Hosts:             []string{"192.168.2.1", "192.168.2.2", "192.168.2.3"},
			DBName:            "configdb",
			AvailabilityRatio: 8.23,
		},
	}
}

func testFileStruct(t *testing.T, s *TestFileServer, d *TestFileServer) {
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


func TestToml(t *testing.T) {
	m := &TOMLLoader{Path: testTOML}

	s := &TestFileServer{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFileStruct(t, s, getDefaultFileServer())
}


func TestJSON(t *testing.T) {
	m := &JSONLoader{Path: testJSON}

	s := &TestFileServer{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}

	testFileStruct(t, s, getDefaultFileServer())
}

func TestINI(t *testing.T) {
	m := &INILoader{Path: testINI}
	s := &TestFileServer{}
	if err := m.Load(s); err != nil {
		t.Error(err)
	}
	
	testFileStruct(t, s, getDefaultFileServer())
}