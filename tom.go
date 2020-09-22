//tom.go, a TOML example

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title   string
	Calc    calculation `toml:"calculation"`
}

type calculation struct {
	Data1   []int      `toml:"data1"`
	Data2   []float64  `toml:"data2"`
}


func main() {
	var config tomlConfig
	var sum1 int
	var sum2 float64
	var strLine = "--------------------"

	if _, err := toml.DecodeFile("tom.settings", &config); err != nil {
		log.Println(err)
		return
	}

	log.Println(os.Args[0], " Title: ", config.Title, "\n")

	fmt.Println("Data1: ", config.Calc.Data1)
	for i, v := range config.Calc.Data1 {
		sum1 += v
		fmt.Printf("%d: %9d\n",i+1,v)
	}
	fmt.Println(strLine)
	fmt.Printf("Sum: %7d\n\n", sum1)

	fmt.Println("Data2: ", config.Calc.Data2)
	for i, v := range config.Calc.Data2 {
		sum2 += v
		fmt.Printf("%d: %9.2f\n",i+1,v)
	}
	fmt.Println(strLine)
	fmt.Printf("Sum: %7.2f\n", sum2)
}
