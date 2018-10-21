package main

import (
	//	"fmt"
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func CompareFiles(file1, file2 string) bool {
	f1, err := ioutil.ReadFile(file1)

	if err != nil {
		log.Fatal(err)
	}

	f2, err := ioutil.ReadFile(file2)

	if err != nil {
		log.Fatal(err)
	}

	return bytes.Equal(f1, f2)
}

func TestCreateSimpleTodo(t *testing.T) {
	assert := assert.New(t)
	todo1, err := CreateTodo("task")

	assert.Equal(nil, err)
	assert.Equal("task", todo1.Text)
	assert.Equal(3, todo1.Priority)
	assert.False(todo1.Done)
}

func TestCreateSimpleTodoWithPriority(t *testing.T) {
	assert := assert.New(t)
	todo1, err := CreateTodoWithPriority("task", 1)

	assert.Equal(nil, err)
	assert.Equal("task", todo1.Text)
	assert.Equal(1, todo1.Priority)
	assert.False(todo1.Done)
}

func TestInvalidPriorityTodo(t *testing.T) {
	assert := assert.New(t)
	todo1, err := CreateTodoWithPriority("task", 24)

	assert.NotEqual(nil, err)
	assert.Equal(Todo{}, todo1)
}

func TestInvalidDescriptionTodo(t *testing.T) {
	assert := assert.New(t)
	todo1, err := CreateTodoWithPriority("áéő Ups", 2)

	assert.NotEqual(nil, err)
	assert.Equal(Todo{}, todo1)
}

func TestReadFile(t *testing.T) {
	assert := assert.New(t)

	todos := Todos{dbFile: "tst/testreadfile.json"}
	todos.readFromFile()

	assert.Equal("16390851413506644199", todos.TodoArray[0].Id)
	assert.Equal("Priority task", todos.TodoArray[0].Text)
	assert.Equal(5, todos.TodoArray[0].Priority)
	assert.False(todos.TodoArray[0].Done)

	assert.Equal("2065259010723891734", todos.TodoArray[1].Id)
	assert.Equal("Less priority task", todos.TodoArray[1].Text)
	assert.Equal(4, todos.TodoArray[1].Priority)
	assert.False(todos.TodoArray[1].Done)

	assert.Equal("1823758784994838383", todos.TodoArray[2].Id)
	assert.Equal("Simple task, which is done already", todos.TodoArray[2].Text)
	assert.Equal(2, todos.TodoArray[2].Priority)
	assert.True(todos.TodoArray[2].Done)
}

func TestWriteFile(t *testing.T) {
	assert := assert.New(t)
	todos := Todos{dbFile: "tst/testwritefile.json"}
	defer os.Remove("tst/testwritefile.json") // remove the test file after this test finished

	todos.readFromFile()
	assert.Equal(0, len(todos.TodoArray), "This file should not exist, therefore todo array should be empty")

	todo, _ := CreateTodo("Task")
	todos.addTodo(todo)

	todo, _ = CreateTodoWithPriority("Very important task", 5)
	todos.addTodo(todo)

	todos.writeToFile()
	assert.True(CompareFiles("tst/testwritefile.json", "tst/testwritefile_expected.json"))
}
