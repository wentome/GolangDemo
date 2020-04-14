package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"../../cdeamon"
	"../../ctask"
	"./define"
	"./dog"
	"./event"
	"./grpc"
	"github.com/astaxie/beego/logs"
)

func init() {
	// info 级别
	//logs.SetLogger(logs.AdapterConsole, `{"level":6,"color":true}`)
	//默认7 debug 级别
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/watch.log","daily":true,"maxdays":10}`)
	logs.Async()
}
func main() {
	var serverMode bool
	flag.BoolVar(&serverMode, "server", false, "server mode")
	clusterHostsP := flag.String("cluster", "", "host1:host2:host3")
	flag.Parse()
	flagArgs := flag.Args()

	//client mode
	if !serverMode {
		dog.Client()
		return
	}

	//server mode
	if len(flagArgs) > 0 {
		if flagArgs[0] == "stop" {
			cdeamon.Stop()
			return
		} else if flagArgs[0] == "restart" {
			cdeamon.Stop()
		}
	}

	if cdeamon.IsRunning() {
		fmt.Println("already run")
		return
	}
	if cdeamon.IsDeamon() {
		return
	}

	hostname, _ := os.Hostname()
	define.Master = ""
	define.Hostname = hostname
	if *clusterHostsP == "" {
		define.ClusterHosts = append(define.ClusterHosts, hostname)
	} else {
		define.ClusterHosts = strings.Split(*clusterHostsP, ":")
	}

	tasker := ctask.NewTasker()
	tasker.SetTaskConsolePort(8888)
	tasker.AddTask(5, 100000000, "RpcServer", grpc.RpcServer)
	tasker.AddTask(5, 20, "ElectMaster", event.ElectMaster)
	tasker.AddTask(5, 20, "CheckDisk", event.CheckDisk)
	tasker.AddTask(5, 20, "Query", event.Query)
	tasker.RunTask()
}
