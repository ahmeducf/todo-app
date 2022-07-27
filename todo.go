package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type Todo struct {
	ID uint			`json:"id"`
	Title string	`json:"title"`
	Completed bool	`json:"completed"`
}

var db *gorm.DB
var db_FILE = os.Getenv("DB_FILE")
var router = mux.NewRouter()

func getAllTodos(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var todos []Todo
	db.Find(&todos)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func addTodoItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var todoItem Todo
	if err := json.NewDecoder(req.Body).Decode(&todoItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	if db.First(&Todo{}, todoItem.ID).Error == nil {
		w.WriteHeader(http.StatusConflict)
	} else {
		db.Create(&todoItem)
		w.WriteHeader(http.StatusOK)
	}
}

func getTodoItemById(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	
}

func RegisterTodoRoutes() {
	router.HandleFunc("/todos", getAllTodos).Methods("GET")
	router.HandleFunc("/todos", addTodoItem).Methods("POST")
	router.HandleFunc("/todos/{id}", getTodoItemById)
}

func main() {
	db, _ = gorm.Open(sqlite.Open(db_FILE), &gorm.Config{})
	db.AutoMigrate(&Todo{})


	RegisterTodoRoutes()
	http.ListenAndServe(":8080", router)
}