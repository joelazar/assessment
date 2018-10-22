package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	initRouter()
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func TestGetAll(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/testwritefile_expected.json"
	todos = CreateTodos()
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	file, _ := ioutil.ReadFile(dbFile)
	expected := string(file)

	assert.Equal(expected, resp.Body.String(), "Body should contain the whole database file.")
}

func TestGetById(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/testwritefile_expected.json"
	todos = CreateTodos()
	req, err := http.NewRequest("GET", "/todos/283230875458061868", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	expected := "{\n  \"id\": \"283230875458061868\",\n  \"text\": \"Task\",\n  \"priority\": 3,\n  \"done\": false\n}"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific todo item.")
}

func TestGetByIdNotFound(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/testwritefile_expected.json"
	todos = CreateTodos()
	req, err := http.NewRequest("GET", "/todos/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusNotFound, resp.Code, "Status code should be 404.")
	expected := "Todo item not found.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
}

func TestModifyById(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	todos.addTodo(Todo{"1234", "Task", 1, false})
	todos.addTodo(Todo{"5678", "Second Task", 4, false})
	payload := []byte(`{"text":"Updated Task","priority":2,"done":true}`)
	req, err := http.NewRequest("PUT", "/todos/1234", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	expected := "{\n  \"id\": \"1234\",\n  \"text\": \"Updated Task\",\n  \"priority\": 2,\n  \"done\": true\n}"
	assert.Equal(expected, resp.Body.String(), "Body should contain the updated todo item.")
}

func TestModifyByIdNotFound(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	todos.addTodo(Todo{"1234", "Task", 1, false})
	todos.addTodo(Todo{"5678", "Second Task", 4, false})
	payload := []byte(`{"text":"Updated Task","priority":2,"done":true}`)
	req, err := http.NewRequest("PUT", "/todos/1235", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusBadRequest, resp.Code, "Status code should be 400.")
	expected := "Todo item not found.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
}

func TestCreateNewTodo(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	payload := []byte(`{"text":"New Task","priority":5,"done":false}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	expected := "{\n  \"id\": \"15924477744071265666\",\n  \"text\": \"New Task\",\n  \"priority\": 5,\n  \"done\": false\n}"
	assert.Equal(expected, resp.Body.String(), "Body should contain the new todo item.")
}

func TestCreateNewTodoWithId(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	payload := []byte(`{"id":"12341234","text":"New Task","priority":5,"done":false}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	expected := "{\n  \"id\": \"15924477744071265666\",\n  \"text\": \"New Task\",\n  \"priority\": 5,\n  \"done\": false\n}"
	assert.Equal(expected, resp.Body.String(), "Body should contain the new todo item with a new reserved id.")
}

func TestCreateNewTodoInvalidDescription(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	payload := []byte(`{"id":"12341234","text":"123 Task","priority":3,"done":false}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusBadRequest, resp.Code, "Status code should be 400.")
	expected := "Failed to add todo item, error : Task description does include non-English letter.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
}

func TestCreateNewTodoInvalidPriority(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	payload := []byte(`{"id":"12341234","text":"Task","priority":99,"done":false}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusBadRequest, resp.Code, "Status code should be 400.")
	expected := "Failed to add todo item, error : Task's priority is invalid.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
}

func TestCreateNewTodoInvalidDone(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	payload := []byte(`{"id":"12341234","text":"Task","priority":3,"done":"not bool"}`)
	req, err := http.NewRequest("POST", "/todos", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusBadRequest, resp.Code, "Status code should be 400.")
	expected := "Failed to decode request body.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
}

func TestDeleteTodo(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	todos.addTodo(Todo{"1234", "Task", 1, false})
	todos.addTodo(Todo{"5678", "Second Task", 4, false})
	req, err := http.NewRequest("DELETE", "/todos/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusOK, resp.Code, "Status code should be 200.")
	assert.Equal("", resp.Body.String(), "Body should be empty.")
	assert.Equal(1, len(todos.TodoArray), "Array should contain one item.")
	assert.Equal("5678", todos.TodoArray[0].Id)
}

func TestDeleteTodoNotFound(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/tempdb.json"
	defer os.Remove(dbFile) // remove the test file after this test finished
	todos = CreateTodos()
	todos.addTodo(Todo{"1234", "Task", 1, false})
	todos.addTodo(Todo{"5678", "Second Task", 4, false})
	req, err := http.NewRequest("DELETE", "/todos/1235", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp := executeRequest(req)

	assert.Equal(http.StatusNotFound, resp.Code, "Status code should be 400.")
	expected := "Todo item not found.\n"
	assert.Equal(expected, resp.Body.String(), "Body should contain the specific error message.")
	assert.Equal(2, len(todos.TodoArray), "Array still should contain two items.")
}
