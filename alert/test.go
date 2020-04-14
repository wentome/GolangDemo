// alert
package main

import (
	"fmt"
	"log"

	"time"

	"./calert"
)

func main() {
	ealert := calert.NewAlert("http://localhost/alert", "code")
	for i := 0; i < 1000000; i++ {
		res, err := ealert.Send("title", fmt.Sprintf(`{"time":"%d","a"："1","b":"2"}`, i))
		if err != nil {
			log.Println(err)
		} else {
			log.Println(res)
		}
		time.Sleep(time.Millisecond * 0)
	}

}
