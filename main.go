package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/KelpGF/Client-Server-API/database"
	"github.com/KelpGF/Client-Server-API/model"
	_ "github.com/mattn/go-sqlite3"
)

type Response struct {
	USDBRL model.Quotation `json:"USDBRL"`
}

func main() {
	database.StartDb()
	defer database.CloseDb()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, mux)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxRequest, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxRequest, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var data Response
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := database.Insert(ctx, data.USDBRL); err != nil {
		http.Error(w, "insert - "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
