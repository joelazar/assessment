package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"strconv"
)

type Todo struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Priority int    `json:"priority"`
	Done     bool   `json:"done"`
}

type Todos struct {
	TodoArray []Todo `json:"todos"`
	dbFile    string
}

func CreateTodo(text string) Todo {
	return Todo{
		Id:       hash(text),
		Text:     text,
		Priority: 3,
		Done:     false,
	}
}

func CreateTodoWithPriority(text string, priority int) Todo {
	return Todo{
		Id:       hash(text),
		Text:     text,
		Priority: priority,
		Done:     false,
	}
}

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 10)
}

func (t Todos) writeToFile() {
	todoJson, _ := json.Marshal(t)
	if err := ioutil.WriteFile(t.dbFile, prettyPrint(todoJson), 0644); err == nil {
		fmt.Println("Successfully written data to todo-database.json file")
	}

}

func (t *Todos) readFromFile() {
	jsonFile, err := os.Open(t.dbFile)
	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &t)
}

func (t *Todos) addTodo(item Todo) {
	t.TodoArray = append(t.TodoArray, item)
}

func prettyPrint(b []byte) []byte {
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		fmt.Println("preetyPrint failed")
	}
	return out.Bytes()
}
