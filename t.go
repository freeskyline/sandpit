// Go example code

package main

import (
	"flag"
	"fmt"
	"kgen"
	"math"
	"os"
	"regexp"
	"runtime"
)

type ST struct {
	a int
	b int
	c int
}

func f() (result int) {
	defer func() {
		result++
	}()

	return 0
}

func testMain() {
	str1 := flag.String("s1", "aaaa", "The String One")
	str2 := flag.String("s2", "bbbb", "The String Two")

	fmt.Printf("Pi: %f\n", math.Pi)
	fmt.Printf("Pi: %e\n", math.Pi)
	fmt.Printf("Pi: %g\n", math.Pi)

	rrr := '中'
	fmt.Printf("%%c: %c\n", rrr)
	fmt.Printf("%%q: %q\n", rrr)

	str := "魑魅魍魉"
	fmt.Printf("%%s: %s\n", str)
	fmt.Printf("%%q: %q\n", str)
	fmt.Printf("GOPATH: %v\n", os.Getenv("GOPATH"))

	fmt.Println("func f(): ", f())
	fmt.Println("runtime.GOMAXPROCS: ", runtime.GOMAXPROCS)
	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())

	fmt.Println(*str1)
	fmt.Println(*str2)
	fmt.Println(os.Args)
}

func main() {
	flag.Parse()
	fmt.Println("-------------------------Begin-------------------------")
	defer fmt.Println("\n-------------------------End---------------------------")
	//testMain()

	//testEllipsis(100, 200, 300, 400, 500, 600, 700)
	testRegex()
	testKgen()
}

func testEllipsis0(a []int) {
	for i, v := range a {
		fmt.Printf("%d\t%v\n", i, v)
	}
}

func testEllipsis(args ...int) {
	testEllipsis0(args)
}

func testRegex() {
	re := regexp.MustCompile("(gopher){2}")
	fmt.Println(re.MatchString("gopher"))
	fmt.Println(re.MatchString("gophergopher"))
	fmt.Println(re.MatchString("gophergophergopher"))

	ptn := "^[0-9]{3}[;]{1}[0-9]{1}[A-HJ-NP-Ya-hj-np-y]{1}[0-5]{1}[0-9]{1}[1-7]{1}[A-Za-z]{1}[A-Za-z]{2}[0-9]{2}$"
	rgx := regexp.MustCompile(ptn)

	var str string
	str = "000;4Q317LAA02"
	fmt.Println(str, rgx.MatchString(str))
	str = "001;4Q317LAA0A"
	fmt.Println(str, rgx.MatchString(str))
	str = "001;4Q311LAA02"
	fmt.Println(str, rgx.MatchString(str))
	str = "001;4Z311LAA02"
	fmt.Println(str, rgx.MatchString(str))
}

func testKgen() {
	kgen.HelloKgen()
}
