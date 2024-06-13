package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/KelpGF/Client-Server-API/server/database"
	"github.com/KelpGF/Client-Server-API/server/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.StartDb()
	defer database.CloseDb()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Handler)

	HOST := os.Getenv("HOST")
	PORT := os.Getenv("PORT")
	if HOST == "" {
		HOST = "localhost"
	}
	if PORT == "" {
		PORT = "8080"
	}

	log.Println("Server started on port " + PORT)
	http.ListenAndServe(HOST+":"+PORT, mux)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := service.GetQuotation(ctx)
	if err != nil {
		http.Error(w, "get - "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = database.Insert(ctx, data)
	if err != nil {
		http.Error(w, "insert - "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
