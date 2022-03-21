package currencyLayer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"CleverIT-challenge/internal/core/domain/currency"
)

// currencyClientHTTP is the http client for communication with Currency Layer Service
type currencyClientHTTP struct {
	url string
}

type Response struct {
	Success bool               `json:"success"`
	Quotes  map[string]float64 `json:"quotes"`
	Error struct{
		Info    string             `json:"info,omitempty"`
	} `json:"error"`
}

func (client *currencyClientHTTP) GetCurrencyPriceInDollar(_ context.Context, currencyID string) (float64, error) {
	currencyURL := client.url + currencyID
	result, err := http.Get(currencyURL)
	if err != nil {
		return 0, err
	}
	bytes, err := io.ReadAll(result.Body)
	if err != nil {
		return 0, err
	}
	response := &Response{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return 0, err
	}
	if !response.Success {
		return 0, fmt.Errorf(response.Error.Info)
	}
	currencyKey := "USD" + currencyID
	for key, value := range response.Quotes {
		if key == currencyKey {
			return value, nil
		}
	}
	log.Println("Error finding curency: ", currencyKey)
	return 0, fmt.Errorf("invalid currencyKey")
}

func (client *currencyClientHTTP) IsValidCurrency(_ context.Context, currencyID string) (bool, error) {
	result, err := http.Get(client.url + currencyID)
	if err != nil {
		return false, err
	}
	bytes, err := io.ReadAll(result.Body)
	if err != nil {
		return false, err
	}
	response := &Response{}
	if err = json.Unmarshal(bytes, response); err != nil {
		return false, err
	}
	if !response.Success {
		return false, fmt.Errorf(response.Error.Info)
	}

	for key, _ := range response.Quotes {
		if key == "USD"+currencyID {
			return true, nil
		}
	}
	log.Println("Error finding curency: ", currencyID)
	return false, nil
}

func NewCurrencyService() currency.Service {
	var accessToken = os.Getenv("CURRENCY_LAYER_TOKEN")

	return &currencyClientHTTP{
		url: fmt.Sprintf("http://api.currencylayer.com/live?access_key=%s&format=1&currencies=", accessToken),
	}
}
