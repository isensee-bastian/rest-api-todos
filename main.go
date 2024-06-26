package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Provide some initial example todos.
var allTodos = []Todo{
	{Title: "study for math exam"},
	{Title: "take the trash out"},
	{Title: "watch golang tutorial"},
}

type Todo struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// handleWelcome serves as a way of checking that the API is up and running.
func handleWelcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to this Todo list application!")
}

// handlePost adds a new todo item to the list.
func handlePost(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(writer, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var todo Todo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		handleError(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	allTodos = append(allTodos, todo)
	writer.WriteHeader(http.StatusCreated)
}

// handeGetAll returns all todos, optionally filtered by their completion state.
func handleGetAll(writer http.ResponseWriter, request *http.Request) {
	filteredTodos, err := filterTodos(writer, request.URL.Query())

	if err != nil {
		handleError(writer, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := json.Marshal(filteredTodos)

	if err != nil {
		handleError(writer, "Failed to create response body", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%s", data)
}

// handleGet returns a specific todo item from the list by index.
func handleGet(writer http.ResponseWriter, request *http.Request) {
	index := parseIndex(request.PathValue("index"), writer)
	if index < 0 {
		return
	}

	todo := allTodos[index]
	data, err := json.Marshal(todo)

	if err != nil {
		handleError(writer, "Error on response writing", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(writer, "%s", data)
}

// handlePut replaces a specific todo item with a new one by index.
func handlePut(writer http.ResponseWriter, request *http.Request) {
	index := parseIndex(request.PathValue("index"), writer)
	if index < 0 {
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(writer, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var todo Todo

	err = json.Unmarshal(body, &todo)
	if err != nil {
		handleError(writer, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	allTodos[index] = todo
}

// handleDelete removes a specific todo item from the list by index.
func handleDelete(writer http.ResponseWriter, request *http.Request) {
	index := parseIndex(request.PathValue("index"), writer)
	if index < 0 {
		return
	}

	// Remove the specified element by reslicing. This is fine here since our todo
	// list is expected to be relatively small. It would not be a good idea for large
	// slices though due to performance reasons.
	allTodos = append(allTodos[:index], allTodos[index+1:]...)
}

func parseIndex(rawIndex string, writer http.ResponseWriter) int {
	if rawIndex == "" {
		handleError(writer, "Missing index value in path", http.StatusBadRequest)
		return -1
	}

	index, err := strconv.Atoi(rawIndex)
	if err != nil {
		handleError(writer, "Index value in path must be an integer", http.StatusBadRequest)
		return -1
	}

	if index < 0 || index >= len(allTodos) {
		handleError(writer, "Index value in path must be greater than zero and smaller than todo list length", http.StatusBadRequest)
		return -1
	}

	return index
}

func filterTodos(writer http.ResponseWriter, queryParams url.Values) ([]Todo, error) {
	values, found := queryParams["completed"]

	if found && len(values) > 0 {
		filterValue, err := strconv.ParseBool(values[0])

		if err != nil {
			return nil, fmt.Errorf("Invalid query string, expecting boolean value")
		}

		filtered := []Todo{}
		for _, todo := range allTodos {
			if todo.Completed == filterValue {
				filtered = append(filtered, todo)
			}
		}

		return filtered, nil
	}

	return allTodos, nil
}

func handleError(writer http.ResponseWriter, message string, statusCode int) {
	fmt.Printf("Response: %d - %s", statusCode, message)
	http.Error(writer, message, statusCode)
}

func makeMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Important: Make sure your Go version is 1.22.1 or higher as
	// there have been multiple changes to the mux routing.
	mux.HandleFunc("GET /{$}", handleWelcome) // Note {$} to match exactly root.
	mux.HandleFunc("GET /todo/{index}", handleGet)
	mux.HandleFunc("GET /todo/all", handleGetAll)
	mux.HandleFunc("POST /todo", handlePost)
	mux.HandleFunc("PUT /todo/{index}", handlePut)
	mux.HandleFunc("DELETE /todo/{index}", handleDelete)

	return mux
}

func main() {
	mux := makeMux()
	fmt.Println("Listening for requests...")
	http.ListenAndServe(":8090", mux)
}
