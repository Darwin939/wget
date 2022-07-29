package client

import "net/http"

type Client struct {
	cli *http.Client
}
