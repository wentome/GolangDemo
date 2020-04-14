package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	CREATE_DB_MANAGER = `CREATE DATABASE IF NOT EXISTS manager`
	CREATE_TB_USER    = `CREATE TABLE IF NOT EXISTS user(
i INT UNSIGNED AUTO_INCREMENT,
code VARCHAR(100) NOT NULL,
name VARCHAR(100) NOT NULL,
message VARCHAR(4096) NOT NULL,
PRIMARY KEY ( i )
)ENGINE=InnoDB DEFAULT CHARSET=utf8`
)

type User struct {
	I       int    `json:"i" orm:"column(i);pk"`
	Code    string `json:"code" orm:"column(code)"`
	Name    string `json:"name" orm:"column(name)"`
	Message string `json:"message" orm:"column(message)"`
}

type userPasswd struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//Remember bool   `json:"remember"`
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func init() {
	orm.RegisterModel(new(User))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "remote:123456@(192.168.220.254:3306)/sys?charset=utf8")
}

var o orm.Ormer

func main() {
	o = orm.NewOrm()
	//create database
	o.Raw(CREATE_DB_MANAGER).Exec()
	orm.RegisterDataBase("manager", "mysql", "remote:123456@(192.168.220.254:3306)/manager?charset=utf8")
	//switch database
	o.Using("manager")
	fmt.Println(o.Raw(CREATE_TB_USER).Exec())
	// fileBytes, err := ioutil.ReadFile("mosquitto.conf")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// lines := bytes.Split(fileBytes, []byte("\n"))
	// for _, line := range lines {
	// 	if len(line) > 1 {
	// 		if line[0] == '#' && line[1] != ' ' {
	// 			fmt.Printf(string(line) + "\n")
	// 		}
	// 	}
	// }
	router := mux.NewRouter()

	router.HandleFunc("/api/user/{code}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		code := vars["code"]
		user := User{}
		fmt.Println(code)
		err := o.Raw("SELECT * FROM  user WHERE code=?", code).QueryRow(&user)
		if err != nil {
			log.Println(err)
			return
		}
		userByte, err := json.Marshal(user)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(userByte))
	})
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/api/userlist", getUserList)
	router.HandleFunc("/api/save", saveUserInfo)

	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("run..")
	log.Fatal(srv.ListenAndServe())

}

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//获取路径
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	//如果都没匹配上就回到index
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//"/" 服务路径
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	state := false
	session, _ := store.Get(r, "cookie-name")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		return
	}
	fmt.Println(string(body))
	up := userPasswd{}
	err = json.Unmarshal(body, &up)
	if err != nil {
		fmt.Printf("User passwd err, %v\n", err)
		return
	}
	//fmt.Println(string(body), up)
	if up.Username == "admin" {
		if up.Password == "admin" {
			state = true
		}
	}
	if state {
		session.Values["authenticated"] = true
		session.Save(r, w)
		fmt.Fprintln(w, `{"state":"pass"}`)
	} else {
		fmt.Fprintln(w, `{"state":"fail"}`)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	fmt.Fprintln(w, "loginout successed!")
}

func getUserList(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		//http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Fprintln(w, `{"state":"unlogin"}`)
		return
	}
	var userlist []User
	num, err := o.Raw("SELECT i, code, name, message FROM user").QueryRows(&userlist)
	if err == nil {
		fmt.Println("user nums: ", num)
	}
	userlistByte, err := json.Marshal(userlist)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(userlistByte))
	fmt.Fprintln(w, string(userlistByte))
}

func saveUserInfo(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "cookie-name")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		//http.Error(w, "Forbidden", http.StatusForbidden)
		fmt.Fprintln(w, `{"state":"unlogin"}`)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("read body err, %v\n", err)
		fmt.Fprintln(w, `{"state":"fail"}`)
		return
	}
	fmt.Println(string(body))
	user := User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Printf("User passwd err, %v\n", err)
		fmt.Fprintln(w, `{"state":"fail"}`)
		return
	}
	tmpuser := User{Code: user.Code}
	err = o.Read(&tmpuser, "Code")
	if err == orm.ErrNoRows {
		fmt.Println("new user")
		id, err := o.Insert(&user)
		if err == nil {
			fmt.Println(id)
			fmt.Fprintln(w, `{"state":"pass"}`)
		} else {
			fmt.Fprintln(w, `{"state":"fail"}`)
		}
	} else if err == nil {
		fmt.Println("update")
		user.I = tmpuser.I
		if num, err := o.Update(&user); err == nil {
			fmt.Println(num)
			fmt.Fprintln(w, `{"state":"pass"}`)

		} else {
			fmt.Fprintln(w, `{"state":"fail"}`)
		}

	}
}
