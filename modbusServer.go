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
	Title   string    `toml:"title"`
	Enable  enable    `toml:"enable"`
	MbTcp   modbusTCP `toml:"modbusTCP"`
	MbRtu   modbusRTU `toml:"modbusRTU"`
	MbReg   registers `toml:"registers"`
	MbSim   simulated `toml:"simulated"`
}

type enable struct {
	LogsEn   bool  `toml:"logAntoGenEnabled"`
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

type registers struct {
	DiscreteInputs    [][2]uint32  `toml:"discreteInputs"`
	Coils             [][2]uint32  `toml:"coils"`
	InputRegisters    [][2]uint32  `toml:"inputRegisters"`
	HoldingRegisters  [][2]uint32  `toml:"holdingRegisters"`
}

type simulated struct {
	DataDiscreteInputs    [][2][]uint32  `toml:"dataDiscreteInputs"`
	DataCoils             [][2][]uint32  `toml:"dataCoils"`
	DataInputRegisters    [][2][]uint32  `toml:"dataInputRegisters"`
	DataHoldingRegisters  [][2][]uint32  `toml:"dataHoldingRegisters"`
}

var applName string
var config  tomlConfig
var server *mbserver.Server
var err error
var buf bytes.Buffer
var fileSet string
var fileLog string

var mapDiscreteInputs    map[uint32][]uint32
var mapCoils             map[uint32][]uint32
var mapInputRegisters    map[uint32][]uint32
var mapHoldingRegisters  map[uint32][]uint32
var lenDiscreteInputs    map[uint32]uint32
var lenCoils             map[uint32]uint32
var lenInputRegisters    map[uint32]uint32
var lenHoldingRegisters  map[uint32]uint32
var idxDiscreteInputs    map[uint32]uint32
var idxCoils             map[uint32]uint32
var idxInputRegisters    map[uint32]uint32
var idxHoldingRegisters  map[uint32]uint32

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

	server = mbserver.NewServer()
	initModbusServer(server)
	defer server.Close()

	executeSettings()
	fmt.Printf(buf.String())

	if config.Enable.LogsEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}

	for {
		updateModbusServer(server)
		time.Sleep(time.Second)
	}
}

func executeSettings() {
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
			buf.WriteString(fmt.Sprintln("ModbusRTU Server listening on: ",
				cnf.Address,
				strconv.Itoa(cnf.BaudRate),
				strconv.Itoa(cnf.DataBits),
				strconv.Itoa(cnf.StopBits),
				cnf.Parity,
				strconv.Itoa(int(cnf.Timeout / time.Second))))
		} else {
			log.Printf("%v\n", err)
			buf.WriteString(fmt.Sprintln("%v", err))
		}
	}

	mapDiscreteInputs = make(map[uint32][]uint32)
	lenDiscreteInputs = make(map[uint32]uint32)
	idxDiscreteInputs = make(map[uint32]uint32)
	for _, v := range config.MbSim.DataDiscreteInputs {
		mapDiscreteInputs[v[0][0]] = v[1]
		lenDiscreteInputs[v[0][0]] = uint32(len(v[1]))
		idxDiscreteInputs[v[0][0]] = 0
	}

	mapCoils = make(map[uint32][]uint32)
	lenCoils = make(map[uint32]uint32)
	idxCoils = make(map[uint32]uint32)
	for _, v := range config.MbSim.DataCoils {
		mapCoils[v[0][0]] = v[1]
		lenCoils[v[0][0]] = uint32(len(v[1]))
		idxCoils[v[0][0]] = 0
	}

	mapInputRegisters = make(map[uint32][]uint32)
	lenInputRegisters = make(map[uint32]uint32)
	idxInputRegisters = make(map[uint32]uint32)
	for _, v := range config.MbSim.DataInputRegisters {
		mapInputRegisters[v[0][0]] = v[1]
		lenInputRegisters[v[0][0]] = uint32(len(v[1]))
		idxInputRegisters[v[0][0]] = 0
	}

	mapHoldingRegisters = make(map[uint32][]uint32)
	lenHoldingRegisters = make(map[uint32]uint32)
	idxHoldingRegisters = make(map[uint32]uint32)
	for _, v := range config.MbSim.DataHoldingRegisters {
		mapHoldingRegisters[v[0][0]] = v[1]
		lenHoldingRegisters[v[0][0]] = uint32(len(v[1]))
		idxHoldingRegisters[v[0][0]] = 0
	}
}

