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
	var expected = []Todo{
		{Name: "study for math exam"},
		{Name: "take the trash out"},
		{Name: "watch golang tutorial"},
	}

	t.Run("Get all initial todos", func(t *testing.T) {
		todoGetAll(t, expected)
	})
	t.Run("Post a new todo", func(t *testing.T) {
		todo := Todo{Name: "eat some vegetables"}
		todoPost(t, todo)

		expected = append(expected, todo)
		todoGetAll(t, expected)
	})
}

func todoPost(t *testing.T, todo Todo) {
	data, err := json.Marshal(todo)

	if err != nil {
		t.Fatalf("Unexpected error on Todo json marshalling: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/todo", bytes.NewBuffer(data))
	rec := httptest.NewRecorder()

	handleTodoPost(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Fatalf("Expected status 201 but got %v", res.StatusCode)
	}
}

func todoGetAll(t *testing.T, expected []Todo) {
	req := httptest.NewRequest(http.MethodGet, "/todo/all", nil)
	rec := httptest.NewRecorder()

	handleTodoGetAll(rec, req)
	res := rec.Result()
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

func TestWelcome(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handleWelcome(rec, req)
	res := rec.Result()
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
