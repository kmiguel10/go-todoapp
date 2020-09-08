package main

import (
	"io"
	"net/http"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

var db, _ = gorm.Open("mysql", "root:root@/todolist?charset=utf8&parseTime=True&loc=Local")

type TodoItemModel struct {
	Id          int    `gorm:"primary_key"` //primary key
	Description string //displays in front-end UI
	Completed   bool   //determines if todo item is done or not
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API health is ok")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("description")
	log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database")
	todo := &TodoItemModel{Description: description, Completed: false}
	db.Create(&todo)
	result := db.Last(&todo)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Value)
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	defer db.Close() //close db connection when main() func is returned

	db.Debug().DropTableIfExists(&TodoItemModel{})
	db.Debug().AutoMigrate(&TodoItemModel{})

	log.Info("Starting TodoList API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	http.ListenAndServe(":8000", router)
}
