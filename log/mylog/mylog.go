package mylog

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"

	"github.com/rifflock/lfshook"
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
	// entry.Message Format("2006-01-02 15:04:05.999999999 -0700 MST")
	b.WriteString(fmt.Sprintf("%s-%s:%s\n",
		entry.Time.Format("2006-01-02 15:04:05.999"),
		strings.ToUpper(entry.Level.String()), entry.Message))
	return b.Bytes(), nil
}

func init() {
	log.SetFormatter(&MyFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	writer, _ := rotatelogs.New(
		"/tmp/golog/test.log.%Y%m%d",
		rotatelogs.WithLinkName("/tmp/golog/test.log"),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(7),
	)
	log.AddHook(lfshook.NewHook(
		lfshook.WriterMap{
			log.DebugLevel: writer,
			log.InfoLevel:  writer,
			log.WarnLevel:  writer,
			log.ErrorLevel: writer,
			log.FatalLevel: writer,
			log.PanicLevel: writer,
		}, &MyFormatter{}))
}

func Newlog() *log.Entry {
	logger := log.WithFields(log.Fields{})
	return logger
}
