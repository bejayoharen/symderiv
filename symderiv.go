package main
import "fmt"
import "os"

func usage() {
	fmt.Println("usage:")
	fmt.Println("\t" + os.Args[0] + " [Polynomial]")
	fmt.Println("for example:")
	fmt.Println("\t" + os.Args[0] + " \"12x^7-3x+5\"")
}

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
