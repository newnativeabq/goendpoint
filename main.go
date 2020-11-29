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
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Error %s.  Could not parse MultipartForm.", err)
		}

		var dp dpacket
		dp.sid = r.FormValue("sid")
		dp.time = time.Now()
		dp.value, _ = strconv.Atoi(r.FormValue("value"))
		dp.attachment = getFileFromRequest(r, "file")

		err = WriteData(db, dp)
		if err != nil {
			log.Printf("Error %s. Could not write to database.", err)
		}

		log.Printf("Data Received- SID: %s, TIME: %s, VALUE: %d", dp.sid, dp.time, dp.value)
	}
}

// getFileFromRequest attempts to read file
func getFileFromRequest(r *http.Request, key string) fileAttachment {
	var newAttachment fileAttachment
	file, header, err := r.FormFile(key)

	if err != nil {
		log.Printf("r.FormFile failed. %s", err)
	} else {
		newAttachment.file = file
		newAttachment.header = header
		log.Printf("File detected. Attempting to write with filename %s", header.Filename)
	}

	return newAttachment
}
