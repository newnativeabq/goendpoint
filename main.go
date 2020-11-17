package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"
)

func main() {
	file, _ := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	log.SetOutput(file)
	config := BuildConfigurations("config", "yml")
	db := GetDB(config.Database)

	handler := http.NewServeMux()
	handler.HandleFunc("/", root)
	handler.HandleFunc("/api/data", dataHandler(db))

	defer db.Close()
	http.ListenAndServe(config.Server.Address, handler)
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Hello World`)
}

func dataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		var dp dpacket
		dp.sid = r.FormValue("sid")
		dp.time = time.Now()
		dp.value, _ = strconv.Atoi(r.FormValue("value"))
		// attachment := r.FormValue("attachment")

		err := WriteData(db, dp)
		if err != nil {
			log.Printf("Error %s. Could not write to database.", err)
		}

		fmt.Fprintf(w, "Data Received", dp.sid, dp.time, dp.value)
	}
}
