// cacser
package main

import (
	"bytes"
	//"errors"
	"encoding/json"
	"fmt"
	"log"

	//"os"
	"time"

	"../../acser"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	CREATE_DB_ALERT          = `CREATE DATABASE IF NOT EXISTS alert`
	CREATE_TB_ALERT_MESSAGES = `CREATE TABLE IF NOT EXISTS alert_messages(
i INT UNSIGNED AUTO_INCREMENT,
id VARCHAR(100) NOT NULL,
title VARCHAR(100) NOT NULL,
time DATETIME,
message VARCHAR(4096) NOT NULL,
PRIMARY KEY ( i )
)ENGINE=InnoDB DEFAULT CHARSET=utf8`
)

type AMessage struct {
	I       int    `json:"i" orm:"column(i);pk"`
	Id      string `json:"id" orm:"column(id)"`
	Title   string `json:"title" orm:"column(title)"`
	Time    string `json:"time" orm:"column(time)"`
	Message string `json:"message" orm:"column(message)"`
}

func (u *AMessage) TableName() string {
	return "alert_messages"
}

var messageStruct AMessage

func myParse(message []byte) {
	messageSeg := bytes.Split(message, []byte("|"))
	if len(messageSeg) == 4 {
		seg := bytes.Split(message, []byte("|"))[2]
		jsonByte, err := acser.UnGzipBase64(string(seg))
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(jsonByte, &messageStruct)
		if err != nil {
			log.Println(err)
			return
		}
		name, node, offset := ac.GetProgress()
		log.Println(name, node, offset, messageStruct)
		o.Insert(&messageStruct)
		time.Sleep(time.Millisecond * 0)
	} else {
		log.Println("Parse file failed!")
	}
}

func init() {
	orm.RegisterModel(new(AMessage))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "remote:123456@(192.168.10.96:3306)/sys?charset=utf8")
}

var o orm.Ormer
var ac acser.Acser

func main() {
	o = orm.NewOrm()
	//create database
	o.Raw(CREATE_DB_ALERT).Exec()
	orm.RegisterDataBase("alert", "mysql", "remote:123456@(192.168.10.96:3306)/alert?charset=utf8")
	//switch database
	o.Using("alert")
	fmt.Println(o.Raw(CREATE_TB_ALERT_MESSAGES).Exec())
	//os.Exit(0)
	log.SetFlags(log.Lshortfile)
	ac = acser.NewAcser()
	ac.SetAcserFile("/root/nginx/logs", "access.log.*", "access.log")
	ac.RegisterParseFunc(myParse)
	err := ac.Run()
	fmt.Println(err)
}
