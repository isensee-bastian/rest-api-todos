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
	fmt.Fprintf(writer, "Welcome to this Todo list application!")
}

func handleTodoPost(writer http.ResponseWriter, request *http.Request) {
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
	data, err := json.Marshal(allTodos)

	if err != nil {
		http.Error(writer, "Failed to create response body", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%s", data)
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Important: Make sure your Go version is 1.22.1 or higher as
	// there have been multiple changes to the mux routing.
	mux.HandleFunc("GET /{$}", handleWelcome) // Note {$} to match exactly root.
	mux.HandleFunc("GET /todo/all", handleTodoGetAll)
	mux.HandleFunc("POST /todo", handleTodoPost)

	return mux
}

func main() {
	mux := makeMux()
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", mux)
}
