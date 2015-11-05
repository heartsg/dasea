package config

import (
	"testing"
)

type TestTagInnerStruct struct {
	A [] float64 `default:"1.0,2.0,3.0"`
	B [] string `default:"abc,def,xyz"`
	C int
}
type TestTagStruct struct {
	A int `default:"1"`
	B int64 `default:"-3"`
	C string `default:"xyz"`
	D float64 `default:"3.14"`
	E []int `default:"1,2,3"`
	F TestTagInnerStruct
}
func TestDefaultValues(t *testing.T) {
	tl := &TagLoader{}
	ts := new(TestTagStruct)
	if err := tl.Load(ts); err != nil {
		t.Error(err)
	}

	if ts.A != 1 {
		t.Errorf("A value is wrong: %d, want: %d", ts.A, 1)
	}
	
	if ts.B != -3 {
		t.Errorf("B value is wrong: %d, want: %d", ts.B, -3)
	}
	
	if ts.C != "xyz" {
		t.Errorf("C value is wrong: %s, want: %s", ts.C, "xyz")
	}
	
	if ts.D != 3.14 {
		t.Errorf("D value is wrong: %f, want: %f", ts.D, 3.14)
	}
	
	if len(ts.E) != 3 || ts.E[0] != 1 || ts.E[1] != 2 || ts.E[2] != 3 {
		t.Errorf("E value is wrong: %s, want: 1,2,3", ts.E)
	}
	
	if len(ts.F.A) != 3 || ts.F.A[0] != 1.0 || ts.F.A[1] != 2.0 || ts.F.A[2] != 3.0 {
		t.Errorf("F.A value is wrong: %s, want: 1.0,2.0,3.0", ts.F.A)
	}
	
	if len(ts.F.B) != 3 || ts.F.B[0] != "abc" || ts.F.B[1] != "def" || ts.F.B[2] != "xyz" {
		t.Errorf("F.B value is wrong: %s, want: abc,def,xyz", ts.F.B)
	}

	if ts.F.C != 0 {
		t.Errorf("F.C value is wrong: %d, want: %d", ts.F.C, 0)
	}
}
