package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var allTodos = []Todo{
	{Name: "study for math exam"},
	{Name: "take the trash out"},
	{Name: "watch golang tutorial"},
}

type Todo struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func handleWelcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to this Todo list application!\n")
}

func handleTodoPost(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		writer.Header().Set("Allow", "POST")
		writer.WriteHeader(405) // Method not allowed.
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var todo Todo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	allTodos = append(allTodos, todo)

	writer.WriteHeader(http.StatusCreated)
}

func handleTodoGetAll(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.Header().Set("Allow", "GET")
		writer.WriteHeader(405) // Method not allowed.
		return
	}

	data, err := json.Marshal(allTodos)

	if err != nil {
		http.Error(writer, "Failed to create response body", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%s", data)
}

func main() {
	http.HandleFunc("/", handleWelcome)
	http.HandleFunc("/todo/all", handleTodoGetAll)
	http.HandleFunc("/todo", handleTodoPost)

	http.ListenAndServe(":8090", nil)
}
