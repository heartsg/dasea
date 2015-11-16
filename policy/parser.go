package policy

import (
	"fmt"
    "strings"
)

//Based on openstack oslo.policy, used to check authorization policy


//PolicyParser is used to parse rules into key, check object map
type policyParser struct {
}


//Conjunction operators are available, allowing for more expressiveness
//in crafting policies. So, in the policy language, the previous check in
//list-of-lists becomes::
//
//   role:admin or ({{.project_id}}:{{.project_id}} and role:projectadmin)
//
//The policy language also has the ``not`` operator, allowing a richer
//policy rule::
//
//   {{.project_id}}:{{.project_id}} and not role:dunce
//
//Finally, two special policy checks should be mentioned; the policy
//check "@" will always accept an access, and the policy check "!" will
//always reject an access.  (Note that if a rule is either the empty
//list ("[]") or the empty string, this is equivalent to the "@" policy
//check.)  Of these, the "!" policy check is probably the most useful,
//as it allows particular rules to be explicitly disabled.
func (p *policyParser) parseTextRule(rule string) PolicyCheck {
    if rule == "" {
        return &TrueCheck{}
    }
    
    //start parsing token stream
    state := &parseState{tokens:make([]*token, 0)}
    for _, tok := range parseTokenize(rule) {
        state.shift(tok)
    }

    if result, ok := state.result(); ok {
        return result.check
    } else {
        return &FalseCheck{}
    }
}

   /* # Parse the token stream
    state = ParseState()
    for tok, value in _parse_tokenize(rule):
        state.shift(tok, value)

    try:
        return state.result
    except ValueError:
        # Couldn't parse the rule
        LOG.exception(_LE('Failed to understand rule %s'), rule)

        # Fail closed
        return _checks.FalseCheck()*/


//In the list-of-lists representation, each check inside the innermost
//list is combined as with an "and" conjunction--for that check to pass,
//all the specified checks must pass.  These innermost lists are then
//combined as with an "or" conjunction. As an example, take the following
//rule, expressed in the list-of-lists representation::
//
//    [["role:admin"], ["{{.project_id}}:{{.project_id}}", "role:projectadmin"]]
//
func (p *policyParser) parseListRule(rule []interface{}) PolicyCheck {
	//empty rule defaults to true
	if rule == nil || len(rule) == 0 {
		return &TrueCheck{}
	}
    
    //outer list joined by or, inner list joined by and
    orList := make([]PolicyCheck, 0)
    for _, innerRule := range rule {
        if innerRule == nil {
            continue
        }
        
        //if it is a string, make it a list with len 1
        if s, ok := innerRule.(string); ok {
            innerRule = []string{ s }
        }
        
        if rules, ok := innerRule.([]string); ok {
            if len(rules) == 0 {
                continue
            }
            
            andList := make([]PolicyCheck, 0)
            for _, r := range rules {
                andList = append(andList, ParseCheck(r))
            }
            
            if len(andList) == 1 {
                orList = append(orList, andList[0])
            } else {
                orList = append(orList, &AndCheck{rules:andList})
            }
        } else {
            //I do not understand this list?
            orList = append(orList, &FalseCheck{})
        }
    }
    
    if len(orList) == 0 {
        return &TrueCheck{}
    } else if len(orList) == 1 {
        return orList[0]
    } else {
        return &OrCheck{rules:orList}
    }
}

func (p *policyParser) parseRule(rule interface{}) (PolicyCheck, error) {
	switch r := rule.(type) {
	case string:
		return p.parseTextRule(r), nil
	case []interface{}:
		return p.parseListRule(r), nil
	}
	
	return nil, fmt.Errorf("policy.policyParser: rule format not supported")
}


//Tokenizer for the policy rules
type token struct {
    id string
    value string
    check PolicyCheck
}
func parseTokenize(rule string) []*token {
    r := make([]*token, 0)
    for _, tok := range strings.Fields(rule) {
        if tok == "" {
            continue
        }
        
        //leading parentheses
        clean := strings.TrimLeft(tok, "(")
        for i := 0; i < len(tok) - len(clean); i++ {
            r = append(r, &token{id:"(", value:"(", check:nil})
        }
        
        if clean == "" {
            continue
        } else {
            tok = clean
        }
        
        //handle trailing parentheses
        clean = strings.TrimRight(tok, ")")
        trail := len(tok) - len(clean)
        
        lowered := strings.ToLower(clean)
        if lowered == "and" || lowered == "or" || lowered == "not" {
            r = append(r, &token{id:lowered, value:clean, check:nil})
        } else if clean != "" {
            if len(clean) >= 2 && ((clean[0] == '"' && clean[len(clean) - 1] == '"') ||
                (clean[0] == '\'' && clean[len(clean) - 1] == '\'') ) {
                r = append(r, &token{id:"string", value:clean[1:len(clean)-1], check:nil})
            } else {
                r = append(r, &token{id:"check", value:clean, check:ParseCheck(clean)})
            }
            
        }
        
        for i := 0; i < trail; i++ {
            r = append(r, &token{id:")", value:")", check:nil})
        }
    }
    return r
}

