package main

import (
	"flag"
	"fmt"
	"github.com/appcelerator/docker-haproxy/core"
	"os"
)

const version string = "1.0.0-1"

func main() {
	args := os.Args[1:]
	fmt.Println("amp-haproxy-controller started with argument: ", args)
	if len(args) > 0 && args[0] == "autotest" {
		//todo
		os.Exit(0)
	}
	disableBackends := flag.Bool("disable-backends", false, "disable all default backends in haproxy configuration")
	flag.Parse()
	var coreArgs []string
	if *disableBackends {
		coreArgs = append(coreArgs, "disableBackends")
	}
	core.Run(version, coreArgs)
}
