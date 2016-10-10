package core

import (
	"fmt"
	"os"
	"time"
)

var ampHAProxyControllerVersion string

//Run launch main loop
func Run(version string) {
	ampHAProxyControllerVersion = version
	conf.load(version)
	err := etcdClient.init()
	if err != nil {
		fmt.Printf("ETCD connection error %v\n", err)
		os.Exit(1)
	}
	haproxy.init()
	haproxy.trapSignal()
	initAPI()
	haproxy.start()
	for {
		time.Sleep(3 * time.Second)
		etcdClient.watchForServicesUpdate()
	}
}
