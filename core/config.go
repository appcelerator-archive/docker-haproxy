package core

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//ControllerConfig Json format of conffile
type ControllerConfig struct {
	etcdEndpoints   	[]string
	haProxyPort     	int
	haProxyConffile 	string
	rootServicesKey 	string
	stackName       	string
	stackID         	string
	noDefaultBackend	bool
}

var conf ControllerConfig

//Load Json conffile and instanciate new Config
func (config *ControllerConfig) load(version string) {
	config.setDefault()
	config.loadConfigUsingEnvVariable()
	config.display(version)
}

//Set default value of configuration
func (config *ControllerConfig) setDefault() {
	list := make([]string, 1)
	config.etcdEndpoints = list
	config.etcdEndpoints[0] = "etcd:2379"
	config.haProxyPort = 8082
	config.stackName = ""
	config.stackID = ""
	config.noDefaultBackend = false
}

//Update config with env variables
func (config *ControllerConfig) loadConfigUsingEnvVariable() {
	config.etcdEndpoints = getStringArrayParameter("ETCD_ENDPOINTS", config.etcdEndpoints)
	config.stackName = getStringParameter("STACKNAME", config.stackName)
	config.noDefaultBackend = getBoolParameter("NO_DEFAULT_BACKEND", config.noDefaultBackend)
}

func (config *ControllerConfig) display(version string) {
	fmt.Printf("HAProxy controller version %s\n", version)
	fmt.Printf("ETCD endpoints %s\n", config.etcdEndpoints)
	if config.stackName == "" {
		fmt.Println("Role: Master")
	} else {
		fmt.Printf("Role: Stack: %s\n", config.stackName)
	}
}

//return env variable value, if empty return default value
func getStringParameter(envVariableName string, def string) string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	return value
}

//return env variable value, if empty return default value
func getBoolParameter(envVariableName string, def bool) bool {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	if (strings.ToLower(value) == "true") {
		return true
	}
	return false
}

//return env variable value convert to int, if empty return default value
func getIntParameter(envVariableName string, def int) int {
	value := os.Getenv(envVariableName)
	if value != "" {
		ivalue, err := strconv.Atoi(value)
		if err != nil {
			return def
		}
		return ivalue
	}
	return def
}

//return env variable value, if empty return default value
func getStringArrayParameter(envVariableName string, def []string) []string {
	value := os.Getenv(envVariableName)
	if value == "" {
		return def
	}
	list := strings.Split(strings.Replace(value, " ", "", -1), ",")
	return list
}
