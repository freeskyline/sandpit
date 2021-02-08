// Go example code linking with a Lib (*.a)

package main

import (
	"fmt"
	"lib"
)

func main() {
	fmt.Println("-------------------------Begin-------------------------")
	defer fmt.Println("\n-------------------------End---------------------------")

	testKgen()
}

func testKgen() {
	lib.HelloKgen()
}
