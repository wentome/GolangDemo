// alert
package main

import (
	"fmt"
	"log"

	"time"

	"./calert"
)

func main() {
	ealert := calert.NewAlert("http://192.168.10.96:83/alert", "CP_TEST")
	for i := 0; i < 100; i++ {
		res, err := ealert.Send("title01", fmt.Sprintf(`{"time":"%d","a"："1","b":"2"}`, i))
		if err != nil {
			log.Println(err)
		} else {
			log.Println(res)
		}
		time.Sleep(time.Millisecond * 3)
	}

}
