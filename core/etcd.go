package core

import (
	"fmt"
	"github.com/appcelerator/amp/api/rpc/service"
	"github.com/appcelerator/amp/api/rpc/stack"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"os"
	"path"
	"strings"
	"time"
)

const stackRootKey = "amp/stacks"
const stackRootNameKey = "amp/stacks/names"
const serviceRootNameKey = "amp/services"

//ETCDClient etcd struct
type ETCDClient struct {
	client *clientv3.Client
}

var etcdClient ETCDClient

func (inst *ETCDClient) init() error {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.etcdEndpoints,
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Connected to etcd: %v\n", conf.etcdEndpoints)
	inst.client = cli
	if conf.stackName != "" {
		conf.stackID = inst.getStackByName(conf.stackName)
		if conf.stackID == "" {
			fmt.Printf("Stack %s not found\n", conf.stackName)
			inst.Close()
			os.Exit(1)
		}
	}
	return err
}

//Close close ETCD client
func (inst *ETCDClient) Close() error {
	if err := inst.client.Close(); err != nil {
		return err
	}
	inst.client = nil
	return nil
}

// Get stack Id using stack name
func (inst *ETCDClient) getStackByName(name string) string {
	stackID := stack.StackID{}
	resp, _ := inst.client.Get(context.Background(), path.Join(stackRootNameKey, name))
	if len(resp.Kvs) == 0 {
		return stackID.Id
	}
	proto.Unmarshal(resp.Kvs[0].Value, &stackID)
	return stackID.Id
}

func (inst *ETCDClient) watchForServicesUpdate() {
	watchKeys := stackRootNameKey
	if conf.stackName != "" {
		watchKeys = path.Join(stackRootKey, conf.stackID, "services")
	}
	fmt.Println("Waiting for update on ", watchKeys)
	rch := inst.client.Watch(context.Background(), watchKeys, clientv3.WithPrefix())
	wresp := <-rch
	for _, ev := range wresp.Events {
		fmt.Printf("Key updated: %q\n", ev.Kv.Key)
	}
	time.Sleep(5 * time.Second)
	haproxy.loadTry = 0
	haproxy.updateConfiguration(true)
}

func (inst *ETCDClient) getAllPublicServices(stackID string) (map[string]*publicService, error) {
	resp, err := inst.client.Get(context.Background(), path.Join(stackRootKey, conf.stackID, "services"))
	if err != nil {
		return nil, err
	}
	serviceMap := make(map[string]*publicService)
	if len(resp.Kvs) == 0 {
		return serviceMap, nil
	}
	servList := stack.IdList{}
	erru := proto.Unmarshal(resp.Kvs[0].Value, &servList)
	if erru != nil {
		return nil, erru
	}
	for _, servID := range servList.List {
		servSpec := service.ServiceSpec{}
		fmt.Printf("process service Id:%s\n", servID)
		resp, err := inst.client.Get(context.Background(), path.Join(serviceRootNameKey, servID))
		if err != nil {
			return nil, err
		}
		merr := proto.Unmarshal(resp.Kvs[0].Value, &servSpec)
		if merr != nil {
			fmt.Printf("Error unmarcharling service id=%s error=%v\n", servID, merr)
		} else {
			if servSpec.PublishSpecs != nil && len(servSpec.PublishSpecs) > 0 {
				pService := publicService{
					name:    servSpec.Name,
					mapping: make(map[string]string),
				}
				serviceMap[servSpec.Name] = &pService
				for _, mapping := range servSpec.PublishSpecs {
					if mapping.Name != "" && mapping.InternalPort != 0 {
						pService.mapping[mapping.Name] = fmt.Sprintf("%d", mapping.InternalPort)
						fmt.Printf("add service %s mapping: %s\n", servSpec.Name, mapping.Name)
					}
				}
			}
		}
	}
	return serviceMap, nil
}

func (inst *ETCDClient) getAllStacks() (map[string]*publicStack, error) {
	resp, err := inst.client.Get(context.Background(), stackRootNameKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	stackMap := make(map[string]*publicStack)
	for _, st := range resp.Kvs {
		fmt.Printf("key:%s\n", string(st.Key))
		list := strings.Split(string(st.Key), "/")
		stackName := list[len(list)-1]
		stackMap[stackName] = &publicStack{
			name: stackName,
		}
	}
	return stackMap, nil
}
