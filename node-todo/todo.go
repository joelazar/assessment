package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

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

func CreateTodo(text string) (Todo, error) {
	return CreateTodoWithPriority(text, 3)
}

func CreateTodoWithPriority(text string, priority int) (Todo, error) {
	if err := TodoValidation(text, priority); err != nil {
		return Todo{}, err
	}

	return Todo{
		Id:       hash(text),
		Text:     text,
		Priority: priority,
		Done:     false,
	}, nil
}

func TodoValidation(text string, priority int) error {
	if !OnlyEnglishLetters(text) {
		return errors.New("Task description does include non-English letter")
	}

	if priority < 1 || priority > 5 {
		return errors.New("Task's prioirty is invalid")
	}

	return nil
}

func hash(s string) string {
	h := fnv.New64a()
	h.Write([]byte(s))
	return strconv.FormatUint(h.Sum64(), 10)
}

func OnlyEnglishLetters(s string) bool {
	s = strings.Replace(s, " ", "", -1)
	return IsLetter(s)
}

func (t Todos) writeToFile() error {
	todoJson, _ := json.Marshal(t)
	return ioutil.WriteFile(t.dbFile, prettyPrint(todoJson), 0644)
}

func (t *Todos) readFromFile() error {
	jsonFile, err := os.Open(t.dbFile)
	if err != nil {
		return err
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteValue, &t)
}

func (t *Todos) addTodo(item Todo) {
	if !reflect.DeepEqual(item, Todo{}) {
		t.TodoArray = append(t.TodoArray, item)
	}
}

func prettyPrint(b []byte) []byte {
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		fmt.Println("preetyPrint failed")
	}
	return out.Bytes()
}
