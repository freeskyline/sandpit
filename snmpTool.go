//snmpTool.go, SNMP Client Tool

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

	"github.com/alouca/gosnmp"
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Title   string    `toml:"title"`
	Enable  enable    `toml:"enable"`
	Target  target    `toml:"target"`
	Oid     oid       `toml:"oid"`
}

type enable struct {
	DbgEn  bool  `toml:"debugDataEnabled"`
	LogEn  bool  `toml:"logAntoGenEnabled"`
}

type target struct {
	Ip         string         `toml:"ip"`
	Port       int            `toml:"port"`
	Timeout    int64          `toml:"timeout"`
	Community  string         `toml:"community"`
	Version    string         `toml:"version"`
}

type oid struct {
	snmpget []string    `toml:"singleCoil"`
	getnext []string    `toml:"singleRegister"`
}

var applName string
var config  tomlConfig
var err error
var results []byte
var buf bytes.Buffer
var fileSet string
var fileLog string

var snmp *gosnmp.GoSNMP

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

	if err == nil {
		executeQueries()
	} else {
		buf.WriteString(fmt.Sprintln(err))
	}

	fmt.Printf(buf.String())

	if config.Enable.LogEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}
}

func executeSettings() {
	str := config.Target.Ip + ":" + strconv.Itoa(config.Target.Port)
	snmp, err = gosnmp.NewGoSNMP(str, 
				config.Target.Community,
				transSNMPVer(config.Target.Version),
				config.Target.Timeout)

	if err == nil {
		snmp.SetDebug(config.Enable.DbgEn)
		snmp.SetVerbose(false)
	} else {
		buf.WriteString(fmt.Sprintln(err))
	}

	buf.WriteString(fmt.Sprintln("Target:", str))
	buf.WriteString(fmt.Sprintln("Community :", config.Target.Community))
	buf.WriteString(fmt.Sprintln("Version :",   config.Target.Version))
	buf.WriteString(fmt.Sprintln("Timeout :", strconv.Itoa(int(config.Target.Timeout)), "second(s)"))
}

func transSNMPVer(ver string) (sv gosnmp.SnmpVersion) {
	sv = gosnmp.Version1

	if ver == "2c" {
		sv = gosnmp.Version2c
	}

	return
}

func executeQueries() {
	for _, oid := range config.Oid.snmpget {
		r, e := snmp.Get(oid)
		str := "\nGet OID: " + oid
		if e == nil {
			buf.WriteString(fmt.Sprintln(str, "OK"))
			
			for _, v := range r.Variables {
			fmt.Printf("%s -> ", v.Name)
				switch v.Type {
				case gosnmp.OctetString:
					if s, ok := v.Value.(string); ok {
						fmt.Printf("%s\n", s)
					} else {
						fmt.Printf("Response is not a string\n")
					}
				default:
					fmt.Printf("Type: %d - Value: %v\n", v.Type, v.Value)
				}
			}
			buf.WriteString(fmt.Sprintln(results))
		} else {
			buf.WriteString(fmt.Sprintln(str, "NG", err))
		}
	}
}

const strDefaultSettings =
`title = "Default Settings for SNMP Client Tool"

[enable]
  debugDataEnabled  = false
  logAntoGenEnabled = false

[target]
  ip   = "127.0.0.1"
  port = 502
  timeout = 3          #Unit: Second
  community = "public"
  version   = "2c"    #"1","2c"

[oid]
  snmpget = [".1.3.6.1.2.1.1.1.0",
             ".1.3.6.1.2.1.1.1.1",
             ".1.3.6.1.2.1.1.1.2",
             ".1.3.6.1.2.1.1.1.3"]

  getnext = [".1.3.6.1.1",
             ".1.3.6.1.4"]
`
