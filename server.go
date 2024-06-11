package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/KelpGF/Client-Server-API/database"
	"github.com/KelpGF/Client-Server-API/service"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database.StartDb()
	defer database.CloseDb()

	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", Handler)

	port := os.Getenv("PORT")
	log.Println("Server started on port " + port)
	http.ListenAndServe(":"+port, mux)
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
