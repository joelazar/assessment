package main

import (
	//	"fmt"
	"github.com/stretchr/testify/assert"
	//	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	assert := assert.New(t)

	todos := Todos{dbFile: "tst/testreadfile.json"}
	todos.readFromFile()

	assert.Equal("16390851413506644199", todos.TodoArray[0].Id)
	assert.Equal("Priority task", todos.TodoArray[0].Text)
	assert.Equal(5, todos.TodoArray[0].Priority)
	assert.Equal(false, todos.TodoArray[0].Done)

	assert.Equal("2065259010723891734", todos.TodoArray[1].Id)
	assert.Equal("Less priority task", todos.TodoArray[1].Text)
	assert.Equal(4, todos.TodoArray[1].Priority)
	assert.Equal(false, todos.TodoArray[1].Done)

	assert.Equal("1823758784994838383", todos.TodoArray[2].Id)
	assert.Equal("Simple task, which is done already", todos.TodoArray[2].Text)
	assert.Equal(2, todos.TodoArray[2].Priority)
	assert.Equal(true, todos.TodoArray[2].Done)
}

//defer os.Remove("test.json")
