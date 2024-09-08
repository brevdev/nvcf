package api

import (
	"net/http"
	"time"

	nvcf "github.com/tmc/nvcf-go"
	"github.com/tmc/nvcf-go/option"
)

type Client struct {
	*nvcf.Client
	httpClient *http.Client
}

type Option (*Client)

// Options to create: WithHTTPClient, header manips, BaseURL, env handling

func NewClient(apiKey string, opts ...Option) *Client {
	return &Client{
		Client: nvcf.NewClient(
			option.WithHeader("Content-Type", "application/json"),
			option.WithHeader("Accept", "application/json"),
			option.WithHeader("Authorization", "Bearer "+apiKey),
			option.WithBaseURL("https://api.ngc.nvidia.com/"),
		),
		httpClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}
