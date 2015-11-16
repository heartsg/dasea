package policy

import (
	"fmt"
	"testing"
)

func TestParseListRule(t *testing.T) {
	parser := &policyParser{}
	
	//test empty rule 
	c1 := parser.parseListRule([]interface{}{})
	if _, ok := c1.(*TrueCheck); !ok {
		t.Errorf("%s should be TrueCheck", c1)
	}
	
	//test nil inner rule
	c2 := parser.parseListRule([]interface{}{nil, []string{}})
	if _, ok := c2.(*TrueCheck); !ok {
		t.Errorf("%s should be TrueCheck", c1)
	}
	
	//test inner rules with string
	c3 := parser.parseListRule([]interface{}{"abc:def"})
	if c, ok := c3.(*GenericCheck); !ok {
		t.Errorf("%s should be Generic", c3)
	} else if c.kind != "abc" || c.match != "def" {
		t.Errorf("List rule failed")
	}
	
	c4 := parser.parseListRule([]interface{}{[]string{"abc:def"}})
	if c, ok := c4.(*GenericCheck); !ok {
		t.Errorf("%s should be Generic", c4)
	} else if c.kind != "abc" || c.match != "def" {
		t.Errorf("List rule failed")
	}
	
	//test multiple checks
	c5 := parser.parseListRule([]interface{}{"xyz:xyz", []string{"abc:def"}})
	if c, ok := c5.(*OrCheck); !ok {
		t.Errorf("%s should be OrCheck", c5)
	} else if len(c.rules) != 2 {
		t.Errorf("List rule failed")
	} else {
		if c_sub0, ok := c.rules[0].(*GenericCheck); !ok {
			t.Errorf("List rule failed")
		} else if c_sub0.kind != "xyz" || c_sub0.match != "xyz" {
			t.Errorf("List rule failed")
		}
		if c_sub1, ok := c.rules[1].(*GenericCheck); !ok {
			t.Errorf("List rule failed")
		} else if c_sub1.kind != "abc" || c_sub1.match != "def" {
			t.Errorf("List rule failed")
		}
	}
	
	//test multiple checks in list (and checks)
	c6 := parser.parseListRule([]interface{}{[]string{"xyz:xyz", "abc:def"}})
	if c, ok := c6.(*AndCheck); !ok {
		t.Errorf("%s should be AndCheck", c6)
	} else if len(c.rules) != 2 {
		t.Errorf("List rule failed")
	} else {
		if c_sub0, ok := c.rules[0].(*GenericCheck); !ok {
			t.Errorf("List rule failed")
		} else if c_sub0.kind != "xyz" || c_sub0.match != "xyz" {
			t.Errorf("List rule failed")
		}
		if c_sub1, ok := c.rules[1].(*GenericCheck); !ok {
			t.Errorf("List rule failed")
		} else if c_sub1.kind != "abc" || c_sub1.match != "def" {
			t.Errorf("List rule failed")
		}
	}
	
	//overall
	c7 := parser.parseListRule([]interface{}{
		[]string{"xyz:xyz", "abc:def", "123:456"}, 
		[]string{"xyz:xyz", "abc:def", "123:456"}, 
		[]string{"xyz:xyz", "abc:def", "123:456"}})
	if c, ok := c7.(*OrCheck); !ok {
		t.Errorf("%s should be OrCheck", c7)
	} else if len(c.rules) != 3 {
		t.Errorf("List rule failed")
	} else {
		for i, _ := range c.rules {
			if c_sub, ok := c.rules[i].(*AndCheck); !ok {
				t.Errorf("List rule failed")
			} else if len(c_sub.rules) != 3 {
				t.Errorf("List rule failed")
			} else {
				for j, _ := range c_sub.rules {
					if c_subsub, ok := c_sub.rules[j].(*GenericCheck); !ok {
						t.Errorf("List rule failed")
					} else {
						if (j == 0 && (c_subsub.kind != "xyz" || c_subsub.match != "xyz")) ||
							(j == 1 && (c_subsub.kind != "abc" || c_subsub.match != "def")) ||
							(j == 2 && (c_subsub.kind != "123" || c_subsub.match != "456")) {
							t.Errorf("List rule failed")
						}
					}
				}
			}
		}
	}
}

func TestTokenize(t *testing.T) {
	t1 := parseTokenize("role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin)")
	if len(t1) != 7 {
		t.Errorf("Parse tokenize error, want length 7, got %d", len(t1))
	}
	if t1[0].id != "check" || t1[1].id != "or" || t1[2].id != "(" || t1[3].id != "check" ||
		t1[4].id != "and" || t1[5].id != "check" || t1[6].id != ")" {
		t.Errorf("Parse tokenize error, id wrong")	
	}
}

func TestShiftReducer(t *testing.T) {
	t1 := parseTokenize("role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin)")
	s1 := &parseState { tokens: make([]*token, 0) }
	for _, tok1 := range t1 {
		s1.shift(tok1)
		s1.reduce()
	}
	if result1, ok := s1.result(); !ok {
		t.Errorf("Shift reduce failed")
	} else if result1.value != "(role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin))" {
		t.Errorf("Shift reduce failed")
	}
}

func TestPolicyParser(t *testing.T) {
	p := &policyParser{}
	r1, err := p.parseRule("role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin)")
	if err != nil {
		t.Errorf("Parse rule failed")
	}
	if fmt.Sprintf("%s", r1) != "(role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin))" {
		t.Errorf("Parse rule failed")
	}
}