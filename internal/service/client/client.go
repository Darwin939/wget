package client

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	cli *http.Client
}

func NewClient(timeout time.Duration) *Client {
	return &Client{
		cli: &http.Client{Timeout: timeout},
	}
}

func (c *Client) SendHttp1(method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return c.cli.Do(req)
}
