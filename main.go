package main

import (
	"bytes"
	"fmt"
	"io"
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
	log.Printf("Server started.")

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
		dp.attachment, err = getFileFromRequest(r, "file")

		err = WriteData(db, dp)
		if err != nil {
			log.Printf("Error %s. Could not write to database.", err)
		}

		log.Printf("Data Received- SID: %s, TIME: %s, VALUE: %d", dp.sid, dp.time, dp.value)
	}
}

// getFileFromRequest attempts to read file
func getFileFromRequest(r *http.Request, key string) (fileAttachment, error) {
	var newAttachment fileAttachment
	file, header, err := r.FormFile(key)
	defer file.Close()

	if err != nil {
		log.Printf("r.FormFile failed. %s", err)
		return newAttachment, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return newAttachment, err
	}

	newAttachment.file = buf.Bytes()
	newAttachment.header = header
	log.Printf("File detected. Attempting to write with filename %s", header.Filename)

	return newAttachment, err
}
