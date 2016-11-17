package main

import (
	"github.com/appcelerator/docker-haproxy/core"
)

const version string = "1.0.2"

func main() {
	core.Run(version)
}
