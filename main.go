package main

import (
	"fmt"
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	handler.HandleFunc("/", root)
	handler.HandleFunc("/api/data", processdata)

	configuration := BuildConfigurations("config", "yml")
	http.ListenAndServe(configuration.Server.Address, handler)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Hello World`)
}

func processdata(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Data Received")
}
