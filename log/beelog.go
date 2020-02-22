// beelog
package main

import (
	"time"

	"github.com/astaxie/beego/logs"
)

func init() {
	// info 级别
	logs.SetLogger(logs.AdapterConsole, `{"level":6,"color":true}`)
	//默认7 debug 级别
	logs.SetLogger(logs.AdapterFile, `{"filename":"beelog.log","daily":true,"maxdays":10}`)
	logs.Async()
}

func main() {
	logs.Debug("my book is bought in the year of ", 2016)
	logs.Info("this %s cat is %v years old", "yellow", 3)
	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	logs.Error(1024, "is a very", "good game")
	logs.Critical("oh,crash")
	time.Sleep(time.Second)
}
