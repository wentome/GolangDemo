// alert
package calert

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AMessage struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Time    string `json:"time"`
	Message string `json:"message"`
}
type AlertMannager struct {
	AlertUrl     string
	AlertId      string
	AlertMessage AMessage
}
type Alert interface {
	gzipBase64(message interface{}) string
	unGzipBase64(message string) interface{}
	post(message string) (string, error)
	Send(title string, message string) (string, error)
}

func NewAlert(url string, id string) Alert {
	m := new(AlertMannager)
	m.AlertUrl = url
	m.AlertId = id
	return m
}

// struct -> jsonString -> Gzip -> base64 -> string
func (m *AlertMannager) gzipBase64(message interface{}) string {
	var gzipBuf bytes.Buffer
	messageJsonByte, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}
	gzipWriter := gzip.NewWriter(&gzipBuf)
	defer gzipWriter.Close()
	gzipWriter.Write(messageJsonByte)
	gzipWriter.Close()
	messageGzipBase64 := base64.URLEncoding.EncodeToString(gzipBuf.Bytes())
	return messageGzipBase64
}

//  string  -> unbase64  -> unGzip -> jsonString -> struct
func (m *AlertMannager) unGzipBase64(message string) interface{} {
	var messageStruct interface{}
	gzipByte, err := base64.URLEncoding.DecodeString(message)
	if err != nil {
		log.Println(err)
	}
	reader := bytes.NewReader(gzipByte)
	gzipReader, _ := gzip.NewReader(reader)
	defer gzipReader.Close()
	jsonByte, _ := ioutil.ReadAll(gzipReader)
	json.Unmarshal(jsonByte, &messageStruct)
	return messageStruct
}

func (m *AlertMannager) post(message string) (string, error) {
	resp, err := http.Post(m.AlertUrl, "application/x-www-form-urlencoded", strings.NewReader(message))
	if err != nil {
		return "", errors.New(fmt.Sprintln(err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New(fmt.Sprintln(err))
	}
	return string(body), nil
}

func (m *AlertMannager) Send(title string, message string) (string, error) {
	m.AlertMessage.Id = m.AlertId
	m.AlertMessage.Title = title
	m.AlertMessage.Time = time.Now().Format("2006-01-02 15:04:05")
	m.AlertMessage.Message = message
	messageGzipBase64 := m.gzipBase64(m.AlertMessage)
	// abc := m.unGzipBase64(messageGzipBase64)
	// log.Println("unGzipBase64:", abc)
	return m.post(messageGzipBase64)

}
