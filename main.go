package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"github.com/k0kubun/pp"
)

var db *gorm.DB

type User struct {
	gorm.Model

	Id   int64
	Name string
}

type ResponseJson struct {
	Id   int64
	Name string
}

func getUser(c web.C, w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(c.URLParams["userid"], 10, 64)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(userID)

	pp.Println(db)

	user := new(User)
	err = db.First(user, userID).Error
	if err  != nil {
		fmt.Println("user Not found")
	}

	var response ResponseJson
	response.Id = userID
	response.Name = user.Name
	json.NewEncoder(w).Encode(response)
}

func getUsers(c web.C, w http.ResponseWriter, r *http.Request) {
}

func getPay(c web.C, w http.ResponseWriter, r *http.Request) {

}

func Route(m *web.Mux) {
	m.Get("/user/:userid", getUser)
	m.Get("/pay/:payid", getPay)
}

func main() {
	var hoge string
	flag.StringVar(&hoge, "hoge", "hoge", "hoge")
	flag.Parse()

	var err error

	db, err = gorm.Open("sqlite3", "user.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		panic(err)
	}

	user := new(User)
	err = db.First(user, 1).Error
	if err != nil {
		fmt.Println(err)
	}

	pp.Println(db)

	Route(goji.DefaultMux)
	goji.Serve()
}