type parseState struct {
    tokens []*token
}
func (s *parseState) shift(tok *token) {
    s.tokens = append(s.tokens, tok)
    s.reduce()
}
func (s *parseState) reduce() {
    if ok := s.reduceWrapCheck(); ok {
        s.reduce()
        return
    }
    if ok := s.reduceMakeAndExpr(); ok {
        s.reduce()
        return
    }
    if ok := s.reduceExtendAndExpr(); ok {
        s.reduce()
        return
    }
    if ok := s.reduceMakeOrExpr(); ok {
        s.reduce()
        return
    }
    if ok := s.reduceExtendOrExpr(); ok {
        s.reduce()
        return
    }
    if ok := s.reduceMakeNotExpr(); ok {
        s.reduce()
        return
    }
    
    //failed
    return
}
func (s *parseState) reduceWrapCheck() bool {
    length := len(s.tokens)
    if length >= 3 && s.tokens[length - 3].id == "(" && s.tokens[length - 1].id == ")" &&
        (s.tokens[length - 2].id == "check" || s.tokens[length - 2].id == "and_expr" ||
        s.tokens[length - 2].id == "or_expr") {    
        s.tokens = append(s.tokens[0:length-3], &token{"check", s.tokens[length-2].value, s.tokens[length-2].check})
        
        return true
    }
    return false
}
func (s *parseState) reduceMakeAndExpr() bool {
    length := len(s.tokens)
    if length >= 3 && s.tokens[length - 3].id == "check" && s.tokens[length - 2].id == "and" && s.tokens[length - 1].id == "check" {
        andCheck := &AndCheck{rules:[]PolicyCheck{s.tokens[length - 3].check, s.tokens[length - 1].check}}
        s.tokens = append(s.tokens[0:length - 3], &token{"and_expr", fmt.Sprintf("%s", andCheck), andCheck})
        
        return true
    }
    return false
}
func (s *parseState) reduceExtendAndExpr() bool {
    length := len(s.tokens)
    if length >= 3 && s.tokens[length - 3].id == "and_expr" && s.tokens[length - 2].id == "and" && s.tokens[length - 1].id == "check" {
        andCheck := s.tokens[length - 3].check.(*AndCheck)
        andCheck.AddCheck(s.tokens[length - 1].check)
        s.tokens = append(s.tokens[0:length - 3], &token{"and_expr", fmt.Sprintf("%s", andCheck), andCheck})
        return true
    }
    return false
}
func (s *parseState) reduceMakeOrExpr() bool {
    length := len(s.tokens)
    if length >= 3 && s.tokens[length - 3].id == "check" && s.tokens[length - 2].id == "or" && s.tokens[length - 1].id == "check" {
        orCheck := &OrCheck{rules:[]PolicyCheck{s.tokens[length - 3].check, s.tokens[length - 1].check}}
        s.tokens = append(s.tokens[0:length - 3], &token{"or_expr", fmt.Sprintf("%s", orCheck), orCheck})
        
        return true
    }
    return false
}
func (s *parseState) reduceExtendOrExpr() bool {
    length := len(s.tokens)
    if length >= 3 && s.tokens[length - 3].id == "or_expr" && s.tokens[length - 2].id == "or" && s.tokens[length - 1].id == "check" {
        orCheck := s.tokens[length - 3].check.(*OrCheck)
        orCheck.AddCheck(s.tokens[length - 1].check)
        s.tokens = append(s.tokens[0:length - 3], &token{"or_expr", fmt.Sprintf("%s", orCheck), orCheck})
        
        return true
    }
    return false
}
func (s *parseState) reduceMakeNotExpr() bool {
    length := len(s.tokens)
    if length >= 2 && s.tokens[length - 2].id == "not" && s.tokens[length - 1].id == "check" {
        notCheck := &NotCheck{rule:s.tokens[length - 1].check}
        s.tokens = append(s.tokens[0:length - 2], &token{"check", fmt.Sprintf("%s", notCheck), notCheck})
        
        return true
    }
    return false
}
func (s *parseState) result() (*token, bool) {
    if len(s.tokens) != 1 {
        return nil, false
    } else {
        return s.tokens[0], true
    }
}