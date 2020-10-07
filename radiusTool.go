//radiusTool.go, Radius Client Tool

package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type tomlConfig struct {
	Title   string    `toml:"title"`
	Enable  enable    `toml:"enable"`
	Server  server    `toml:"server"`
	User    user      `toml:"user"`
}

type enable struct {
	LogEn  bool  `toml:"logAntoGenEnabled"`
	AccEn  bool  `toml:"accountingEnabled"`
}

type server struct {
	Ip      string    `toml:"ip"`
	Port    int       `toml:"port"`
	Secret  string    `toml:"secret"`
}

type user struct {
	Usr  string  `toml:"usr"`
	Pwd  string  `toml:"pwd"`
}

var applName string
var config tomlConfig
var err error
var buf bytes.Buffer
var fileSet string
var fileLog string
var remote  string

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

	remote = config.Server.Ip + ":" + strconv.Itoa(config.Server.Port)
	buf.WriteString(fmt.Sprintln("Remote:", remote))
	buf.WriteString(fmt.Sprintln("Secret:", config.Server.Secret))
	buf.WriteString(fmt.Sprintln("User:", config.User.Usr))
	buf.WriteString(fmt.Sprintln("Pass:", config.User.Pwd))
	buf.WriteString("\n")

	executeSettings()
	fmt.Printf(buf.String())

	if config.Enable.LogEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}
}

func executeSettings() {
	if config.Enable.AccEn {
	}

	packet := radius.New(radius.CodeAccessRequest, []byte(config.Server.Secret))
	rfc2865.UserName_SetString(packet, config.User.Usr)
	rfc2865.UserPassword_SetString(packet, config.User.Pwd)

	response, err := radius.Exchange(context.Background(), packet, remote)
	if err == nil {
		buf.WriteString(fmt.Sprintln("Code:", response.Code))
	} else {
		buf.WriteString(fmt.Sprintln("%v", err))
		log.Printf("%v\n", err)
	}

	log.Println("Code:", response.Code)
}

const strDefaultSettings =
`title = "Default Settings for Radius Client Tool"

[enable]
  logAntoGenEnabled = false
  accountingEnabled = false

[server]
  ip     = "127.0.0.1"
  port   = 1812
  secret = "secret1234"

[user]
  usr = "peter"
  pwd = "Password1234"
`
