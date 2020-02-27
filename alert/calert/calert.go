// alert
package calert

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"

	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type AlertMessage struct {
	Id      string
	Title   string
	Time    string
	Message string
}
type AlertMannager struct {
	Url string
	Id  string
	AM  AlertMessage
}
type Alert interface {
	gzipBase64(message interface{}) string
	unGzipBase64(message string) interface{}
	post(message string) string
	Send(title string, message string)
}

func NewAlert(url string, id string) Alert {
	am := new(AlertMannager)
	am.Url = url
	am.Id = id
	return am
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

func (m *AlertMannager) post(message string) string {
	resp, err := http.Post(m.Url, "application/x-www-form-urlencoded", strings.NewReader(message))
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error"
	}
	return string(body)
}

func (m *AlertMannager) Send(title string, message string) {
	m.AM.Id = m.Id
	m.AM.Title = title
	m.AM.Time = time.Now().Format("2006-01-02 15:04:05")
	m.AM.Message = message
	messageGzipBase64 := m.gzipBase64(m.AM)
	abc := m.unGzipBase64(messageGzipBase64)
	log.Println("unGzipBase64:", abc)
}
