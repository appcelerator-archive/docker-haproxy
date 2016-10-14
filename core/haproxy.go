package core

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

//HAProxy haproxy struct
type HAProxy struct {
	exec          *exec.Cmd
	isLoadingConf bool
	loadTry       int
}

type publicService struct {
	name    string
	mapping map[string]string
}

type publicStack struct {
	name string
}

var (
	haproxy HAProxy
)

//Set app mate initial values
func (app *HAProxy) init() {
	app.isLoadingConf = false
	app.loadTry = 0
	if conf.noDefaultBackend {
		fmt.Println("Default backends are disabled by request")
		err := app.disableDefaultBackends("/usr/local/etc/haproxy/haproxy.cfg")
		if err != nil {
			log.Fatalf("Failed to disable backends in %s: %v\n", "/usr/local/etc/haproxy/haproxy.cfg", err)
		}
		err = app.disableDefaultBackends("/usr/local/etc/haproxy/haproxy-main.cfg.tpt")
		if err != nil {
			log.Fatalf("Failed to disable backends in %s: %v\n", "/usr/local/etc/haproxy/haproxy-main.cfg.tpt", err)
		}
	}
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
			fmt.Printf("HAProxy configuration reloaded")
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

//update HAProxy configuration regarding ETCD keys values and make HAProxy reload its configuration if reload is true
func (app *HAProxy) updateConfiguration(reload bool) error {
	if conf.stackName == "" {
		//pp.executeDNSPatch()
		return app.updateConfigurationMaster(reload)
	}
	return app.updateConfigurationStack(reload)
}

// disable a backend in the configuration file
// useful for tests purpose when the full stack is not deployed
func (app *HAProxy) disableDefaultBackends(target string) error {
	file, err := os.Open(target)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	var bufferedLine []byte

	fmt.Printf("Updating %s...\n", target)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match, err := regexp.MatchString(`^\s+server\s+`, line)
		if err != nil {
			return err
		}
		matchLocalhost, err := regexp.MatchString(`\s+localhost:`, line)
		if err != nil {
			return err
		}
		if match && !matchLocalhost {
			fmt.Printf("Disabling %s\n", line)
			bufferedLine = []byte(`#` + line + "\n")
			if _, err = buf.Write(bufferedLine); err != nil {
				return err
			}
		} else {
			bufferedLine = []byte(line + "\n")
			if _, err = buf.Write(bufferedLine); err != nil {
				return err
			}
		}
	}
	if err = scanner.Err(); err != nil {
		return err
	}
	if err = file.Close(); err != nil {
		return err
	}

	fmt.Println("Writing...")
	file, err = os.Create(target)
	if err != nil {
		return err
	}
	if _, err = buf.WriteTo(file); err != nil {
		_ = file.Close()
		return err
	}
	fmt.Println("done")
	err = file.Close()
	return err
}

//update HAProxy configuration for master regarding ETCD keys values and make HAProxy reload its configuration if reload is true
func (app *HAProxy) updateConfigurationMaster(reload bool) error {
	fmt.Println("update HAProxy configuration")
	list, err := etcdClient.getAllStacks()
	if err != nil {
		fmt.Println("Erreur on get stacks list: ", err)
	}
	fileNameTarget := "/usr/local/etc/haproxy/haproxy.cfg"
	fileNameTpt := "/usr/local/etc/haproxy/haproxy-main.cfg.tpt"
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
	for scanner.Scan() {
		line := scanner.Text()
		if line == "[frontend]" {
			app.writeStackFrontend(file, list)
		} else if line == "[backends]" {
			app.writeStackBackend(file, list)
		} else {
			file.WriteString(line + "\n")
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading haproxy conffile template: %s %v\n", fileNameTpt, err)
		file.Close()
		return err
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

//update HAProxy configuration for stack regarding ETCD keys values and make HAProxy reload its configuration if reload is true
func (app *HAProxy) updateConfigurationStack(reload bool) error {
	fmt.Println("update HAProxy configuration")
	list, err := etcdClient.getAllPublicServices(conf.stackName)
	if err != nil {
		fmt.Println("Erreur on get services list: ", err)
	}
	fileNameTarget := "/usr/local/etc/haproxy/haproxy.cfg"
	fileNameTpt := "/usr/local/etc/haproxy/haproxy-stack.cfg.tpt"
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
	for scanner.Scan() {
		line := scanner.Text()
		if line == "[frontend]" {
			app.writeServiceFrontend(file, list)
		} else if line == "[backends]" {
			app.writeServiceBackend(file, list)
		} else {
			file.WriteString(line + "\n")
		}
	}
	if err = scanner.Err(); err != nil {
		fmt.Printf("Error reading haproxy conffile template: %s %v\n", fileNameTpt, err)
		file.Close()
		return err
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

func (app *HAProxy) writeServiceFrontend(file *os.File, serviceMap map[string]*publicService) error {
	fmt.Println("Update frontend")
	for _, service := range serviceMap {
		for extName, intPort := range service.mapping {
			line := fmt.Sprintf("use_backend bk_%s%s if { hdr_beg(host) -i %s. }\n", service.name, intPort, extName)
			file.WriteString("\t" + line)
			fmt.Printf(line)
		}

	}
	return nil
}

func (app *HAProxy) writeServiceBackend(file *os.File, serviceMap map[string]*publicService) error {
	fmt.Println("Update backends")
	for _, service := range serviceMap {
		for _, intPort := range service.mapping {
			file.WriteString("\n")
			line1 := fmt.Sprintf("\nbackend bk_%s%s\n", service.name, intPort)
			file.WriteString(line1)
			fmt.Printf(line1)
			line2 := fmt.Sprintf("server %s_1 %s:%s check resolvers docker resolve-prefer ipv4\n", service.name, service.name, intPort)
			file.WriteString("\t" + line2)
			fmt.Printf(line2)
		}
	}
	file.WriteString("\n")
	return nil
}

func (app *HAProxy) writeStackFrontend(file *os.File, stackMap map[string]*publicStack) error {
	fmt.Println("Update frontend")
	for _, stack := range stackMap {
		line := fmt.Sprintf("use_backend bk_%s if { hdr_dom(host) -i .%s. }\n", stack.name, stack.name)
		file.WriteString("\t" + line)
		fmt.Printf(line)
	}
	return nil
}

func (app *HAProxy) writeStackBackend(file *os.File, stackMap map[string]*publicStack) error {
	fmt.Println("Update backends")
	for _, stack := range stackMap {
		file.WriteString("\n")
		line1 := fmt.Sprintf("\nbackend bk_%s\n", stack.name)
		file.WriteString(line1)
		fmt.Printf(line1)
		line2 := fmt.Sprintf("server %s_1 %s-haproxy:80 check resolvers docker resolve-prefer ipv4\n", stack.name, stack.name)
		file.WriteString("\t" + line2)
		fmt.Printf(line2)
	}
	file.WriteString("\n")
	return nil
}
