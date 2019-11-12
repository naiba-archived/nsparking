package main

import (
	"log"

	"github.com/bobesa/go-domain-util/domainutil"
)

func main() {
	var domains = []string{
		"baidu.com.cn",
		"g.com.ms",
		"x.com.ax",
		"a.baidu.com.cn",
	}
	for i := 0; i < len(domains); i++ {
		log.Println(domains[i], domainutil.Domain(domains[i]))
	}
}
