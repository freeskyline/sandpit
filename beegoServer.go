// Beego Web Server

package main

import (
	"fmt"
	"github.com/astaxie/beego"
)

func banner() func() {
	const (
		strApp = "beegoServer"
		strVer = "v1.0.0"
		strUTC = "Build LIHUI 2020-09-09 22:51:33.302 +0000 UTC"
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

func main() {
	defer banner()()
	beego.Run()
}
