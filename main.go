package main

import (
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

func main() {
	// Set the file name of the configuration file
	viper.SetConfigName("config")

	handler := http.NewServeMux()
	handler.HandleFunc("/", root)
	handler.HandleFunc("/api/data", processdata)

	http.ListenAndServe("0.0.0.0:8080", handler)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Hello World`)
}

func processdata(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Data Received")
}
