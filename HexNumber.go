// HexNumber.go

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func banner() func() {
	const (
		strApp = "HexNumber"
		strVer = "v1.0.0"
		strUTC = "Build LIHUI 2020-07-25 16:44:09.20214347 +0000 UTC"
		str = "--------------------------------------------------------------------------------"
	)

	fmt.Println(str,"\n")

	return func() {
		fmt.Println()
		fmt.Println("", strApp, strVer)
		fmt.Println("", strUTC)
		fmt.Println(str)
	}
}

func strToNumber() {
	if len(os.Args) > 1 {
		val, err := strconv.ParseUint(os.Args[1], 0, strconv.IntSize)

		log.Printf("%q", os.Args[1])
		if err == nil {
			v := uint32(val)
			fmt.Printf(" Dec: %d\n", v)
			fmt.Printf(" Hex: %08X\n",v)

			v1 := (v <<  0) >> 24
			v2 := (v <<  8) >> 24
			v3 := (v << 16) >> 24
			v4 := (v << 24) >> 24
			fmt.Printf(" Bin: %08b  %08b %08b %08b\n", v1, v2, v3, v4)
		} else {
			log.Println(err)
		}
	}
}

func main() {
	defer banner()()

	strToNumber()
}
