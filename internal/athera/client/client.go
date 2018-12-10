package client

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/oauth2"
)

// Client ...
type Client struct {
	baseURL *url.URL
	tok     *oauth2.Token

	httpClient *http.Client
}

// New ...
func New(url *url.URL, tok *oauth2.Token) *Client {

	timeout, err := time.ParseDuration("30s")

	if err != nil {
		fmt.Println("Failed parsing timeout:", err)
		os.Exit(1)
	}

	c := &Client{
		baseURL: url,
		tok:     tok,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}

	return c
}

func isSuccessStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
