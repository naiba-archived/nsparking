package main

import (
	"github.com/naiba/nsparking/cmd/parking/controller"
)

func main() {
	go controller.ServeDNS()
	go controller.ServeWeb()
	select {}
}
