package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/KelpGF/Client-Server-API/model"
)

type Response struct {
	USDBRL model.Quotation `json:"USDBRL"`
}

func GetQuotation(ctx context.Context) (model.Quotation, error) {
	var data Response

	ctxRequest, cancel := initContext(ctx)
	defer cancel()

	request, err := http.NewRequestWithContext(ctxRequest, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return data.USDBRL, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return data.USDBRL, err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return data.USDBRL, err
	}

	return data.USDBRL, nil
}

func initContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctxRequest, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	go func() {
		<-ctxRequest.Done()
		if ctxRequest.Err() == context.DeadlineExceeded {
			fmt.Println("http request context - ", ctxRequest.Err())
		}
	}()

	return ctxRequest, cancel
}
