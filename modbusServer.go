//modbusServer.go, ModbusTCP and ModbusRTU Server Simulator

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/goburrow/serial"
	"github.com/tbrandon/mbserver"
)

type tomlConfig struct {
	Title   string
	Enable  enable `toml:"enable"`
	MbTcp   modbusTCP `toml:"modbusTCP"`
	MbRtu   modbusRTU `toml:"modbusRtu"`
}

type enable struct {
	LogsEn   bool  `toml:"logsAntoGenEnabled"`
	MbTcpEn  bool  `toml:"modbusTCPEnabled"`
	MbRtuEn  bool  `toml:"modbusRTUEnabled"`
}

type modbusTCP struct {
	Ip    string    `toml:"ip"`
	Port  int       `toml:"port"`
}

type modbusRTU struct {

}

var applName string
var config  tomlConfig
var err error
var buf bytes.Buffer
var fileSet string
var fileLog string

var mbServerRTU *mbserver.Server
var mbServerTCP *mbserver.Server

func init() {
	strApp := path.Base(os.Args[0])
	strExt := path.Ext(strApp)

	applName = strings.TrimSuffix(strApp, strExt)
	fileSet  = applName + ".settings"

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

	for {
		time.Sleep(time.Second)
	}
}

func executeSettings() {
	if config.Enable.MbTcpEn { 
		mbServerTCP = mbserver.NewServer()

		str := config.MbTcp.Ip + ":" + strconv.Itoa(config.MbTcp.Port)
		err = mbServerTCP.ListenTCP(str)
		if err != nil {
			log.Printf("%v\n", err)
			buf.WriteString(fmt.Sprintln("%v", err))
		} else {
			buf.WriteString(fmt.Sprintln("ModbusTCP Server listening on: ", str))
		}

		defer mbServerTCP.Close()
	}

	if config.Enable.MbRtuEn { 
		mbServerRTU = mbserver.NewServer()

		err = mbServerRTU.ListenRTU(&serial.Config{
			Address:  "COM1",
			BaudRate: 115200,
			DataBits: 8,
			StopBits: 1,
			Parity:   "N",
			Timeout:  10 * time.Second})
		if err != nil {
			log.Printf("%v\n", err)
			buf.WriteString(fmt.Sprintln("%v", err))
		} else {
			buf.WriteString(fmt.Sprintln("ModbusRTU Server listening on: ", "COM1"))
		}

		defer mbServerRTU.Close()
	}
}

const strDefaultSettings =
`title = "Default Settings for ModbusTCP and ModbusRTU Server Simulator"

[enable]
  logsAntoGenEnabled = false
  modbusTCPEnabled = true
  modbusRTUEnabled = false

[modbusTCP]
  ip = "127.0.0.1"
  port = 502

[modbusRTU]
`
