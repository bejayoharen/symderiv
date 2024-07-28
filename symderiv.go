package main
import "fmt"
import "os"
import "strconv"
import "strings"
import "regexp"

func usage() {
	fmt.Println("usage:")
	fmt.Println("\t" + os.Args[0] + " [Polynomial]")
	fmt.Println("for example:")
	fmt.Println("\t" + os.Args[0] + " \"12x^7-3x+5\"")
}

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

// -- Polynomial
// regexp for matching either + or -
var plusMinusRe *regexp.Regexp
func init() {
	plusMinusRe = regexp.MustCompile(`[+-]`)
}
type Polynomial []Term
func parsePolynomial( p string ) (Polynomial, error) {
	// clean and split on +/-
	pb := []byte(strings.ReplaceAll(p, " ", ""))
	idx := plusMinusRe.FindAllIndex(pb, -1)

	// make the first term
	ret := make( []Term, 0, len(idx) )
	end := len( pb )
	if len(idx) != 0 {
		end = idx[0][0]
	}
	//println( string(pb) )
	//println( end )
	//println( string(pb[:end]) )
	if end != 0 { // end will equal 0 if it's a unary + or -
        	t, err := parseTerm( string(pb[:end]) )
		if( err != nil ) {
			var t [0]Term
			return t[0:0], err
		}
		ret = append( ret, t )
	}

	// subsequent terms
	for i := 0; i<len(idx); i++ {
		s := idx[i][0]
		e := len( pb )
		if len( idx ) > i + 1 {
			e = idx[i+1][0]
		}
		ts := string(pb[s:e])
		if ts == "+" {
			continue;
		}
                t, err := parseTerm( ts )
		if( err != nil ) {
			var t [0]Term
			return t[0:0], err
		}
                ret = append( ret, t )
        }
	return ret, nil
}
func (p Polynomial) String() string {
	var ret = make([]string, 0, len(p))
	for i := 0; i<len(p); i++ {
		t := p[i]
		s := t.String()
		ret = append( ret, s )
	}
	return strings.Join( ret, " + " )
}
func (p Polynomial) Derivative() Polynomial {
	var ret = make([]Term, 0, len(p))
	for i := 0; i<len(p); i++ {
		if p[i].isZero() || p[i].isConstant() {
			continue
		}
                t := p[i].Derivative()
		if( !t.isZero() ) {
                	ret = append( ret, t )
		}
        }
	return ret
}
func (p Polynomial) isSingleVar() bool {
	variable := ""
	for i := 0; i<len(p); i++ {
                v := p[i].variable
                if v == "" {
                        continue
                }
		if variable == "" {
			variable = v
		} else if variable != v {
			return false
		} else {
		}
        }
	return true
}

// -- main function
func main() {
	if len( os.Args ) != 2 {
		usage()
		return
	}
	if os.Args[1] == "help" || os.Args[1] == "-h" {
		usage()
		return
	}
	polyString := os.Args[1]

	poly, err := parsePolynomial( polyString )
	if err != nil {
		fmt.Println( "Error: ", err )
		usage()
		return
	}
	if !poly.isSingleVar() {
		fmt.Println( "Error: Polynomials must be of a single variable." )
		usage()
		return
	}
	deriv := poly.Derivative();

	fmt.Println( "Original Polynomial: ", poly )
	fmt.Println( "Derivative of orig : ", deriv )

}
