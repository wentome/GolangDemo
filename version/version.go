// version
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Mi struct {
	Id         int    `json:"id" orm:"column(id)"`
	Product    string `json:"product" orm:"column(product)"`
	Version    string `json:"version" orm:"column(version)"`
	Changelist string `json:"changelist" orm:"column(changelist)"`
	Url        string `json:"url" orm:"column(url)"`
	Md5        string `json:"md5" orm:"column(md5)"`
}

func init() {
	orm.RegisterModel(new(Mi))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "remote:123456@(192.168.220.254:3306)/version?charset=utf8")
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

func main() {
	o := orm.NewOrm()
	router := mux.NewRouter()

	router.HandleFunc("/version/{ver}/{pro}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ver := vars["ver"]
		pro := vars["pro"]
		mis := make([]Mi, 100)
		_, err := o.Raw(fmt.Sprintf("SELECT version, changelist,url,md5 FROM  mi WHERE product='%s' and version LIKE '%s%%' order by id desc limit 100", pro, ver)).QueryRows(&mis)
		if err != nil {
			log.Println(err)
		}
		res, err := json.Marshal(mis)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, string(res))
	})

	spa := spaHandler{staticPath: "build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    ":80",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("run..")
	log.Fatal(srv.ListenAndServe())
}
