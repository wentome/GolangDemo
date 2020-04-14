// alert
package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"../define"
)

func PackAlertMessage(shipid string, title string, alertMessage map[string]string) string {
	var gzipBuf bytes.Buffer
	message := define.AlertMessageStruct{}
	message.Shipid = shipid
	message.Title = title
	message.Time = time.Now().Format("2006-01-02 15:04:05")
	message.Message = alertMessage
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

func SendAlert(message string) string {
	resp, err := http.Post(define.AlertUrl, "application/x-www-form-urlencoded", strings.NewReader(message))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Error"
	}
	return string(body)
}
func SendAlertTest() {
	log.Println("SendAlertTest")
}
