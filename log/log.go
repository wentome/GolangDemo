package main

import (
	"bytes"
	// "encoding/json"
	"fmt"
	"os"

	"strings"

	log "github.com/sirupsen/logrus"
)

type MyFormatter struct{}

func (mf *MyFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	// entry.Message 就是需要打印的日志Format("2006-01-02 15:04:05.999999999 -0700 MST")
	b.WriteString(fmt.Sprintf("%s-%s:%s\n", entry.Time.Format("2006-01-02 15:04:05.999"),
		strings.ToUpper(entry.Level.String()), entry.Message))
	return b.Bytes(), nil
}
func init() {
	log.SetFormatter(&MyFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func main() {

	contextLogger := log.WithFields(log.Fields{})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Warn("Me too")

}
