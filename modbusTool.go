//modbusServer.go, ModbusTCP and ModbusRTU Client Tool

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
	"github.com/goburrow/modbus"
)

type tomlConfig struct {
	Title   string    `toml:"title"`
	Enable  enable    `toml:"enable"`
	MbTcp   modbusTCP `toml:"modbusTCP"`
	MbRtu   modbusRTU `toml:"modbusRTU"`
	MbReg   registers `toml:"registers"`
}

type enable struct {
	LogsEn   bool  `toml:"logsAntoGenEnabled"`
	MbRtuEn  bool  `toml:"modbusRTUEnabled"`
}

type modbusTCP struct {
	Ip       string    `toml:"ip"`
	Port     int       `toml:"port"`
	SlaveId  byte
	Timeout  time.Duration
}

type modbusRTU struct {
	Address  string         `toml:"address"`
	BaudRate int            `toml:"baudRate"`
	DataBits int            `toml:"dataBits"`
	StopBits int            `toml:"stopBits"`
	Parity   string         `toml:"parity"`
	SlaveId  byte
	Timeout  time.Duration  `toml:"timeout"`
}

type registers struct {
	DiscreteInputs    [][2]uint16  `toml:"discreteInputs"`
	Coils             [][2]uint16  `toml:"coils"`
	HoldingRegisters  [][2]uint16  `toml:"holdingRegisters"`
	InputRegisters    [][2]uint16  `toml:"inputRegisters"`
}

var applName string
var config  tomlConfig
var err error
var results []byte
var buf bytes.Buffer
var fileSet string
var fileLog string

var tcpHandler *modbus.TCPClientHandler
var rtuHandler *modbus.RTUClientHandler
var client modbus.Client

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

	if(tcpHandler != nil) {
		err = tcpHandler.Connect()
		defer tcpHandler.Close()
		client = modbus.NewClient(tcpHandler)
	} else {
		err = rtuHandler.Connect()
		defer rtuHandler.Close()
		client = modbus.NewClient(rtuHandler)
	}

	if err == nil {
		executeQueries()
	} else {
		buf.WriteString(fmt.Sprintln(err))
	}

	fmt.Printf(buf.String())

	if config.Enable.LogsEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}
}

func executeSettings() {
	if config.Enable.MbRtuEn {
		rtuHandler = modbus.NewRTUClientHandler(config.MbRtu.Address)
		rtuHandler.BaudRate = config.MbRtu.BaudRate
		rtuHandler.DataBits = config.MbRtu.DataBits
		rtuHandler.StopBits = config.MbRtu.DataBits
		rtuHandler.Parity   = config.MbRtu.Parity
		rtuHandler.SlaveId  = config.MbRtu.SlaveId
		rtuHandler.Timeout  = config.MbRtu.Timeout * time.Second
		buf.WriteString(fmt.Sprintln("ModbusRTU:", config.MbRtu.Address,
				strconv.Itoa(config.MbRtu.BaudRate),
				strconv.Itoa(config.MbRtu.DataBits),
				strconv.Itoa(config.MbRtu.StopBits),
				config.MbRtu.Parity))
		buf.WriteString(fmt.Sprintln("Slave ID :", strconv.Itoa(int(config.MbRtu.SlaveId))))
		buf.WriteString(fmt.Sprintln("Timeout  :", strconv.Itoa(int(config.MbRtu.Timeout)), "second(s)"))
	} else {
		str := config.MbTcp.Ip + ":" + strconv.Itoa(config.MbTcp.Port)
		tcpHandler = modbus.NewTCPClientHandler(str)
		tcpHandler.Timeout = config.MbTcp.Timeout * time.Second
		tcpHandler.SlaveId = config.MbTcp.SlaveId
		//tcpHandler.Logger = log.New(os.Stdout, "test: ", log.LstdFlags)
		buf.WriteString(fmt.Sprintln("ModbusTCP:", str))
		buf.WriteString(fmt.Sprintln("Slave ID :", strconv.Itoa(int(config.MbTcp.SlaveId))))
		buf.WriteString(fmt.Sprintln("Timeout  :", strconv.Itoa(int(config.MbTcp.Timeout)), "second(s)"))
	}
}

func executeQueries() {
	for _, v := range config.MbReg.DiscreteInputs {
		results, err = client.ReadDiscreteInputs(v[0], v[1])
		str := "\nReadDiscreteInputs: " + strconv.Itoa(int(v[0])) + "," + strconv.Itoa(int(v[1]))
		if(err == nil) {
			buf.WriteString(fmt.Sprintln(str, "OK"))
			buf.WriteString(fmt.Sprintln(results))
		} else {
			buf.WriteString(fmt.Sprintln(str, "NG", err))
		}
	}

	for _, v := range config.MbReg.Coils {
		results, err = client.ReadCoils(v[0], v[1])
		str := "\nReadCoils: " + strconv.Itoa(int(v[0])) + "," + strconv.Itoa(int(v[1]))
		if(err == nil) {
			buf.WriteString(fmt.Sprintln(str, "OK"))
			buf.WriteString(fmt.Sprintln(results))
		} else {
			buf.WriteString(fmt.Sprintln(str, "NG", err))
		}
	}

	for _, v := range config.MbReg.HoldingRegisters {
		results, err = client.ReadHoldingRegisters(v[0], v[1])
		str := "\nReadHoldingRegisters: " + strconv.Itoa(int(v[0])) + "," + strconv.Itoa(int(v[1]))
		if(err == nil) {
			buf.WriteString(fmt.Sprintln(str, "OK"))
			buf.WriteString(fmt.Sprintln(results))
		} else {
			buf.WriteString(fmt.Sprintln(str, "NG", err))
		}
	}

	for _, v := range config.MbReg.InputRegisters {
		results, err = client.ReadInputRegisters(v[0], v[1])
		str := "\nReadInputRegisters: " + strconv.Itoa(int(v[0])) + "," + strconv.Itoa(int(v[1]))
		if(err == nil) {
			buf.WriteString(fmt.Sprintln(str, "OK"))
			buf.WriteString(fmt.Sprintln(results))
		} else {
			buf.WriteString(fmt.Sprintln(str, "NG", err))
		}
	}
}

const strDefaultSettings =
`title = "Default Settings for ModbusTCP and ModbusRTU Client Tool"

[enable]
  logsAntoGenEnabled = false
  modbusRTUEnabled = false

[modbusTCP]
  ip   = "127.0.0.1"
  port = 502
  slaveId = 1
  timeout = 3          #Unit: Second

[modbusRTU]
  address  = "COM1"
  baudRate = 2400
  dataBits = 8         #5, 6, 7 or 8
  stopBits = 1         #1 or 2
  parity   = "N"       #"N" - None, "E" - Even, "O" - Odd
  slaveId  = 2
  timeout  = 5         #Unit: Second

[registers]
  discreteInputs   = [[0,  16], [1,   1], [2,   1], [3,   1]]
  coils            = [[0,  64], [1,   1], [2,   1], [3,   1]]
  holdingRegisters = [[0, 100], [1,   3], [2,   1], [3,   1]]
  inputRegisters   = [[0, 100], [1,   3], [2,   1], [3,   1]]
`
