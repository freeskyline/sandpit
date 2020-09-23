//tom.go, a TOML example

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title   string
	Enable  enable `toml:"enable"`
	Calc    calculation `toml:"calculation"`
}

type enable struct {
	LogsEn   bool  `toml:"logsAntoGenEnabled"`
	CalcEn   bool  `toml:"calculationEnabled"`
}

type calculation struct {
	Data1   []int      `toml:"data1"`
	Data2   []float64  `toml:"data2"`
}

var applName string
var config  tomlConfig
var err error
var buf bytes.Buffer
var fileSet string
var fileLog string



func init() {
	applName = os.Args[0]
	fileSet  = "tom.settings"

	if _, err = os.Stat(fileSet); os.IsNotExist(err) {
		data := []byte(strDefaultSettings)
		ioutil.WriteFile(fileSet, data, 0644)
	}
}

func main() {
	if _, err = toml.DecodeFile(fileSet, &config); err != nil {
		log.Println(err)
		return
	}

	strTime := time.Now().Format(time.RFC3339)
	buf.WriteString(fmt.Sprintln("Application  :", applName))
	buf.WriteString(fmt.Sprintln("Configuration:", fileSet))
	buf.WriteString(fmt.Sprintln("Config Title :", config.Title))
	buf.WriteString(fmt.Sprintln("Date Time    :", strTime))
	buf.WriteString("\n")

	executeSettings()
	fmt.Printf(buf.String())

	if config.Enable.LogsEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}
}

func executeSettings() {
	if config.Enable.CalcEn { fnCalculation() }
}

func fnCalculation() {
	const strLine = "--------------------\n"
	var sum1 int
	var sum2 float64

	buf.WriteString(fmt.Sprintln("Data1: ", config.Calc.Data1))
	for i, v := range config.Calc.Data1 {
		sum1 += v
		buf.WriteString(fmt.Sprintf("%d:\t%9d\n",i+1,v))
	}
	buf.WriteString(strLine)
	buf.WriteString(fmt.Sprintf("Sum:\t%9d\n\n", sum1))

	buf.WriteString(fmt.Sprintln("Data2: ", config.Calc.Data2))
	for i, v := range config.Calc.Data2 {
		sum2 += v
		buf.WriteString(fmt.Sprintf("%d:\t%9.2f\n",i+1,v))
	}
	buf.WriteString(strLine)
	buf.WriteString(fmt.Sprintf("Sum:\t%9.2f\n\n", sum2))
}

const strDefaultSettings =
`title = "TOM Default Settings"

[enable]
  logsAntoGenEnabled = true
  calculationEnabled = true

[calculation]
  data1 = [11, 22, 33]
  data2 = [100.01, 200.20, 300.50, 100.01, 200.20, 300.50, 100.01, 200.20, 300.50]
`