func initModbusServer(s *mbserver.Server) {
	for _, v := range config.MbReg.DiscreteInputs {
		s.DiscreteInputs[v[0]] = byte(v[1])
	}

	for _, v := range config.MbReg.Coils {
		s.Coils[v[0]] = byte(v[1])
	}

	for _, v := range config.MbReg.InputRegisters {
		s.InputRegisters[v[0]] = uint16(v[1])
	}

	for _, v := range config.MbReg.HoldingRegisters {
		s.HoldingRegisters[v[0]] = uint16(v[1])
	}
}

func updateModbusServer(s *mbserver.Server) {
	for k, v := range idxDiscreteInputs{
		s.DiscreteInputs[k] = byte(mapDiscreteInputs[k][v])
	}
	for k, v := range idxDiscreteInputs {
		tmp := v + 1
		tmp %= lenDiscreteInputs[k]
		idxDiscreteInputs[k] = tmp
	}

	for k, v := range idxCoils {
		s.Coils[k] = byte(mapCoils[k][v])
	}
	for k, v := range idxCoils {
		tmp := v + 1
		tmp %= lenCoils[k]
		idxCoils[k] = tmp
	}

	for k, v := range idxInputRegisters {
		s.InputRegisters[k] = uint16(mapInputRegisters[k][v])
	}
	for k, v := range idxInputRegisters {
		tmp := v + 1
		tmp %= lenInputRegisters[k]
		idxInputRegisters[k] = tmp
	}

	for k, v := range idxHoldingRegisters {
		s.HoldingRegisters[k] = uint16(mapHoldingRegisters[k][v])
	}
	for k, v := range idxHoldingRegisters {
		tmp := v + 1
		tmp %= lenHoldingRegisters[k]
		idxHoldingRegisters[k] = tmp
	}
}

const strDefaultSettings =
`title = "Default Settings for ModbusTCP and ModbusRTU Server Simulator"

[enable]
  logAntoGenEnabled = false
  modbusTCPEnabled = true
  modbusRTUEnabled = false

[modbusTCP]
  ip   = "0.0.0.0"
  port = 502

[modbusRTU]
  address  = "COM1"
  baudRate = 2400
  dataBits = 8         #5, 6, 7 or 8
  stopBits = 1         #1 or 2
  parity   = "N"       #"N" - None, "E" - Even, "O" - Odd
  timeout  = 10        #Unit: Second

[registers]
  discreteInputs   = [[0,     6], [1,     1], [2,     2], [3,     3]]
  coils            = [[0,    60], [1,    10], [2,    20], [3,    30]]
  inputRegisters   = [[0,   600], [1,   100], [2,   200], [3,   300]]
  holdingRegisters = [[0, 60000], [1, 10000], [2, 20000], [3, 30000]]

[simulated]
  dataDiscreteInputs   = [[[100], [0,1]],
                          [[101], [0,1,1,0]],
                          [[102], [0,1,1,0,0,0]],
                          [[103], [0,1,1,0,1,0,0,1,1,1,1]]]

  dataCoils            = [[[200], [0,0,0,0,0,0]],
                          [[201], [0,1,1,0]],
                          [[202], [0,1,1,0,0,0]],
                          [[203], [0,1,1,0,1,0,0,1,1,1]]]
  dataInputRegisters   = [[[300], [ 0,1]],
                          [[301], [10,11,12,13]],
                          [[302], [20,21,22,23,24,25]],
                          [[303], [30,31,32,33,34,35,36,37,38,39]]]

  dataHoldingRegisters = [[[400], [ 0,99]],
                          [[401], [10,11,12,13]],
                          [[402], [20,21,22,23,24,25]],
                          [[403], [30,31,32,33,34,35,36,37,38,39]]]
`
