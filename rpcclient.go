package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type RPCClient struct {
	client http.Client
}

func NewRPCClient() RPCClient {
	return RPCClient{
		client: http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *RPCClient) do(service, method string, request interface{}, response interface{}) error {
	reqURL := fmt.Sprintf("http://%s/%s", service, method)
	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error request marshaling: %s", err)
	}

	resp, err := c.client.Post(reqURL, "application/json", bytes.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("http error: %s", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(resp)
	if err != nil {
		return fmt.Errorf("error body unmarshaling: %s", err)
	}

	return nil
}
