// client
package grpc

import (
	"context"
	"time"

	"github.com/astaxie/beego/logs"

	"../define"
	pb "./protoc"

	"google.golang.org/grpc"
)

func QueryAgent(host string, data1 string, data2 string, data3 string) string {
	conn, err := grpc.Dial(host+define.Port, grpc.WithInsecure())
	if err != nil {
		logs.Info("did not connect: %v", err)
		return ""
	}
	defer conn.Close()
	c := pb.NewWatchDogClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ForAgent(ctx, &pb.Request{Data1: data1, Data2: data2, Data3: data3})
	if err != nil {
		logs.Info("could not greet: %v", err)
		return ""
	}
	return r.Data1
}

func QueryMaster(host string, data1 string, data2 string, data3 string) string {
	conn, err := grpc.Dial(host+define.Port, grpc.WithInsecure())
	if err != nil {
		logs.Info("did not connect: %v", err)
		return ""
	}
	defer conn.Close()
	c := pb.NewWatchDogClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.ForMaster(ctx, &pb.Request{Data1: data1, Data2: data2, Data3: data3})
	if err != nil {
		logs.Info("could not greet: %v", err)
		return ""
	}
	return r.Data1
}

func QueryClusterAgent(hosts []string, data1 string, data2 string, data3 string) map[string]string {
	resCluster := make(map[string]string)
	chs := make(map[string]chan string)
	for _, host := range define.ClusterHosts {
		chs[host] = make(chan string)
		go func() {
			res := QueryAgent(host, data1, data2, data3)
			chs[host] <- res

		}()
	}

	for host, ch := range chs {
		resCluster[host] = <-ch
	}
	return resCluster
}
