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
/*    if rule == "" {
        return &TrueCheck{}
    }
    
    //start parsing token stream
    state := &parseState{}
    for _ tok := range parseTokenize(rule) {
        state.shift(tok)
    }

    return state.result*/
    return &FalseCheck{}
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
    tokens []token
}
func (s *parseState) shift(tok *token) {
    s.tokens = append(s.tokens, token)
    s.reduce()
}
func (s *parseState) reduce() {
    
}


def reducer(*tokens):
    """Decorator for reduction methods.

    Arguments are a sequence of tokens, in order, which should trigger running
    this reduction method.
    """

    def decorator(func):
        # Make sure we have a list of reducer sequences
        if not hasattr(func, 'reducers'):
            func.reducers = []

        # Add the tokens to the list of reducer sequences
        func.reducers.append(list(tokens))

        return func

    return decorator


class ParseStateMeta(type):
    """Metaclass for the :class:`.ParseState` class.

    Facilitates identifying reduction methods.
    """

    def __new__(mcs, name, bases, cls_dict):
        """Create the class.

        Injects the 'reducers' list, a list of tuples matching token sequences
        to the names of the corresponding reduction methods.
        """

        reducers = []

        for key, value in cls_dict.items():
            if not hasattr(value, 'reducers'):
                continue
            for reduction in value.reducers:
                reducers.append((reduction, key))

        cls_dict['reducers'] = reducers

        return super(ParseStateMeta, mcs).__new__(mcs, name, bases, cls_dict)


@six.add_metaclass(ParseStateMeta)
class ParseState(object):
    """Implement the core of parsing the policy language.

    Uses a greedy reduction algorithm to reduce a sequence of tokens into
    a single terminal, the value of which will be the root of the
    :class:`Check` tree.

    .. note::

        Error reporting is rather lacking.  The best we can get with this
        parser formulation is an overall "parse failed" error. Fortunately, the
        policy language is simple enough that this shouldn't be that big a
        problem.
    """

    def __init__(self):
        """Initialize the ParseState."""

        self.tokens = []
        self.values = []

    def reduce(self):
        """Perform a greedy reduction of the token stream.

        If a reducer method matches, it will be executed, then the
        :meth:`reduce` method will be called recursively to search for any more
        possible reductions.
        """

        for reduction, methname in self.reducers:
            if (len(self.tokens) >= len(reduction) and
                    self.tokens[-len(reduction):] == reduction):
                # Get the reduction method
                meth = getattr(self, methname)

                # Reduce the token stream
                results = meth(*self.values[-len(reduction):])

                # Update the tokens and values
                self.tokens[-len(reduction):] = [r[0] for r in results]
                self.values[-len(reduction):] = [r[1] for r in results]

                # Check for any more reductions
                return self.reduce()

    def shift(self, tok, value):
        """Adds one more token to the state.

        Calls :meth:`reduce`.
        """

        self.tokens.append(tok)
        self.values.append(value)

        # Do a greedy reduce...
        self.reduce()

    @property
    def result(self):
        """Obtain the final result of the parse.

        :raises ValueError: If the parse failed to reduce to a single result.
        """

        if len(self.values) != 1:
            raise ValueError('Could not parse rule')
        return self.values[0]

    @reducer('(', 'check', ')')
    @reducer('(', 'and_expr', ')')
    @reducer('(', 'or_expr', ')')
    def _wrap_check(self, _p1, check, _p2):
        """Turn parenthesized expressions into a 'check' token."""

        return [('check', check)]

    @reducer('check', 'and', 'check')
    def _make_and_expr(self, check1, _and, check2):
        """Create an 'and_expr'.

        Join two checks by the 'and' operator.
        """

        return [('and_expr', _checks.AndCheck([check1, check2]))]

    @reducer('and_expr', 'and', 'check')
    def _extend_and_expr(self, and_expr, _and, check):
        """Extend an 'and_expr' by adding one more check."""

        return [('and_expr', and_expr.add_check(check))]

    @reducer('check', 'or', 'check')
    def _make_or_expr(self, check1, _or, check2):
        """Create an 'or_expr'.

        Join two checks by the 'or' operator.
        """

        return [('or_expr', _checks.OrCheck([check1, check2]))]

    @reducer('or_expr', 'or', 'check')
    def _extend_or_expr(self, or_expr, _or, check):
        """Extend an 'or_expr' by adding one more check."""

        return [('or_expr', or_expr.add_check(check))]

    @reducer('not', 'check')
    def _make_not_expr(self, _not, check):
        """Invert the result of another check."""

        return [('check', _checks.NotCheck(check))]
