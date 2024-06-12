package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/KelpGF/Client-Server-API/server/model"
	"github.com/google/uuid"
)

var db *sql.DB

func StartDb() error {
	log.Println("Starting database...")

	sqlite3, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return err
	}

	db = sqlite3

	db.Exec("CREATE TABLE IF NOT EXISTS quotation (id TEXT PRIMARY KEY, code TEXT, codein TEXT, name TEXT, high TEXT, low TEXT, varBid TEXT, pctChange TEXT, bid TEXT, ask TEXT, timestamp TEXT, create_date TEXT)")

	return nil
}

func CloseDb() {
	db.Close()
}

func Insert(ctx context.Context, quotation model.Quotation) error {
	ctxDB, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	go func() {
		<-ctxDB.Done()
		if ctxDB.Err() == context.DeadlineExceeded {
			fmt.Println("database context - ", ctxDB.Err())
		}
	}()

	stmt, err := db.PrepareContext(ctxDB, "INSERT INTO quotation (id, code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(
		ctxDB,
		uuid.New().String(),
		quotation.Code,
		quotation.Codein,
		quotation.Name,
		quotation.High,
		quotation.Low,
		quotation.VarBid,
		quotation.PctChange,
		quotation.Bid,
		quotation.Ask,
		quotation.Timestamp,
		quotation.CreateDate,
	)
	if err != nil {
		return err
	}

	return nil
}
