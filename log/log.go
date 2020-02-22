package main

import (
	"time"

	"./mylog"
	"./sub"
)

func main() {
	logger := mylog.Newlog()
	for i := 0; i < 100; i++ {
		logger.Info("I'll be logged with common and other field")
		logger.Warn("Me too")
		sub.LogTest()
		time.Sleep(time.Second)
	}

}
