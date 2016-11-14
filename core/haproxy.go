package core

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//HAProxy haproxy struct
type HAProxy struct {
	exec                *exec.Cmd
	isLoadingConf       bool
	dnsRetryLoopId      int
	loadTry             int
	dnsNotResolvedList  []string
	defaultInfraService []defaultInfraService
}

type publicService struct {
	name    string
	mapping map[string]string
}

type publicStack struct {
	name string
}

type defaultInfraService struct {
	name string
	port int
	mode string
}

var (
	haproxy HAProxy
)

//Set app mate initial values
func (app *HAProxy) init() {
	app.isLoadingConf = false
	app.loadTry = 0
	app.dnsNotResolvedList = []string{}

	app.defaultInfraService = []defaultInfraService{
		defaultInfraService{name: "amplifier", port: 50101, mode: "tcp"},
		defaultInfraService{name: "grafana", port: 3000},
		defaultInfraService{name: "elasticsearch", port: 9200},
		defaultInfraService{name: "amp-ui", port: 8080},
		defaultInfraService{name: "registry", port: 5000},
		defaultInfraService{name: "amplifier-gateway", port: 3000},
	}
	haproxy.updateConfiguration(false)
}

//Launch a routine to catch SIGTERM Signal
func (app *HAProxy) trapSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)
	go func() {
		<-ch
		fmt.Println("\namp-haproxy-controller received SIGTERM signal")
		etcdClient.Close()
		os.Exit(1)
	}()
}

//Launch HAProxy using cmd command
func (app *HAProxy) start() {
	go func() {
		fmt.Println("launching HAProxy on initial configuration")
		app.exec = exec.Command("haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg")
		app.exec.Stdout = os.Stdout
		app.exec.Stderr = os.Stderr
		err := app.exec.Run()
		if err != nil {
			fmt.Printf("HAProxy exit with error: %v\n", err)
			etcdClient.Close()
			os.Exit(1)
		}
	}()
}

//Stop HAProxy
func (app *HAProxy) stop() {
	fmt.Println("Send SIGTERM signal to HAProxy")
	if app.exec != nil {
		app.exec.Process.Kill()
	}
}

//Launch HAProxy using cmd command
func (app *HAProxy) reloadConfiguration() {
	app.isLoadingConf = true
	fmt.Println("reloading HAProxy configuration")
	pid := app.exec.Process.Pid
	fmt.Printf("Execute: %s %s %s %s %d\n", "haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-sf", pid)
	app.exec = exec.Command("haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-sf", fmt.Sprintf("%d", pid))
	app.exec.Stdout = os.Stdout
	app.exec.Stderr = os.Stderr
	go func() {
		err := app.exec.Run()
		app.isLoadingConf = false
		if err == nil {
			fmt.Println("HAProxy configuration reloaded")
			return
		}
		app.loadTry++
		fmt.Printf("HAProxy reload configuration error, try=%s: %v\n", app.loadTry, err)
		if app.loadTry > 6 {
			os.Exit(1)
		}
		time.Sleep(10 * time.Second)
		app.updateConfiguration(true)
	}()
}

// update configuration managing the isUpdatingConf flag
func (app *HAProxy) updateConfiguration(reload bool) error {
	app.dnsRetryLoopId++
	app.dnsNotResolvedList = []string{}
	err := app.updateConfigurationEff(reload)
	if err == nil {
		app.startDNSRevolverLoop(app.dnsRetryLoopId)
	}
	return err
}

