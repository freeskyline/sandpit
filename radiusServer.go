//radiusServer.go, Radius Server Simulator

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
	Ip         string    `toml:"ip"`
	AutPort    int       `toml:"autPort"`
	AccPort    int       `toml:"accPort"`
	AutSecret  string    `toml:"autSecret"`
	AccSecret  string    `toml:"accSecret"`
}

type user struct {
	Allowlist  [][2]string  `toml:"allowlist"`
	Blocklist  [][2]string  `toml:"blocklist"`
}

var applName string
var config tomlConfig
var err error
var buf bytes.Buffer
var fileSet string
var fileLog string

var autServer radius.PacketServer
var accServer radius.PacketServer

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



	if config.Enable.LogEn {
		strTmp := strings.Replace(strTime, ":", "", -1)
		fileLog = applName +"_" + strTmp + ".log"
		ioutil.WriteFile(fileLog, buf.Bytes(), 0644)
	}

	if err = autServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func executeSettings() {
	handler1 := func(w radius.ResponseWriter, r *radius.Request) {
		usr := rfc2865.UserName_GetString(r.Packet)
		pwd := rfc2865.UserPassword_GetString(r.Packet)

		var code radius.Code = radius.CodeAccessReject
		for _, v := range config.User.Allowlist {
			if usr == v[0] && pwd == v[1] {
				code = radius.CodeAccessAccept
			}
		}

		log.Printf("Writing %v to %v %v %v", code, r.RemoteAddr, usr, pwd)
		w.Write(r.Response(code))
	}

	autServer.Addr = config.Server.Ip + ":" +strconv.Itoa(config.Server.AutPort)
	autServer.SecretSource = radius.StaticSecretSource([]byte(config.Server.AutSecret))
	autServer.Handler      = radius.HandlerFunc(handler1)

	buf.WriteString(fmt.Sprintln("Authen radius server listening on:",autServer.Addr))
	buf.WriteString(fmt.Sprintln("Shared secret:", config.Server.AutSecret))

	if config.Enable.AccEn {
	}
}

const strDefaultSettings =
`title = "Default Settings for Radius Server Simulator"

[enable]
  logAntoGenEnabled = false
  accountingEnabled = false

[server]
  ip = "0.0.0.0"
  autPort = 1812
  accPort = 1813
  autSecret = "secret1234"
  accSecret = "secret5678"

[user]
  allowlist = [["peter", "Password1234"],
               ["lihui", "Password5678"]]

  blocklist = [["PETER", "Password1234"],
               ["lihui", "Password5678"]]
`
