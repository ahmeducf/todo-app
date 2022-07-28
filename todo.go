package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)


type TodoItem struct {
	ID uint			`json:"id"`
	Title string	`json:"title"`
	Completed bool	`json:"completed"`
}

const (
	listenPort    = ":8080"
	badRequest 	  = "Request is not valid\n"
	JsonIndent = "    "
)

type TodoServer struct {
	todoItem TodoItem
	db *gorm.DB
	dbFilePath string
}

func New() *TodoServer {
	return &TodoServer{todoItem: TodoItem{}, db: nil}
}

func (app *TodoServer) initDb() error {
	db, err := gorm.Open(sqlite.Open(app.dbFilePath), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&TodoItem{})
	app.db = db

	return err
}

func (app *TodoServer) GetAllTodos(w http.ResponseWriter, req *http.Request) {
	var todos []TodoItem
	if err := app.db.Find(&todos).Error ;err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	data, err := json.MarshalIndent(todos, "", JsonIndent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeHttpResponse(w, http.StatusOK, data)
}

func isBadRequest(item TodoItem) bool{
	if item.ID == 0 || len(item.Title) == 0 {
		return true
	}

	return false
}
func (app *TodoServer) AddTodoItem(w http.ResponseWriter, req *http.Request) {
	var todoItem TodoItem
	if err := json.NewDecoder(req.Body).Decode(&todoItem); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if isBadRequest(todoItem) {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	if result := app.db.First(&TodoItem{}, todoItem.ID); result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		return 

	} else if result := app.db.Create(&todoItem); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return 
	}

	data, err := json.MarshalIndent(todoItem, "", JsonIndent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	writeHttpResponse(w, http.StatusCreated, data)
}

func (app *TodoServer) UpdateTodoItem(w http.ResponseWriter, req *http.Request) {
	var item TodoItem
	var newItem TodoItem
	err := json.NewDecoder(req.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if newItem.ID == 0 || len(newItem.Title) == 0 {
		http.Error(w, badRequest, http.StatusBadRequest)
		return
	}

	err = app.db.First(&item, newItem.ID).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		item.Title = newItem.Title
		item.Completed = newItem.Completed
		if err := app.db.Save(item).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		data, err := json.MarshalIndent(item, "", JsonIndent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeHttpResponse(w, http.StatusCreated, data)
	}
}

func (app *TodoServer) GetTodoItemById(w http.ResponseWriter, req *http.Request) {
	var id int
	var err error
	vars := mux.Vars(req)

	if id, err = strconv.Atoi(vars["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item TodoItem
	err = app.db.First(&item, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		data, err := json.MarshalIndent(item, "", JsonIndent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		writeHttpResponse(w, http.StatusOK, data)
	}
}

func (app *TodoServer) DeleteTodoItemById(w http.ResponseWriter, req *http.Request) {
	var id int
	var err error
	vars := mux.Vars(req)
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var item TodoItem
	err = app.db.First(&item, id).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else {
		if err := app.db.Delete(item).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func writeHttpResponse(w http.ResponseWriter, statusCode int, data []byte) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func (app *TodoServer)registerTodoRoutes(router *mux.Router) {
	router.HandleFunc("/todos", app.GetAllTodos).Methods("GET")
	router.HandleFunc("/todos", app.AddTodoItem).Methods("POST")
	router.HandleFunc("/todos", app.UpdateTodoItem).Methods("PATCH")
	router.HandleFunc("/todos/{id}", app.GetTodoItemById).Methods("GET")
	router.HandleFunc("/todos/{id}", app.DeleteTodoItemById).Methods("DELETE")
}

func main() {
	app := New()
	app.dbFilePath = path.Join("database", "todo.db")
	err := app.initDb()
	if err != nil {
		panic(err)
	}

	var router = mux.NewRouter()
	app.registerTodoRoutes(router)
	http.ListenAndServe(listenPort, router)
}