//update HAProxy configuration for master regarding ETCD keys values and make HAProxy reload its configuration if reload is true
func (app *HAProxy) updateConfigurationEff(reload bool) error {
	fmt.Println("update HAProxy configuration")
	var listStack map[string]*publicStack
	var listService map[string]*publicService
	if conf.stackName == "" {
		list, err := etcdClient.getAllStacks()
		if err != nil {
			fmt.Println("Erreur on get stacks list: ", err)
			return err
		}
		listStack = list
	} else {
		list, err := etcdClient.getAllPublicServices(conf.stackName)
		if err != nil {
			fmt.Println("Erreur on get services list: ", err)
			return err
		}
		listService = list
	}
	fileNameTarget := "/usr/local/etc/haproxy/haproxy.cfg"
	fileNameTpt := "/usr/local/etc/haproxy/haproxy.cfg.tpt"
	file, err := os.Create(fileNameTarget + ".new")
	if err != nil {
		fmt.Printf("Error creating new haproxy conffile for creation: %v\n", err)
		return err
	}
	filetpt, err := os.Open(fileNameTpt)
	if err != nil {
		fmt.Printf("Error opening conffile template: %s : %v\n", fileNameTpt, err)
		return err
	}
	scanner := bufio.NewScanner(filetpt)
	skip := false
	for scanner.Scan() {
		line := scanner.Text()
		skip = hasToBeSkipped(line, skip)
		if conf.debug {
			fmt.Printf("line: %t: %s\n", skip, line)
		}
		if !skip {
			if strings.HasPrefix(strings.Trim(line, " "), "[frontend]") {
				if conf.stackName == "" {
					app.writeStackFrontend(file, listStack)
				} else {
					app.writeServiceFrontend(file, listService)
				}
			} else if strings.HasPrefix(strings.Trim(line, " "), "[backends]") {
				if conf.stackName == "" {
					app.writeStackBackend(file, listStack)
				} else {
					app.writeServiceBackend(file, listService)
				}
			} else {
				file.WriteString(line + "\n")
			}
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading haproxy conffile template: %s %v\n", fileNameTpt, err)
		file.Close()
		return err
	}
	if conf.stackName == "" && !conf.noDefaultBackend {
		for _, serv := range app.defaultInfraService {
			app.writeInfraServiceBackend(file, serv)
		}
	}
	file.Close()
	os.Remove(fileNameTarget)
	err2 := os.Rename(fileNameTarget+".new", fileNameTarget)
	if err2 != nil {
		fmt.Printf("Error renaming haproxy conffile .new: %v\n", err)
		return err
	}
	fmt.Println("HAProxy configuration updated")
	if reload {
		app.reloadConfiguration()
	}
	return nil
}

// compute if the line is the begining of a block which should be skipped or not
func hasToBeSkipped(line string, skip bool) bool {
	if line == "" {
		//if blanck line then end of skip
		return false
	} else if skip {
		// if skipped mode and line not black then continue to skip
		return true
	}

	ref := strings.Trim(line, " ")
	if conf.noDefaultBackend {
		if strings.HasPrefix(ref, "frontend main_") || strings.HasPrefix(ref, "backend main_") || strings.HasPrefix(ref, "backend stack_") {
			//if nodefaultbackend activated and default frontend or backend then skip
			return true
		}
	} else {
		if conf.stackName != "" {
			if strings.HasPrefix(ref, "frontend main_") || strings.HasPrefix(ref, "backend main_") {
				//if stack mode and main frontend or backend then skip
				return true
			}
		} else {
			if strings.HasPrefix(ref, "frontend stack_") || strings.HasPrefix(ref, "backend stack_") {
				//if main mode and stack frontend or backend then skip
				return true
			}
		}
	}
	//if line not "" and not skip then continue to not skip
	return false
}

// write backends for main service configuration
func (app *HAProxy) writeServiceFrontend(file *os.File, serviceMap map[string]*publicService) error {
	for _, service := range serviceMap {
		for extName, intPort := range service.mapping {
			line := fmt.Sprintf("    use_backend bk_%s%s if { hdr_beg(host) -i %s. }\n", service.name, intPort, extName)
			file.WriteString(line)
			fmt.Printf(line)
		}

	}
	return nil
}

// write backends for main haproxy configuration
func (app *HAProxy) writeStackFrontend(file *os.File, stackMap map[string]*publicStack) error {
	for _, stack := range stackMap {
		line := fmt.Sprintf("    use_backend bk_%s if { hdr_dom(host) -i .%s. }\n", stack.name, stack.name)
		file.WriteString(line)
		fmt.Printf(line)
	}
	return nil
}

// write backends for stack haproxy configuration
func (app *HAProxy) writeServiceBackend(file *os.File, serviceMap map[string]*publicService) error {
	for _, service := range serviceMap {
		for _, intPort := range service.mapping {
			dnsResolved := app.tryToResolvDNS(service.name)
			line1 := fmt.Sprintf("\nbackend bk_%s%s\n", service.name, intPort)
			file.WriteString(line1)
			fmt.Printf(line1)
			//if dns name is not resolved haproxy (v1.6) won't start or accept the new configuration so server is disabled
			//to be removed when haproxy will fixe this bug
			if dnsResolved {
				line2 := fmt.Sprintf("    server %s_1 %s:%s check resolvers docker resolve-prefer ipv4\n", service.name, service.name, intPort)
				file.WriteString(line2)
				fmt.Printf(line2)
			} else {
				line2 := "    #dns name not resolved\n"
				file.WriteString(line2)
				fmt.Printf(line2)
				line3 := fmt.Sprintf("    #server %s_1 %s:%s check resolvers docker resolve-prefer ipv4\n", service.name, service.name, intPort)
				file.WriteString(line3)
				fmt.Printf(line3)
				app.addDNSNameInRetryList(service.name)
			}
		}
	}
	return nil
}

// write infra backends for main haproxy configuration
func (app *HAProxy) writeInfraServiceBackend(file *os.File, service defaultInfraService) error {
	line1 := fmt.Sprintf("\nbackend infra_%s\n", service.name)
	file.WriteString(line1)
	fmt.Printf(line1)
	//if dns name is not resolved haproxy (v1.6) won't start or accept the new configuration so server is disabled
	//to be removed when haproxy will fixe this bug
	if app.tryToResolvDNS(service.name) {
		if service.mode != "" {
			file.WriteString("    mode " + service.mode + "\n")
		}
		line2 := fmt.Sprintf("    server %s_1 %s:%d check resolvers docker resolve-prefer ipv4\n", service.name, service.name, service.port)
		file.WriteString(line2)
		fmt.Printf(line2)
	} else {
		line2 := "    #dns name not resolved\n"
		file.WriteString(line2)
		fmt.Printf(line2)
		line3 := fmt.Sprintf("    #server %s_1 %s:%d check resolvers docker resolve-prefer ipv4\n", service.name, service.name, service.port)
		file.WriteString(line3)
		fmt.Printf(line3)
		app.addDNSNameInRetryList(service.name)
	}
	return nil
}

// write backends for main haproxy configuration
func (app *HAProxy) writeStackBackend(file *os.File, stackMap map[string]*publicStack) error {
	for _, stack := range stackMap {
		line1 := fmt.Sprintf("\nbackend bk_%s\n", stack.name)
		file.WriteString(line1)
		fmt.Printf(line1)
		//if dns name is not resolved haproxy (v1.6) won't start or accept the new configuration so server is disabled
		//to be removed when haproxy will fixe this bug
		if app.tryToResolvDNS(fmt.Sprintf("%s-haproxy", stack.name)) {
			line2 := fmt.Sprintf("    server %s_1 %s-haproxy:80 check resolvers docker resolve-prefer ipv4\n", stack.name, stack.name)
			file.WriteString(line2)
			fmt.Printf(line2)
		} else {
			line2 := "    #dns name not resolved\n"
			file.WriteString(line2)
			fmt.Println(line2)
			line3 := fmt.Sprintf("    #server %s_1 %s-haproxy:80 check resolvers docker resolve-prefer ipv4\n", stack.name, stack.name)
			file.WriteString(line3)
			fmt.Printf(line3)
			app.addDNSNameInRetryList(fmt.Sprintf("%s-haproxy", stack.name))
		}
	}
	return nil
}

// test if a dns name is resolved or not
func (app *HAProxy) tryToResolvDNS(name string) bool {
	_, err := net.LookupIP(name)
	if err != nil {
		return false
	}
	return true
}

// add unresolved dns name in list to be retested later
func (app *HAProxy) addDNSNameInRetryList(name string) {
	app.dnsNotResolvedList = append(app.dnsNotResolvedList, name)
}

// on regular basis try to see if one of the unresolved dns become resolved, if so execute a configuration update.
// need to have only one loop at a time, if the id change then the current loop should stop
// id is incremented at each configuration update which can be trigger by ETCD wash also
func (app *HAProxy) startDNSRevolverLoop(loopId int) {
	//if no unresolved DNS name then not needed to start the loop
	if len(haproxy.dnsNotResolvedList) == 0 {
		return
	}
	fmt.Printf("Start DNS resolver id: %d\n", loopId)
	go func() {
		for {
			for _, name := range haproxy.dnsNotResolvedList {
				if app.tryToResolvDNS(name) {
					if haproxy.dnsRetryLoopId == loopId {
						fmt.Printf("DNS %s resolved, update configuration\n", name)
						app.updateConfiguration(true)
					}
					fmt.Printf("Stop DNS resolver id: %d\n", loopId)
					return
				}
			}
			time.Sleep(10)
			if haproxy.dnsRetryLoopId != loopId {
				fmt.Printf("Stop DNS resolver id: %d\n", loopId)
				return
			}
		}
	}()
}
