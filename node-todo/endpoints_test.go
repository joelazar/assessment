package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAll(t *testing.T) {
	assert := assert.New(t)
	dbFile = "tst/testwritefile_expected.json"
	todos = CreateTodos()
	req, err := http.NewRequest("GET", "/todos", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HttpGetAll)

	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code, "Status code should be 200.")
	file, _ := ioutil.ReadFile(dbFile)
	expected := string(file)

	assert.Equal(rr.Body.String(), expected, "Body should contain the whole database file.")
}
