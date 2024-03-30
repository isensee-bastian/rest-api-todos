package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestTodoApi(t *testing.T) {
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
		getAll(t, expected, server.URL)
	})

	t.Run("Get a specific todo", func(t *testing.T) {
		lastIndex := len(expected) - 1
		get(t, lastIndex, expected[lastIndex], server.URL)
	})

	t.Run("Post a new todo", func(t *testing.T) {
		todo := Todo{Name: "eat some vegetables"}
		post(t, todo, server.URL)

		expected = append(expected, todo)
		getAll(t, expected, server.URL)

		lastIndex := len(expected) - 1
		get(t, lastIndex, expected[lastIndex], server.URL)
	})

	t.Run("Put (replace) a todo", func(t *testing.T) {
		todo := Todo{Name: "go for a walk", Completed: true}
		put(t, 1, todo, server.URL)

		expected[1] = todo
		getAll(t, expected, server.URL)
	})

	t.Run("Delete a specific todo", func(t *testing.T) {
		lastIndex := len(expected) - 1
		del(t, lastIndex, server.URL)

		expected := expected[:lastIndex]
		getAll(t, expected, server.URL)
	})

	t.Run("Delete todos until list is empty", func(t *testing.T) {
		for index := 0; index < len(expected)-1; index++ {
			del(t, 0, server.URL)
		}

		expected = []Todo{}
		getAll(t, expected, server.URL)
	})
}

func post(t *testing.T, todo Todo, baseUrl string) {
	data, err := json.Marshal(todo)
	check(t, err, "json marshalling")

	res, err := http.Post(baseUrl+"/todo", "application/json", bytes.NewBuffer(data))
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Fatalf("Expected status 201 but got %v", res.StatusCode)
	}
}

func put(t *testing.T, index int, todo Todo, baseUrl string) {
	data, err := json.Marshal(todo)
	check(t, err, "json marshalling")

	url := fmt.Sprintf("%s/todo/%d", baseUrl, index)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	check(t, err, "http request creation")

	res, err := http.DefaultClient.Do(req)
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}
}

func del(t *testing.T, index int, baseUrl string) {
	url := fmt.Sprintf("%s/todo/%d", baseUrl, index)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	check(t, err, "http request creation")

	res, err := http.DefaultClient.Do(req)
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}
}

func get(t *testing.T, index int, expected Todo, baseUrl string) {
	res, err := http.Get(fmt.Sprintf("%s/todo/%d", baseUrl, index))
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	check(t, err, "response body reading")

	var actual Todo
	err = json.Unmarshal(data, &actual)
	check(t, err, "response body parsing")

	if expected != actual {
		t.Fatalf("Expected todo %+v but got %+v", expected, actual)
	}
}

func getAll(t *testing.T, expected []Todo, baseUrl string) {
	res, err := http.Get(baseUrl + "/todo/all")
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	check(t, err, "response body reading")

	var actual []Todo
	err = json.Unmarshal(data, &actual)
	check(t, err, "response body parsing")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected todos %+v but got %+v", expected, actual)
	}

	t.Logf("Received todo list: %v", actual)
}

func check(t *testing.T, err error, contextInfo string) {
	t.Helper()
	if err != nil {
		t.Fatalf("Error on %s: %v", contextInfo, err)
	}
}

func welcome(t *testing.T, baseUrl string) {
	res, err := http.Get(baseUrl + "/")
	check(t, err, "http request sending")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		t.Fatalf("Expected status 200 but got %v", res.StatusCode)
	}

	data, err := ioutil.ReadAll(res.Body)
	check(t, err, "response body reading")

	if expected := "Welcome to this Todo list application!"; string(data) != expected {
		t.Fatalf("Expected body '%v' but got '%v'", expected, string(data))
	}
}
