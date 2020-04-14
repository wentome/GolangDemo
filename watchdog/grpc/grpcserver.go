// client
package grpc

import (
	"context"
	"encoding/hex"
	"net"
	"strings"

	"github.com/astaxie/beego/logs"

	"encoding/json"

	"../define"
	"../utils"

	pb "./protoc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) ForMaster(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	switch in.Data1 {
	case "get":
		switch in.Data2 {
		case "hosts":
			{
				clusterHosts, _ := json.Marshal(define.ClusterHosts)
				return &pb.Response{Data1: string(clusterHosts)}, nil
			}
		case "code":
			{
				host := define.ClusterHosts[0]
				code := QueryAgent(host, "get", "code", "")
				return &pb.Response{Data1: code}, nil
			}
		case "status":
			{
				clusterStatus := make(map[string]string)
				for _, host := range define.ClusterHosts {
					logs.Info(host)
					status := QueryAgent(host, "get", "status", "")
					clusterStatus[host] = status
				}
				clusterStatusBytes, err := json.Marshal(clusterStatus)
				if err != nil {
					logs.Info(err)
				}
				return &pb.Response{Data1: string(clusterStatusBytes)}, nil
			}
		default:
			return &pb.Response{Data1: "Unknow Comand"}, nil
		}
	case "set":
		switch in.Data2 {
		case "license":
			{
				license := in.Data3
				index := strings.LastIndex(license, "/")
				licenseText := license[0:index]
				code := utils.GetCode()
				licenseText += "/" + code
				licenseSignString := license[index+1 : len(license)]

				licenseSign, _ := hex.DecodeString(licenseSignString)
				if utils.VerifySign("", define.PublicKeyString, licenseSign, []byte(licenseText)) {
					return &pb.Response{Data1: "License successed"}, nil
				} else {
					return &pb.Response{Data1: "License failed"}, nil
				}
			}

		default:
			return &pb.Response{Data1: "Unknow Comand"}, nil
		}
	case "test":
		switch in.Data2 {
		case "alert":
			message := map[string]string{
				"disk": "80%",
				"mem":  "90%",
			}
			alertMessage := utils.PackAlertMessage("uat", "test", message)
			res := utils.SendAlert(alertMessage)
			return &pb.Response{Data1: res}, nil
		}

	case "sync_vote":
		json.Unmarshal([]byte(in.Data2), &define.Vote)
		return &pb.Response{Data1: "ok"}, nil
	}
	return &pb.Response{Data1: "Unknow Comand"}, nil
}

func (s *server) ForAgent(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	switch in.Data1 {
	case "get":
		switch in.Data2 {
		case "status":
			return &pb.Response{Data1: "normal"}, nil
		case "master":
			return &pb.Response{Data1: define.Master}, nil
		case "code":
			code := utils.GetCode()
			return &pb.Response{Data1: code}, nil
		default:
			return &pb.Response{Data1: "Unknow Comand"}, nil
		}
	case "set":
		return &pb.Response{Data1: "i don not know set"}, nil
	case "sync_vote":
		json.Unmarshal([]byte(in.Data2), &define.Vote)
		return &pb.Response{Data1: "ok"}, nil
	}
	return &pb.Response{Data1: "Unknow Comand"}, nil
}

func RpcServer(ch chan string) {
	lis, err := net.Listen("tcp", define.Port)
	if err != nil {
		logs.Error("failed to listen: %v", err)
	}
	s := grpc.NewServer() //起一个服务
	pb.RegisterWatchDogServer(s, &server{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logs.Error("failed to serve: %v", err)
	}
}
