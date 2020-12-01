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
	fmt.Println("Server started.")

	config := BuildConfigurations("config", "yml")
	db := GetDB(config.Database)

	handler := http.NewServeMux()
	handler.HandleFunc("/", root)
	handler.HandleFunc("/api/data/", dataHandler(db))

	defer db.Close()
	log.Fatal(http.ListenAndServe(config.Server.Address, handler))
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `Hello World`)
}

func dataHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dp dpacket
		err := r.ParseMultipartForm(0)
		if err != nil {
			log.Printf("Error %s.  Could not parse MultipartForm.", err)
		}

		file, fh, err := r.FormFile("files")
		_ = fh

		if err != nil {
			log.Printf("r.FormFile failed. %s", err)
		} else {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, file); err != nil {
				log.Printf("Error %s. Could not copy file contents to byte buffer.", err)
			}
			dp.attachment = buf.Bytes()
			file.Close()
		}

		dp.sid = r.FormValue("sid")
		dp.time = time.Now()
		dp.value, _ = strconv.Atoi(r.FormValue("value"))

		err = WriteData(db, dp)
		if err != nil {
			log.Printf("Error %s. Could not write to database.", err)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("201 - Successfully Logged Data with DB"))
	}
}
