// disk
package event

import (
	"strings"
	"time"

	"../../../ctask"
	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/disk"
)

func CheckDisk(ch chan string) {
	count0 := 0
	count1 := 0
	for {
		ctask.TaskContor(ch)

		if count0 > 3 {
			count0 = 0
			ctask.TaskFeedDog(ch)
		}

		if count1 > 600 {
			count1 = 0
			mountPoints := getDiskMountPoint()
			for _, mountPoint := range mountPoints {
				diskUsage, _ := disk.Usage(mountPoint)
				if diskUsage.UsedPercent > 80 {
					logs.Warn("not enough disk space %s:%.1f", mountPoint, diskUsage.UsedPercent)
				}
			}
		}

		time.Sleep(time.Second)
		count0++
		count1++
	}
}

func getDiskMountPoint() []string {
	var mountPoint []string
	partitions, err := disk.Partitions(true)
	if err != nil {
		logs.Error(err)
	} else {
		for _, partition := range partitions {
			if strings.Contains(partition.Device, "/dev/") {
				mountPoint = append(mountPoint, partition.Mountpoint)

			}

		}
		return mountPoint
	}
	return mountPoint
}
