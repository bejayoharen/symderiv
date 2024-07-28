package main
import "strings"
import "regexp"

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
