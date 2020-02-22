package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type Version struct {
	Id         int
	Version    string
	Changelist string
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "remote:123456@(192.168.220.254:3306)/version?charset=utf8")
}
func main() {
	o := orm.NewOrm()
	o.Using("default")
	var version Version
	err := o.Raw("SELECT id, version, changelist FROM  mi WHERE id = ? limit 3", 98).QueryRow(&version)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(version)

	var versions []Version
	_, err = o.Raw("SELECT id, version, changelist FROM  mi WHERE product='release' order by id desc  limit 3 ").QueryRows(&versions)
	if err != nil {
		fmt.Println(err)
	}
	for i, info := range versions {
		fmt.Println(i, info)
	}

	l := logs.GetLogger()
	l.Println("this is a message of http")
	//an official log.Logger with prefix ORM
	logs.GetLogger("ORM").Println("this is a message of orm")
	logs.Async()

	logs.Debug("my book is bought in the year of ", 2016)
	logs.Info("this %s cat is %v years old", "yellow", 3)
	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	logs.Error(1024, "is a very", "good game")
	logs.Critical("oh,crash")

}
