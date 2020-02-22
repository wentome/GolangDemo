// watchdog
package main

import (
	// "log"
	"fmt"
	"time"

	"./wd"
)

func main() {
	tasker := wd.NewTasker()
	tasker.AddTask(5, 10, "Task1", Task1, "args")
	for i := 0; i < 10; i++ {
		taskId := fmt.Sprintf("Task2-%d", i)
		tasker.AddTask(5, 10, taskId, Task2)
	}
	tasker.RunTask()
}

func Task1(ch chan string, arg1 string) {
	for i := 0; i < 50; i++ {
		time.Sleep(time.Second * 1)
		// log.Printf("Task1")
		ch <- "feeddog"
	}
	ch <- "goodby"
}

func Task2(ch chan string) {
	for i := 0; i < 100000; i++ {
		time.Sleep(time.Second * 2)
		// log.Printf("Task2")
		ch <- "feeddog"
	}
	ch <- "goodby"
}
