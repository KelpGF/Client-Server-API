package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/KelpGF/Client-Server-API/client/model"
)

func main() {
	data, err := GetQuotation(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	writeQuotation(data.Bid)
}

func GetQuotation(ctx context.Context) (model.Quotation, error) {
	var data model.Quotation

	ctxRequest, cancel := initContext(ctx)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxRequest, "GET", "http://go-server:8080/cotacao", nil)
	if err != nil {
		return data, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return data, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func initContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctxRequest, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	go func() {
		<-ctxRequest.Done()
		if ctxRequest.Err() == context.DeadlineExceeded {
			fmt.Println("http request context - ", ctxRequest.Err())
		}
	}()

	return ctxRequest, cancel
}

func writeQuotation(bid string) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Panicf("Failed to create file: %v", err)
	}

	text := "Dólar: " + bid
	size, err := f.WriteString(text)
	if err != nil {
		log.Panicf("Failed to write file: %v", err)
	}

	fmt.Printf("Wrote %d bytes\n", size)
}
