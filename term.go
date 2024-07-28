package main
import "fmt"
import "strconv"
import "strings"
import "regexp"

// -- Term is a single term in a polynomial
// regexp to parse a single term
var fullTermRe *regexp.Regexp
func init() {
	fullTermRe  = regexp.MustCompile(`(\+|\-)?(?P<c>[\d]*)((?P<v>[a-zA-Z]+)(\^(?P<e>[\d]+))?)?(?P<x>.*)`)
}
type Term struct {
	coefficient int
	exponent int
	variable string
}
func parseTerm( t string ) (Term, error) {
	t = strings.ReplaceAll(t, " ", "")
	//fmt.Println( "parsing term: " + t )
	m := fullTermRe.FindStringSubmatch( t )
	//fmt.Println( len(m), ": " , m )
        if len(m) != 8 {
		return Term{ coefficient: 0, exponent: 0, variable: "x" }, fmt.Errorf("Could not parse term: %s", t)
	}
	var c int = 1
	if( len(strings.TrimSpace(m[2])) > 0 ) {
		var er error = nil
		c, er = strconv.Atoi(strings.TrimSpace(m[2]))
       		if er != nil {
			return Term{ coefficient: 0, exponent: 0, variable: "x" },
				fmt.Errorf("Could not parse coefficient (%s) in term: %s", m[2], t)
		}
	}
	var ex int = 1
	if( m[6] != "" ) {
		var er error = nil
		ex, er = strconv.Atoi(strings.TrimSpace(m[6]))
        	if er != nil {
			return Term{ coefficient: 0, exponent: 0, variable: "x" },
				fmt.Errorf("Could not parse exponenent (%s) in term: %s", m[6], t)
		}
	}
	if m[4] == "" {
		ex = 0
	}
	if m[1] == "-" {
		c = -c
	}
	ret := Term{ coefficient: c, exponent: ex, variable: m[4] }
	//fmt.Println( ret )
	return ret, nil
}
func (t Term) String() string {
	if t.isZero() {
		return "0"
	}
	if t.isConstant() {
		return fmt.Sprintf( "%d", t.coefficient )
	}
	if( t.coefficient == 1 && t.exponent == 0 ) {
		return t.variable
	}
	if( t.coefficient == -1 && t.exponent == 0 ) {
                return "-" + t.variable
        }
	ret := ""
	switch( t.coefficient ) {
	case 1:
		break
	case -1:
		ret = ret + "-"
		break
	default:
		ret = ret + fmt.Sprintf( "%d", t.coefficient )
	}
	switch( t.exponent ) {
	case 0:
		break
	case 1:
		ret = ret + t.variable
		break
	default:
		ret = ret + fmt.Sprintf( "%s^%d",
                        t.variable,
                        t.exponent )
	}

	return ret
}
func (t Term) isZero() bool {
	return t.coefficient == 0;
}
func (t Term) isConstant() bool {
	return t.exponent == 0 || t.variable == "";
}
func (t Term) Derivative() Term {
	if t.isConstant() || t.isZero() {
		return Term{
			coefficient: 0,
			exponent: 1,
			variable: t.variable }
	}
	return Term{
		coefficient: t.coefficient * t.exponent,
		exponent: t.exponent - 1,
		variable: t.variable }
}
