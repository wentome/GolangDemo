package msql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "remote:123456@(192.168.220.254:3306)/version")
	if err != nil {
		log.Fatal(err)
	}
}

func QuerySingleRows(fruit interface{}, query string, args ...interface{}) {
	rValue := reflect.ValueOf(fruit).Elem()
	values := make([]interface{}, rValue.NumField())
	for i := 0; i < rValue.NumField(); i++ {
		values[i] = rValue.Field(i).Addr().Interface()
	}
	fmt.Println(values)
	err := db.QueryRow(query, args...).Scan(values...)
	if err != nil {
		panic(err.Error())
	}
	t := reflect.ValueOf(values)
	fmt.Println(t)
}

type oneRow struct {
	Id      int
	Version string
}

func QueryRows(containers ...interface{}) {
	// rows, err := db.Query(`SELECT id, version FROM mi`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// columns, _ := rows.Columns()
	for _, container := range containers {
		fmt.Println("-", container)
		reflect.Indirect
	}

	// rValue := reflect.ValueOf(fruit).Elem()
	// // rValue = reflect.Append(rValue, reflect.ValueOf(e))
	// fmt.Println(rValue)

	// rValue := reflect.ValueOf(fruit).Elem()

	// values := make([]interface{}, len(columns))
	// for i := 0; i < rValue.NumField(); i++ {
	// 	values[i] = rValue.Field(i).Addr().Interface()
	// }
	// for rows.Next() {
	// 	err := rows.Scan(values...)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	// fmt.Println(fruit)
}

func GetValue(fruit interface{}) interface{} {
	rType := reflect.TypeOf(fruit).Elem()
	rValue := reflect.ValueOf(fruit).Elem()
	values := make([]interface{}, rType.NumField())
	for i := 0; i < rType.NumField(); i++ {
		values[i] = rValue.FieldByName(rType.Field(i).Name).Addr().Interface()
	}
	return values
}

func GetType(target interface{}) interface{} {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	fmt.Println(t.Elem(), v.Elem())
	if t.Kind() == reflect.Ptr { // 如果传入的是指针，则需要获取指针指向的类型 兼容指针或者变量
		t = t.Elem()
		v = v.Elem()
	}
	fmt.Println(t.Kind(), v.Kind())

	t0 := reflect.New(t).Elem()
	t0.Field(0).SetInt(1)
	t0.Field(1).SetString("abc")

	fmt.Println(t0)

	ts := reflect.SliceOf(t)
	// ts = append(ts, reflect.ValueOf())
	// ts = append(ts, reflect.ValueOf(200, ""))
	fmt.Println(ts)

	// newStruc := reflect.New(t) // 调用反射创建对象
	// res := make([]interface{}, 3)
	// fmt.Println(newStruc.Elem())
	// return res

	type T struct {
		Age  int
		Name string
	}
	// 初始化测试用例
	ts1 := T{}
	t1 := reflect.TypeOf(ts1)
	v1 := reflect.ValueOf(ts1)
	if t1.Kind() == reflect.Ptr { // 如果传入的是指针，则需要获取指针指向的类型 兼容指针或者变量
		t1 = t.Elem()
		v1 = v.Elem()
	}

	// v1.Field(0).SetInt(123) // 内置常用类型的设值方法，利用Field序号get
	fmt.Println(ts1, t1.Field(0), v1, v1.Field(0))

	return reflect.SliceOf(t)
}
