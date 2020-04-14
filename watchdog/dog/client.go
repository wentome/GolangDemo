// client
package dog

import (
	"fmt"
	"os"

	"../define"
	"../grpc"
)

func getMaster() string {
	return grpc.QueryAgent("localhost", "get", "master", "")
}

func getHosts() string {
	master := getMaster()
	hosts := grpc.QueryMaster(master, "get", "hosts", "")
	return hosts
}

func Client() {
	var (
		res string
	)

	switch os.Args[1] {
	case "get":
		switch os.Args[2] {
		case "hosts":
			{
				hosts := getHosts()
				fmt.Printf("hosts : " + hosts + "\n")
				return
			}
		case "master":
			{
				master := getMaster()
				fmt.Printf("master : " + master + "\n")
				return
			}
		case "status":
			{
				master := getMaster()
				res = grpc.QueryMaster(master, "get", "status", "")
				fmt.Printf("status : " + res + "\n")
				return
			}
		case "code":
			{
				master := getMaster()
				code := grpc.QueryMaster(master, "get", "code", "")
				fmt.Printf("%s\n", code)
				return
			}
		case "help":
			{
				fmt.Printf(define.HelpInfo + "\n")
				return
			}
		}
	case "set":
		switch os.Args[2] {
		case "status":
			{
				master := getMaster()
				res = grpc.QueryMaster(master, "get", "status", "")
				fmt.Printf("master : " + master + "\n")
				fmt.Printf("status : " + res + "\n")
				return
			}
		case "license":
			{
				//`shipid/c/mem/start/end/totle/code/sign`
				license := os.Args[3]
				master := getMaster()
				res = grpc.QueryMaster(master, "set", "license", license)

			}
		}
	case "test":
		switch os.Args[2] {
		case "alert":
			{
				master := getMaster()
				res = grpc.QueryMaster(master, "test", "alert", "")
				fmt.Printf(res + "\n")
				return
			}
		}

	}
	fmt.Printf(res + "\n")
}
