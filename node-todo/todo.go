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

func CreateTodo(text string, priority int, done bool) (Todo, error) {
	if priority == 0 {
		priority = 3
	}

	if err := TodoValidation(text, priority); err != nil {
		return Todo{}, err
	}

	return Todo{
		Id:       hash(text),
		Text:     text,
		Priority: priority,
		Done:     done,
	}, nil
}

func TodoValidation(text string, priority int) error {
	if !OnlyEnglishLetters(text) {
		return errors.New("Task description does include non-English letter")
	}

	if priority < 1 || priority > 5 {
		return errors.New("Task's priority is invalid")
	}

	return nil
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

func (t *Todos) lookForId(id string) (bool, int, Todo) {
	for index, element := range t.TodoArray {
		if element.Id == id {
			return true, index, element
		}
	}
	return false, 0, Todo{}
}

func (t *Todos) addTodo(item Todo) error {
	if err := t.validateNewTodoItem(item); err != nil {
		fmt.Println("Failed to add todo item to list!")
		return err
	}

	t.TodoArray = append(t.TodoArray, item)
	return nil
}

func (t *Todos) validateNewTodoItem(item Todo) error {
	if reflect.DeepEqual(item, Todo{}) {
		return errors.New("Todo is an empty struct.")
	}

	if err := TodoValidation(item.Text, item.Priority); err != nil {
		return err
	}

	if !t.isIdUnique(item.Id) {
		fmt.Println("Id is not unique, let's give it another try until we find a free id.")
		t.reserveNewId(&item.Id)
	}

	return nil
}

func (t *Todos) reserveNewId(id *string) {
	for ok := true; ok; ok = !t.isIdUnique(*id) {
		id_to_int, _ := strconv.ParseUint(*id, 10, 64)
		id_to_int += 1
		*id = strconv.FormatUint(id_to_int, 10)
	}
}

func (t *Todos) isIdUnique(id string) bool {
	if found, _, _ := t.lookForId(id); found {
		return false
	}
	return true
}

// Support functions
func prettyPrint(b []byte) []byte {
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		fmt.Println("preetyPrint failed")
	}
	return out.Bytes()
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
