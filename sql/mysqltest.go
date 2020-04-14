package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"

	//"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

const (
	CREATE_DB_ALERT      = `CREATE DATABASE IF NOT EXISTS alert`
	CREATE_TB_ALERT_USER = `CREATE TABLE IF NOT EXISTS user(
id INT UNSIGNED AUTO_INCREMENT,
name VARCHAR(100) NOT NULL,
date DATETIME,
PRIMARY KEY ( id )
)ENGINE=InnoDB DEFAULT CHARSET=utf8`
)

type Version struct {
	Id         int
	Version    string
	Changelist string
}

type User struct {
	Id   int    `orm:"column(id)"`
	Name string `orm:"column(name)"`
	Date string `orm:"column(date)"`
}

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "remote:123456@(192.168.220.254:3306)/sys?charset=utf8")
}

func main() {
	o := orm.NewOrm()
	//create database
	o.Raw(CREATE_DB_ALERT).Exec()
	orm.RegisterDataBase("alert", "mysql", "remote:123456@(192.168.220.254:3306)/alert?charset=utf8")
	//switch database
	o.Using("alert")
	fmt.Println(o.Raw(CREATE_TB_ALERT_USER).Exec())

	//单条插入
	start := time.Now()
	user := new(User)
	for i := 0; i < 100; i++ {
		user.Name = "slene"
		user.Date = time.Now().Format("2006-01-02 15:04:05")
		o.Insert(user)
	}
	cost := time.Since(start)
	fmt.Printf("cost=[%s]", cost)

	//并发批量插入
	start = time.Now()
	users := make([]User, 100)
	for i := 0; i < 100; i++ {
		users[i].Name = "slene"
		users[i].Date = time.Now().Format("2006-01-02 15:04:05")
	}
	//
	num, _ := o.InsertMulti(20, users)
	cost = time.Since(start)
	fmt.Printf("%d cost=[%s]", num, cost)

	// var version Version
	// err := o.Raw("SELECT id, version, changelist FROM  mi WHERE id = ? limit 3", 98).QueryRow(&version)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(version)

	// var versions []Version
	// _, err = o.Raw("SELECT id, version, changelist FROM  mi WHERE product='release' order by id desc  limit 3 ").QueryRows(&versions)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for i, info := range versions {
	// 	fmt.Println(i, info)
	// }

}
