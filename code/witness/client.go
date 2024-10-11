package main

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	http.Client

    url string
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Do(ctx, context.Context, req rpc.JSONRequest) (rpc.JSONResponse, error) {
    rpc.DialOptions(ctx, c.url)
}
