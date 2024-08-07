package api

import (
	"net/http"
	"time"

	"github.com/your-username/nvcf-cli/config"
	"github.com/your-username/nvcf-cli/nvcf"
)

type Client struct {
	nvcfClient *nvcf.Client
	httpClient *http.Client
}

func NewClient() *Client {
	apiKey := config.GetAPIKey()
	return &Client{
		nvcfClient: nvcf.NewClient(nvcf.WithAPIKey(apiKey)),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Implement API methods (ListFunctions, CreateFunction, etc.) here
