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
	Address  string         `toml:"address"`
	BaudRate int            `toml:"baudRate"`
	DataBits int            `toml:"dataBits"`
	StopBits int            `toml:"stopBits"`
	Parity   string         `toml:"parity"`
	Timeout  time.Duration  `toml:"timeout"`
}

var applName string
var config  tomlConfig
var server *mbserver.Server
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
	buf.WriteString(fmt.Sprintln("Date Time    :", strTime, "LIHUI"))
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
	server = mbserver.NewServer()

	if config.Enable.MbTcpEn {
		str := config.MbTcp.Ip + ":" + strconv.Itoa(config.MbTcp.Port)
		err = server.ListenTCP(str)

		if err == nil {
			buf.WriteString(fmt.Sprintln("ModbusTCP Server listening on: ", str))
		} else {
			log.Printf("%v\n", err)
			buf.WriteString(fmt.Sprintln("%v", err))
		}
	}

	if config.Enable.MbRtuEn {
		var cnf serial.Config

		cnf.Address  = config.MbRtu.Address
		cnf.BaudRate = config.MbRtu.BaudRate
		cnf.DataBits = config.MbRtu.DataBits
		cnf.StopBits = config.MbRtu.StopBits
		cnf.Parity   = config.MbRtu.Parity
		cnf.Timeout  = config.MbRtu.Timeout * time.Second

		err = server.ListenRTU(&cnf)
		if err == nil {
			buf.WriteString(fmt.Sprintln("ModbusRTU Server listening on: ", "COM1"))
		} else {
			log.Printf("%v\n", err)
			buf.WriteString(fmt.Sprintln("%v", err))
		}
	}

	defer server.Close()
}

const strDefaultSettings =
`title = "Default Settings for ModbusTCP and ModbusRTU Server Simulator"

[enable]
  logsAntoGenEnabled = false
  modbusTCPEnabled = true
  modbusRTUEnabled = false

[modbusTCP]
  ip   = "127.0.0.1"
  port = 502

[modbusRTU]
  address  = "COM1"
  baudRate = 2400
  dataBits = 8         #Data bits: 5, 6, 7 or 8
  stopBits = 1         #Stop bits: 1 or 2
  parity   = "N"       #Parity: "N" - None, "E" - Even, "O" - Odd
  timeout  = 10        #Timeout Unit: Second
`
