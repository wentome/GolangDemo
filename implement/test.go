// test
package main

import (
	"fmt"
)

// type method interface {
// 	method1()
// 	method2()
// }
type ts struct{}

func (t *ts) method1() {
	fmt.Println("method1")
}
func (t *ts) method2() {
	fmt.Println("method2")
}
func (t *ts) register() {
	fmt.Println("method2")
}

type ts2 struct {
	ts //继承了ts的方法
}

func (t *ts2) method2() {
	fmt.Println("method2")
}

func main() {
	t := ts{}
	t.method1()
	t.method2()
}
