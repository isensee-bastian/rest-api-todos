package main

import (
	"encoding/json"
	"fmt"
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

func handleTodoAll(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		writer.Header().Set("Allow", "GET")
		writer.WriteHeader(405) // Method not allowed.
		return
	}

	data, err := json.Marshal(allTodos)

	if err != nil {
		fmt.Println("Error on json marshalling:", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%s", data)
}

func main() {
	http.HandleFunc("/", handleWelcome)
	http.HandleFunc("/todo/all", handleTodoAll)

	http.ListenAndServe(":8090", nil)
}
