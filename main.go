package main

import (
	"fmt"
	"net/http"
)

func handleWelcome(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Welcome to this Todo list application!\n")
}

func main() {
	http.HandleFunc("/", handleWelcome)
	http.ListenAndServe(":8090", nil)
}
