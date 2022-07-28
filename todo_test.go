package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	testDbPath = path.Join("database", "test_todo.db")

	todoItem = TodoItem{ID: 2, Title: "clean room", Completed: false}

	newItemJson = `{"id": 2, "title": "clean room", "completed": false}`
	invalidJson = `{"id" 2 "title": "clean room", "completed": false}`
	invalidFieldsJson = `{"invalid": 2, "invalid": "clean room", "completed": false}`
	toggoleCompletedJson = `{"id": 2, "title": "clean room", "completed": true}`
	nonExistingItemJson = `{"id": 10000, "title": "clean room", "completed": false}`
)

func createTestingDb(tstmp int) *TodoServer {
	app := New()
	app.dbFilePath = testDbPath + strconv.Itoa(tstmp)
	app.initDb()

	return app
}

func removeDatabase(app *TodoServer) {
	os.RemoveAll(app.dbFilePath)
}

func assertStatusCode(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

func assertEqualTodoItems(t *testing.T, got, want any) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n%v\nwant:\n%v\n", got, want)
	}
}

func TestGetAllTodos(t *testing.T) {
	t.Run("get list of todos items with response status code StatusOK", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		request := httptest.NewRequest(http.MethodGet, "localhost:8080/todos", nil)
		response := httptest.NewRecorder()

		app.db.Create(&TodoItem{ID: 1, Title: "clean the room", Completed: false})
		app.GetAllTodos(response, request)
		var got []TodoItem
		json.NewDecoder(response.Body).Decode(&got)
		want := []TodoItem{{ID: 1, Title: "clean the room", Completed: false}}

		assertStatusCode(t, response.Code, http.StatusOK)
		assertEqualTodoItems(t, got, want)
	})
}



func TestAddTodoItem(t *testing.T) {
	t.Run("add new todo item with code StatusCreated", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todos", bytes.NewBuffer([]byte(newItemJson)))
		response := httptest.NewRecorder()

		app.AddTodoItem(response, request)
		var got TodoItem
		json.NewDecoder(response.Body).Decode(&got)
		want := TodoItem{ID: 2, Title: "clean room", Completed: false}

		assertStatusCode(t, response.Code, http.StatusCreated)
		assertEqualTodoItems(t, got, want)
	})

	t.Run("add existing todo item with StatusConflict", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		app.db.Create(&TodoItem{ID: 2, Title: "clean room", Completed: false})

		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todos", bytes.NewBuffer([]byte(newItemJson)))
		response := httptest.NewRecorder()

		app.AddTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusConflict)
	})

	t.Run("invalid json syntax with StatusBadRequest", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todos", bytes.NewBuffer([]byte(invalidJson)))
		response := httptest.NewRecorder()

		app.AddTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("invalid json fields with StatusBadRequest", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		request := httptest.NewRequest(http.MethodPost, "localhost:8080/todos", bytes.NewBuffer([]byte(invalidFieldsJson)))
		response := httptest.NewRecorder()

		app.AddTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})
}

func TestUpdateTodoItem(t *testing.T) {
	t.Run("update todo item with StatusCreated", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		app.db.Create(&TodoItem{ID: 2, Title: "clean room", Completed: false})

		request := httptest.NewRequest(http.MethodPatch, "localhost:8080/todos", bytes.NewBuffer([]byte(toggoleCompletedJson)))
		response := httptest.NewRecorder()

		app.UpdateTodoItem(response, request)

		var got TodoItem
		app.db.Find(&got, todoItem)
		want := TodoItem{ID: 2, Title: "clean room", Completed: true}

		assertStatusCode(t, response.Code, http.StatusCreated)
		assertEqualTodoItems(t, got, want)
	})

	t.Run("invalid todo item with StatusBadRequest", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		app.db.Create(&TodoItem{ID: 2, Title: "clean room", Completed: false})

		request := httptest.NewRequest(http.MethodPatch, "localhost:8080/todos", bytes.NewBuffer([]byte(invalidJson)))
		response := httptest.NewRecorder()

		app.UpdateTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("invalid todo item fields with StatusBadRequest", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		app.db.Create(&TodoItem{ID: 2, Title: "clean room", Completed: false})

		request := httptest.NewRequest(http.MethodPatch, "localhost:8080/todos", bytes.NewBuffer([]byte(invalidFieldsJson)))
		response := httptest.NewRecorder()

		app.UpdateTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
	})

	t.Run("non existing todo item fields with StatusNotFound", func(t *testing.T) {
		app := createTestingDb(int(time.Now().UnixMilli()))
		defer removeDatabase(app)

		app.db.Create(&TodoItem{ID: 2, Title: "clean room", Completed: false})

		request := httptest.NewRequest(http.MethodPatch, "localhost:8080/todos", bytes.NewBuffer([]byte(nonExistingItemJson)))
		response := httptest.NewRecorder()

		app.UpdateTodoItem(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
	})
}