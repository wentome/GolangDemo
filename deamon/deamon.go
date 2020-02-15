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
