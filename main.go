package main

import (
	"fmt"
	"net/http"

	"database/sql"
)

func main() {
	config := BuildConfigurations("config", "yml")
	db := GetDB(config.Database)

	handler := http.NewServeMux()
	handler.HandleFunc("/", root)
	handler.HandleFunc("/api/data", dataHandler(db))

	http.ListenAndServe(config.Server.Address, handler)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Hello World`)
}

func dataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		sid := r.FormValue("sid")
		time := r.FormValue("time")
		value := r.FormValue("value")
		// attachment := r.FormValue("attachment")

		fmt.Fprintf(w, "Data Received", sid, time, value)
	}
}
