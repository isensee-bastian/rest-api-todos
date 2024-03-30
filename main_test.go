package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCrudApi(t *testing.T) {
	server := httptest.NewServer(makeMux())
	defer server.Close()
	t.Logf("Server started at %s", server.URL)

	t.Run("Check service available", func(t *testing.T) {
		welcome(t, server.URL)
	})

	var expected = []Todo{
		{Name: "study for math exam"},
		{Name: "take the trash out"},
		{Name: "watch golang tutorial"},
	}

	t.Run("Get all initial todos", func(t *testing.T) {
		todoGetAll(t, expected, server.URL)
	})
	t.Run("Post a new todo", func(t *testing.T) {
		todo := Todo{Name: "eat some vegetables"}
		todoPost(t, todo, server.URL)

		expected = append(expected, todo)
		todoGetAll(t, expected, server.URL)
	})
}

func todoPost(t *testing.T, todo Todo, baseUrl string) {
	data, err := json.Marshal(todo)

	if err != nil {
		t.Fatalf("Unexpected error on Todo json marshalling: %v", err)
	}

	res, err := http.Post(baseUrl+"/todo", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatalf("Error on http request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Fatalf("Expected status 201 but got %v", res.StatusCode)
	}
}

func todoGetAll(t *testing.T, expected []Todo, baseUrl string) {
	res, err := http.Get(baseUrl + "/todo/all")
	if err != nil {
		t.Fatalf("Error on http request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error on body reading: %v", err)
	}

	var actual []Todo
	err = json.Unmarshal(data, &actual)
	if err != nil {
		t.Fatalf("Unexpected error on body parsing: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected todos %+v but got %+v", expected, actual)
	}
}

func welcome(t *testing.T, baseUrl string) {
	res, err := http.Get(baseUrl + "/")
	if err != nil {
		t.Fatalf("Error on request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if expected := "Welcome to this Todo list application!"; string(data) != expected {
		t.Fatalf("Expected body '%v' but got '%v'", expected, string(data))
	}
}