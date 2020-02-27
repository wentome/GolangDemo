// alert
package main

import (
	"./calert"
)

func main() {
	ealert := calert.NewAlert("http://ealert.com/alert", "code")
	ealert.Send("title", `{"a"："1","b":"2"}`)
}
