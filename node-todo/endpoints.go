package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

var (
	router *mux.Router
	mutex  *sync.Mutex
)

func initRouter() {
	mutex = &sync.Mutex{}
	router = mux.NewRouter()
	router.HandleFunc("/todos", HttpGetAll).Methods("GET")
	router.HandleFunc("/todos", HttpCreateTodo).Methods("POST")
	router.HandleFunc("/todos/{id}", HttpGetTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", HttpModifyTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", HttpDeleteTodo).Methods("DELETE")
}

func HttpGetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /todos received - query all todo items")
	mutex.Lock()
	defer mutex.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write(todos.getAll())
}

func HttpCreateTodo(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /todos received - try to create new todo item")
	mutex.Lock()
	defer mutex.Unlock()
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode request body.", http.StatusBadRequest)
		return
	}
	response, err := todos.CreateTodoItem(todo)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to add todo item, error : %v", err), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func HttpGetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Printf("GET /todos/%s received - look for given id in database", params["id"])
	mutex.Lock()
	defer mutex.Unlock()
	response, err := todos.getById(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func HttpModifyTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Printf("PUT /todos/%s received - try to modify specific todo item", params["id"])
	mutex.Lock()
	defer mutex.Unlock()
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Failed to decode request body.", http.StatusBadRequest)
		return
	}

	response, err := todos.ModifyTodo(params["id"], todo)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func HttpDeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Printf("DELETE /todos/%s received - delete specific todo from database", params["id"])
	mutex.Lock()
	defer mutex.Unlock()
	err := todos.deleteById(params["id"])
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
