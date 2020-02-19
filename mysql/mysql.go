package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "remote:123456@(192.168.220.254:3306)/version")
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	{ // Query all users
		type mi struct {
			id         int
			product    string
			version    string
			changelist string
			md5        string
			user       string
			url        string
		}

		rows, err := db.Query(`SELECT id, product, version, changelist,md5,user,url FROM mi`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var mis []mi
		for rows.Next() {
			var m mi

			err := rows.Scan(&m.id, &m.product, &m.version, &m.changelist, &m.md5, &m.user, &m.url)
			if err != nil {
				log.Fatal(err)
			}
			mis = append(mis, m)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%#v", mis)
	}

}
