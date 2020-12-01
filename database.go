package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type dpacket struct {
	sid        string
	time       time.Time
	value      int
	attachment []byte
}

// MakePSQLString Make a connection string for PostgreSQL database
func MakePSQLString(dc DatabaseConfigurations) string {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dc.DBHost, dc.DBPort, dc.DBUser, dc.DBPassword, dc.DBName)
	return psqlInfo
}

// GetDB Get database driver
func GetDB(dc DatabaseConfigurations) *sql.DB {
	db, err := sql.Open(dc.DBType, MakePSQLString(dc))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

// WriteData writes sensor data to database.
//	Params:
//		sid: sensor ID
//		time: timestamp
//		value: sensor value
func WriteData(db *sql.DB, dp dpacket) error {
	rootCtx := context.Background()
	ctx, cancel := context.WithTimeout(rootCtx, 1*time.Second)
	defer cancel()

	query := fmt.Sprintf("INSERT INTO data.sdata(sid, time, value, attachment) VALUES ($1, $2, $3, $4)")
	statement, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer statement.Close()
	res, err := statement.ExecContext(ctx, dp.sid, dp.time, dp.value, dp.attachment)
	if err != nil {
		log.Printf("Error %s when inserting row into data.sdata", err)
		return err
	}

	_ = res

	return err
}
