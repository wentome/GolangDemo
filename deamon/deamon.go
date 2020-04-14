package main

import (
	"flag"
	"fmt"
	"time"

	"../../cdeamon"
)

func main() {
	flag.Parse()
	flagArgs := flag.Args()
	if len(flagArgs) > 0 {
		if flagArgs[0] == "stop" {
			cdeamon.Stop()
			return
		} else if flagArgs[0] == "restart" {
			cdeamon.Stop()
		}
	}

	if cdeamon.IsRunning() {
		fmt.Println("already run")
		return
	}
	if cdeamon.IsDeamon() {
		return
	}
	for i := 1; i < 10; i++ {
		time.Sleep(time.Second)
	}
}
