// wd
package wd

import (
	"log"
	"reflect"
	"time"
)

const (
	Initial = iota
	Started
	Stopped
	Suspend
)

type Task struct {
	taskId       string
	taskFunc     interface{}
	taskArgs     []string
	taskch       chan string
	taskState    int
	taskInitTime int
	taskWarnTime int
	Timer        int
}

type TaskManager struct {
	TaskQueue []*Task
	TaskCount int
	TaskState int
}
type Tasker interface {
	findTaskInQueue(taskId string) int
	AddTask(taskInitTime int, MaxTimeout int, taskId string, taskFunc interface{}, args ...string)
	RemoveTask(taskId string)
	taskStart()
	RunTask()
	// TaskTimer(ch chan string)
	// Dog()
	// GetTask()
	// DecTimer()
	// GetTaskStatus()
	// FeedDog(task string)
	// Run()
}

func NewTasker() Tasker {
	tasker := new(TaskManager)
	tasker.TaskCount = 0
	tasker.TaskState = Initial
	return tasker
}

func (t *TaskManager) findTaskInQueue(taskId string) int {
	for i, task := range t.TaskQueue {
		if task.taskId == taskId {
			return i
		}
	}
	return -1
}

func (t *TaskManager) AddTask(taskWarnTime int, taskInitTime int, taskId string, taskFunc interface{}, args ...string) {
	index := t.findTaskInQueue(taskId)
	if index >= 0 {
		log.Printf("AddTask Failed:%s already in TaskQueue", taskId)
	} else {
		ch := make(chan string)
		var argsList []string
		argsList = append(argsList, args...)
		if len(t.TaskQueue) == t.TaskCount {
			task := Task{taskId, taskFunc, argsList, ch, Initial, taskInitTime, taskWarnTime, taskInitTime}
			t.TaskQueue = append(t.TaskQueue, &task)
			t.TaskCount++
		}
		// log.Println("add t:", *t)

	}

}

func (t *TaskManager) RemoveTask(taskId string) {
	index := t.findTaskInQueue(taskId)
	if index < 0 {
		log.Printf("RemoveTask Failed:%s not in TaskQueue", taskId)
	} else {
		close(t.TaskQueue[index].taskch)
		for i := index; i < t.TaskCount-1; i++ {
			t.TaskQueue[i] = t.TaskQueue[i+1]
		}
		t.TaskCount--
		t.TaskQueue[t.TaskCount] = nil
	}

	// log.Println("remove t:", *t)
}

func (t *TaskManager) taskStart() {
	for _, task := range t.TaskQueue {
		if task.taskState != Started {
			// log.Printf("TaskStart:%s", task.taskId)
			taskFunc := task.taskFunc
			fv := reflect.ValueOf(taskFunc)
			params := make([]reflect.Value, len(task.taskArgs)+1)
			params[0] = reflect.ValueOf(task.taskch)
			for i, arg := range task.taskArgs {
				params[i+1] = reflect.ValueOf(arg)
			}
			go fv.Call(params)
			task.taskState = Started
			t.TaskState = Started
		}
	}
}
func (t *TaskManager) RunTask() {
	t.taskStart()
	for {
		cases := make([]reflect.SelectCase, t.TaskCount)
		for i := 0; i < t.TaskCount; i++ {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.TaskQueue[i].taskch)}
		}
		if t.TaskCount > 0 {
			taskIndex, value, _ := reflect.Select(cases)
			msg := value.String()
			if msg == "goodby" {
				t.RemoveTask(t.TaskQueue[taskIndex].taskId)
				log.Println("goodby", taskIndex)
			} else if msg == "feeddog" {
				log.Println("feeddog", t.TaskQueue[taskIndex].taskId)
				// w.TaskMap[w.TaskIndexToId[int32(taskIndex)]].Timer = w.TaskMap[w.TaskIndexToId[int32(taskIndex)]].MaxTimeout
				// if w.TaskIndexToId[int32(taskIndex)] == "TaskTimer" {
				// 	w.Dog()
				// }
			}

		} else {
			time.Sleep(time.Second)
			log.Println("TaskQueue is empty")
		}

	}
}

// func (t *watch) RemoveTask(taskId string) {
// 	close(w.TaskMap[taskId].ch)
// 	time.Sleep(time.Second)
// 	delete(w.TaskMap, taskId)
// 	w.TaskCount--
// 	var index int32 = 0
// 	for taskId, taskMap := range w.TaskMap {
// 		w.TaskIndexToId[index] = taskId
// 		taskMap.taskIndex = index
// 		index++
// 	}
// 	delete(w.TaskIndexToId, w.TaskCount)
// }

// func (w *watch) GetTask() {
// 	for task, TaskManager := range w.TaskMap {
// 		fmt.Println(task, *TaskManager, w.TaskCount, w.TaskIndexToId)
// 	}
// }
// func (w *watch) DecTimer() {
// 	for _, TaskManager := range w.TaskMap {
// 		TaskManager.Timer--
// 	}
// }

// func (w *watch) GetTaskStatus() {
// 	for task, TaskManager := range w.TaskMap {
// 		if TaskManager.Timer == 0 {
// 			fmt.Println(task, "Excp", TaskManager.Timer)
// 		} else if TaskManager.Timer < TaskManager.WarnTimeout {
// 			fmt.Println(task, "Warn", TaskManager.Timer)
// 		} else {
// 			fmt.Println(task, "Normal", TaskManager.Timer)
// 		}
// 	}

// }
// func (w *watch) FeedDog(task string) {
// 	w.TaskMap[task].Timer = w.TaskMap[task].MaxTimeout
// }

// func (w *watch) TaskTimer(ch chan string) {
// 	for {
// 		time.Sleep(time.Second)
// 		ch <- "feeddog"
// 	}
// }

// func (w *watch) Dog() {
// 	for _, taskMap := range w.TaskMap {
// 		taskMap.Timer--
// 	}

// 	for task, taskMap := range w.TaskMap {
// 		log.Println("Dog", task, taskMap.Timer)
// 	}
// }
