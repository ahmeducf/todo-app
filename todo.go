package main

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type Todo struct {
	ID uint			`json:"id"`
	Title string	`json:"title"`
	Completed bool	`json:"completed"`
}

var (
	dbFilePath 	  = path.Join("database", "todo.db")
	listenPort    = ":8080"
)

var db *gorm.DB
var router = mux.NewRouter()

func getAllTodos(w http.ResponseWriter, req *http.Request) {
	var todos []Todo
	if err := db.Find(&todos).Error ;err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	data, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHttpResponse(w, http.StatusOK, data)
}

func addTodoItem(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var todoItem Todo
	if err := json.NewDecoder(req.Body).Decode(&todoItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.First(&Todo{}, todoItem.ID); result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		return 

	} else if result := db.Create(&todoItem); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusConflict)
		return 
	}

	w.WriteHeader(http.StatusCreated)
}

func getTodoItemById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	
	var item Todo
	err := db.First(&item, vars["id"]).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		data, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeHttpResponse(w, http.StatusOK, data)
	}
}

func deleteTodoItemById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var item Todo
	err := db.First(&item, vars["id"]).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		if err := db.Delete(item).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func updateTodoItemById(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	var item Todo
	err := db.First(&item, vars["id"]).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		var newItem Todo
		json.NewDecoder(req.Body).Decode(&newItem)

		item.Title = newItem.Title
		item.Completed = newItem.Completed

		if err := db.Save(item).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		data, err := json.Marshal(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeHttpResponse(w, http.StatusCreated, data)
	}
}

func writeHttpResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func RegisterTodoRoutes() {
	router.HandleFunc("/todos", getAllTodos).Methods("GET")
	router.HandleFunc("/todos", addTodoItem).Methods("POST")
	router.HandleFunc("/todos/{id}", getTodoItemById).Methods("GET")
	router.HandleFunc("/todos/{id}", deleteTodoItemById).Methods("DELETE")
	router.HandleFunc("/todos/{id}", updateTodoItemById).Methods("PATCH")
}

func main() {
	db, _ = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	db.AutoMigrate(&Todo{})


	RegisterTodoRoutes()
	http.ListenAndServe(listenPort, router)
}