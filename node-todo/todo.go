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
	"time"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

var cancelChannel chan string

type Todo struct {
	Id       string `json:"id"`
	Text     string `json:"text"`
	Priority int    `json:"priority"`
	Done     bool   `json:"done"`
}

type Todos struct {
	TodoArray []Todo `json:"todos"`
}

func CreateTodos() Todos {
	cancelChannel = make(chan string, 1)
	t := Todos{}
	t.readFromFile()
	return t
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
		return errors.New("Task description does include non-English letter.")
	}

	if priority < 1 || priority > 5 {
		return errors.New("Task's priority is invalid.")
	}

	return nil
}

func (t Todos) writeToFile() error {
	todoJson, _ := json.Marshal(t)
	return ioutil.WriteFile(dbFile, prettyPrint(todoJson), 0644)
}

func (t *Todos) readFromFile() error {
	jsonFile, err := os.Open(dbFile)
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
	if err := t.validateTodoItem(&item); err != nil {
		return err
	}

	t.TodoArray = append(t.TodoArray, item)
	return nil
}

func (t *Todos) validateTodoItem(item *Todo) error {
	if reflect.DeepEqual(item, Todo{}) {
		return errors.New("Todo is an empty struct.")
	}

	if err := TodoValidation(item.Text, item.Priority); err != nil {
		return err
	}

	if !t.isIdUnique(item.Id) {
		t.reserveNewId(&(item.Id)) // let's give it another try until we find a free id
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

func (t *Todos) getAll() []byte {
	todoJson, _ := json.Marshal(t)
	return prettyPrint(todoJson)
}

func (t *Todos) getById(id string) ([]byte, error) {
	found, _, todo := t.lookForId(id)
	if found {
		todoJson, _ := json.Marshal(todo)
		return prettyPrint(todoJson), nil
	} else {
		return []byte{}, errors.New("Todo item not found.")
	}
}

func (t *Todos) CreateTodoItem(new_todo Todo) ([]byte, error) {
	todo, err := CreateTodo(new_todo.Text, new_todo.Priority, new_todo.Done)

	if err != nil {
		return []byte{}, err
	}

	err = t.addTodo(todo)

	if err != nil {
		return []byte{}, err
	}

	todoJson, _ := json.Marshal(t.TodoArray[len(t.TodoArray)-1]) // the id created by CreateTodo can be different from the id which is in the database
	t.writeToFile()
	return prettyPrint(todoJson), nil
}

func (t *Todos) ModifyTodo(id string, modified_todo Todo) ([]byte, error) {
	found, index, todo := t.lookForId(id)
	if found {
		if err := TodoValidation(modified_todo.Text, modified_todo.Priority); err != nil {
			return []byte{}, err
		}

		if modified_todo.Done && !todo.Done {
			go func() {
				for {
					select {
					case id := <-cancelChannel:
						if todo.Id == id {
							fmt.Printf("Cancel thread for %s.", todo.Id)
							break
						}
					case <-time.After(5 * time.Minute):
						doneChannel <- todo.Id
						break
					}
				}
			}()
		} else if !modified_todo.Done && todo.Done {
			cancelChannel <- todo.Id
		}

		todo.Done = modified_todo.Done
		todo.Text = modified_todo.Text
		todo.Priority = modified_todo.Priority

		t.TodoArray[index] = todo
		todoJson, _ := json.Marshal(todo)

		t.writeToFile()
		return prettyPrint(todoJson), nil
	} else {
		return []byte{}, errors.New("Todo item not found.")
	}
}

func (t *Todos) deleteById(id string) error {
	found, index, _ := t.lookForId(id)
	if found {
		t.TodoArray = append(t.TodoArray[:index], t.TodoArray[index+1:]...)
		t.writeToFile()
		return nil
	} else {
		return errors.New("Todo item not found.")
	}
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
