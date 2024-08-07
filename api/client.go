package api

import (
	"fmt"
	"net/http"
	"time"

	nvcf "github.com/stainless-sdks/nvcf-go"
	"github.com/tmc/nvcf/config"
)

type Client struct {
	nvcfClient *nvcf.Client
	httpClient *http.Client
}

func NewClient() *Client {
	apiKey := config.GetAPIKey()
	fmt.Println("apiKey:", apiKey)
	return &Client{
		nvcfClient: nvcf.NewClient(),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Implement API methods (ListFunctions, CreateFunction, etc.) here
