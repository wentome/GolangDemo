// elect
package event

import (
	"encoding/json"

	"time"

	"../../../ctask"
	"../define"
	"../grpc"
	"../utils"
	"github.com/astaxie/beego/logs"
)

type ClusterStr struct {
	Master       string
	ClusterHosts []string
	HostName     string
}

func ElectMaster_(ch chan string) {
	count0 := 0
	for {
		ctask.TaskContor(ch)
		if count0 > 3 {
			count0 = 0
			ctask.TaskFeedDog(ch)
		}
		time.Sleep(time.Second * 1)
		count0++

	}

}
func ElectMaster(ch chan string) {
	count0 := 0
	var nobodyKnowsMaster bool
	//default vote
	for _, host := range define.ClusterHosts {
		define.Vote = append(define.Vote, host)
	}
	for {
		ctask.TaskContor(ch)
		if count0 > 3 {
			count0 = 0
			ctask.TaskFeedDog(ch)
		}
		if define.Master == "" {
			define.IsMaster = false
			logs.Info("no master")
			//who is master
			nobodyKnowsMaster = true
			hostsWithoutThis := utils.HostsWithoutThis(define.Hostname, define.Vote)
			hostsMaster := grpc.QueryClusterAgent(hostsWithoutThis, "get", "master", "")
			for host, master := range hostsMaster {
				if master != "" {
					define.Master = master
					nobodyKnowsMaster = false
					res := grpc.QueryAgent(define.Master, "get", "vote", "")
					if res != "" {
						json.Unmarshal([]byte(res), &define.Vote)
					}
					logs.Info("%s kown master is %s", host, define.Master)
					break
				}
			}

			// nobodyKnowsMaster
			if nobodyKnowsMaster {
				define.Master = define.Vote[0]
				num := len(define.Vote)
				for i := 0; i < num-1; i++ {
					define.Vote[i] = define.Vote[i+1]
				}
				define.Vote[num-1] = define.Master
				logs.Info("Nobody kown master, Elect %s as Master:", define.Master)
			}
			//sync_vote
			if define.Master == define.Hostname {
				define.IsMaster = true
				for _, host := range define.Vote {
					if host == define.Hostname {
						continue
					}
					voteBytes, err := json.Marshal(define.Vote)
					if err != nil {
						logs.Error(err)
					}
					res := grpc.QueryAgent(host, "sync_vote", string(voteBytes), "")
					if res == "" {
						logs.Info("sync_vote to %s failed", host)
						break
					}
				}
			}
		} else {
			res := grpc.QueryAgent(define.Master, "get", "master", "")
			define.Master = res
		}
		time.Sleep(define.ElectMasterInterval * time.Second)
		count0++
	}
}
