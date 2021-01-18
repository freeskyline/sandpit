// Go example code

package main

import (
	"flag"
	"fmt"
	"math"
	"os"
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

	fmt.Printf("Pi: %f\n", math.Pi);
	fmt.Printf("Pi: %e\n", math.Pi);
	fmt.Printf("Pi: %g\n", math.Pi);

	rrr := '中'
	fmt.Printf("%%c: %c\n", rrr);
	fmt.Printf("%%q: %q\n", rrr);

	str := "魑魅魍魉"
	fmt.Printf("%%s: %s\n", str);
	fmt.Printf("%%q: %q\n", str);
	fmt.Printf("GOPATH: %v\n", os.Getenv("GOPATH"));

	fmt.Println("func f(): ", f())
	fmt.Println("runtime.GOMAXPROCS: ", runtime.GOMAXPROCS)
	fmt.Println("runtime.NumCPU(): ", runtime.NumCPU())

	fmt.Println(*str1)
	fmt.Println(*str2)
	fmt.Println(os.Args)
}

func main() {
	flag.Parse();
	fmt.Println("-------------------------Begin-------------------------")
	defer fmt.Println("\n-------------------------End---------------------------")
	//testMain()

	testEllipsis(100, 200, 300)
}

func testEllipsis(args ...int) {
	for i, v := range args {
		fmt.Printf("%d\t%v\n", i, v)
	}
}
