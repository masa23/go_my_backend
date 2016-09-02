package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/k0kubun/pp"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var db *gorm.DB

type User struct {
	Name  string `sql:"not null; unique"`
	Email string
}

func addUser(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := new(User)
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tx := db.Begin()
	err = tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	tx.Commit()

	fmt.Fprintf(w, "%s", "{\"status\": true}")
}

func getUser(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := strconv.ParseInt(c.URLParams["userid"], 10, 64)
	pp.Println(db)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println(userID)

	user := new(User)
	err = db.Find(user, userID).Error
	if err != nil {
		fmt.Println("user Not found")
	}

	json.NewEncoder(w).Encode(user)
}

func getUsers(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	err := db.Find(&users).Error
	if err != nil {
		fmt.Println("user Not found")
	}

	for i := 0; i < len(users); i++ {
		if users[i].Email == "" {
			users[i].Email = "null"
		}
	}

	json.NewEncoder(w).Encode(users)
}

func Route(m *web.Mux) {
	m.Post("/user", addUser)
	m.Get("/user/:userid", getUser)
	m.Get("/users", getUsers)
}

func main() {
	var err error

	db, err = gorm.Open("sqlite3", "api.db")
	if err != nil {
		panic(err)
	}

	db.LogMode(true)
	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		panic(err)
	}

	Route(goji.DefaultMux)
	goji.Serve()
}
