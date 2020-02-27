// cacser
package acser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"time"
)

type ParseFunc func(message []byte) ([]byte, error)
type AcserManager struct {
	acserPath             string
	acserFileName         string
	acserFileOffset       int64
	acserRecordFileName   string
	acserRecordFileOffset int64
	acserParseFunc        ParseFunc
}

type Acser interface {
	SetAcserPath(acserPath string)
	RegisterParseFunc(parseFunc ParseFunc)
	Run() error
	getAcserFile() (string, error)
}

func NewAcser() Acser {
	acser := new(AcserManager)
	return acser
}

func (a *AcserManager) SetAcserPath(acserPath string) {
	a.acserPath = acserPath
}

func (a *AcserManager) RegisterParseFunc(parseFunc ParseFunc) {
	a.acserParseFunc = parseFunc
}

//如果路径为空返回错误
//如果a.acserFileName为"" 返回最后一个
//如果a.acserFileName 没在路径里 则返回错误
func (a *AcserManager) getAcserFile() (string, error) {
	files, err := ioutil.ReadDir(a.acserPath)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Open %s failed!", a.acserPath))
	}
	filesNum := len(files)
	if filesNum == 0 {
		return "", errors.New(fmt.Sprintf("no file in %s!", a.acserPath))
	}
	if a.acserFileName == "" {
		return files[filesNum-1].Name(), nil
	}
	for i := 0; i < filesNum; i++ {
		if files[i].Name() == a.acserFileName {
			if i > 0 {
				return files[i-1].Name(), nil
			} else {
				return a.acserFileName, nil
			}
		}
	}
	return "", errors.New(fmt.Sprintf("%s not in %s!", a.acserFileName, a.acserPath))
}

func (a *AcserManager) Run() error {
	var acserFile *os.File
	var acserBufReader *bufio.Reader
	var acserFileOpened bool
	var acserCount int
	var acserCount2 int
	//获取进度记录
	a.acserRecordFileName = ""
	a.acserRecordFileOffset = 0
	a.acserFileName = ""
	a.acserFileOffset = 0
	acserFileOpened = false
	for {
		if acserFileOpened == false {
			file, err := a.getAcserFile()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 3)
				continue
			}
			log.Println("Open:", file)
			a.acserFileName = file
			acserFile, err = os.Open(path.Join(a.acserPath, a.acserFileName))
			defer acserFile.Close()
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * 3)
				continue
			}
			// 从文件开头移动 read offset 到进度文件记录的 record_offset
			a.acserFileOffset, _ = acserFile.Seek(a.acserRecordFileOffset, os.SEEK_SET)
			acserBufReader = bufio.NewReaderSize(acserFile, 4096*4)
			acserFileOpened = true

		}
		// 循环读取 解析文件
		line, _, _ := acserBufReader.ReadLine()
		// 如果读取到新行
		if len(line) != 0 {
			a.acserFileOffset, _ = acserFile.Seek(0, os.SEEK_CUR)
			res, err := a.acserParseFunc(line)
			if err != nil {
				fmt.Println(err)
			} else {
				messageMap, err := a.unGzipBase64(string(res))
				if err != nil {
					log.Println(err)
				} else {
					acserCount++
					if acserCount-acserCount2 > 100000 {
						fmt.Println(time.Now().Unix(), acserCount, messageMap)
						acserCount2 = acserCount

					}

				}
			}
		} else {
			//需要加入 判断读完了
			info, _ := acserFile.Stat()
			size := info.Size()
			if size == a.acserFileOffset {
				file, err := a.getAcserFile()
				if err != nil {
					log.Println(err)
				} else {
					if file != a.acserFileName {
						acserFile.Close()
						acserFileOpened = false
						continue
					}
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}

}
func (a *AcserManager) unGzipBase64(message string) (interface{}, error) {
	var err error
	var messageStruct interface{}
	gzipByte, err := base64.URLEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(gzipByte)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()
	jsonByte, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonByte, &messageStruct)
	return messageStruct, err
}

func callFunc(function interface{}, args ...interface{}) {
	fv := reflect.ValueOf(function)
	params := make([]reflect.Value, len(args))
	for i, arg := range args {
		params[i] = reflect.ValueOf(arg)
	}
	fv.Call(params)
}
