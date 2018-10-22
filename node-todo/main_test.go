package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

func readBody(resp *http.Response) string {
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return string(bodyBytes)
}

func initMainTest() {
	go func() {
		dbFile = "tst/tempdb.json"
		removePeriod = 1 * time.Second
		startMain()
	}()
	time.Sleep(1 * time.Second) // let's wait for main to bind for localhost:8000
}

func TestMainBasicFunctions(t *testing.T) {
	assert := assert.New(t)
	initMainTest()
	file, _ := os.Open(dbFile)

	resp, _ := http.Get("http://127.0.0.1:8000/todos")

	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("{\n  \"todos\": []\n}", readBody(resp))

	emptyStat, _ := file.Stat()
	assert.True(CompareFiles("tst/empty-database.json", dbFile), "Database file should empty.")

	postBody := []byte(`{"text":"New Task","priority":5,"done":false}`)
	resp, _ = http.Post("http://127.0.0.1:8000/todos", "application/json", bytes.NewBuffer(postBody))

	assert.Equal(http.StatusOK, resp.StatusCode)
	expected := "{\n  \"id\": \"15924477744071265666\",\n  \"text\": \"New Task\",\n  \"priority\": 5,\n  \"done\": false\n}"
	assert.Equal(expected, readBody(resp))
	afterCreateStat, _ := file.Stat()
	assert.True(emptyStat.Size() < afterCreateStat.Size(), "Database file updated.")

	client := http.Client{}

	putBody := []byte(`{"text":"New tasks priority just got decreased","priority":2,"done":false}`)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8000/todos/15924477744071265666", bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	assert.Equal(http.StatusOK, resp.StatusCode)
	expected = "{\n  \"id\": \"15924477744071265666\",\n  \"text\": \"New tasks priority just got decreased\",\n  \"priority\": 2,\n  \"done\": false\n}"
	assert.Equal(expected, readBody(resp))
	afterModifyStat, _ := file.Stat()
	assert.True(afterCreateStat.Size() < afterModifyStat.Size(), "Database file updated, (new task description is longer now).")

	resp, _ = http.Get("http://127.0.0.1:8000/todos/15924477744071265666")
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal(expected, readBody(resp))
	afterGetStat, _ := file.Stat()
	assert.True(afterGetStat.Size() == afterModifyStat.Size(), "Database file should not be updated after a simple get.")

	req, _ = http.NewRequest(http.MethodDelete, "http://127.0.0.1:8000/todos/15924477744071265666", nil)
	resp, _ = client.Do(req)
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("", readBody(resp))
	assert.True(CompareFiles("tst/empty-database.json", dbFile), "Database file should empty after delete operation.")
}

func TestMainDeleteDone(t *testing.T) {
	assert := assert.New(t)

	postBody := []byte(`{"text":"First Task","priority":5,"done":false}`)
	resp, _ := http.Post("http://127.0.0.1:8000/todos", "application/json", bytes.NewBuffer(postBody))
	assert.Equal(http.StatusOK, resp.StatusCode)

	client := http.Client{}

	putBody := []byte(`{"text":"First Task","priority":5,"done":true}`)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8000/todos/16882564041262248690", bytes.NewBuffer(putBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = client.Do(req)
	assert.Equal(http.StatusOK, resp.StatusCode)

	time.Sleep(2 * time.Second) // let's wait for done task to be deleted

	resp, _ = http.Get("http://127.0.0.1:8000/todos")
	assert.Equal(http.StatusOK, resp.StatusCode)
	assert.Equal("{\n  \"todos\": []\n}", readBody(resp))
}
