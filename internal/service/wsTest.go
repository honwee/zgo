package service

import (
	"net/http"

	"zgo/pkg/ws"
)

type Client struct {
	WsClient *ws.Client
}

func (c *Client) Open(writer http.ResponseWriter, request *http.Request) (*Client, bool) {
	wc := &ws.Client{}
	if client, ok := wc.OpenWs(writer, request); ok {
		c.WsClient = client
		return c, true
	}

	return nil, false
}
