// deamon
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("main run ...", os.Getppid())
	if os.Getppid() != 1 {
		filePath, _ := filepath.Abs(os.Args[0])
		args := append([]string{filePath}, os.Args[1:]...)
		os.StartProcess(filePath, args, &os.ProcAttr{Files: []*os.File{os.Stdin, os.Stdout, os.Stderr}})
		return
	}

	for i := 1; i < 10; i++ {
		fmt.Printf("deamon run ... %d\n", i)
	}
	fmt.Printf("\n")
}

func FindProcessPidByName(processName string) []int {
	var pids []int
	fd, _ := ioutil.ReadDir("/proc")
	for _, fi := range fd {
		fiName := fi.Name()
		pid, err := strconv.Atoi(fiName)
		if err == nil {
			statusFile := path.Join("/proc", fiName, "status")
			f, err := ioutil.ReadFile(statusFile)
			if err != nil {
				continue
			}

			name := string(f[6:bytes.IndexByte(f, '\n')])
			if name == processName {
				pids = append(pids, pid)
			}

		} else {
			continue
		}
	}
	return pids
}
func KillProcess(pid int) {
	proc, _ := os.FindProcess(pid)
	proc.Kill()
}
