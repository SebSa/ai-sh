package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiURLv1                 = "https://api.openai.com/v1"
	createCompletionEndpoint = "/completions"
)

// Client is the client for the open ai api.
type Client struct {
	apiURL    string
	authToken string
	httpCli   *http.Client
}

// New creates a new open ai client.
func New(httpCli *http.Client, authToken string) *Client {
	return &Client{
		apiURL:    apiURLv1,
		authToken: authToken,
		httpCli:   httpCli,
	}
}

// CreateCompletion creates a completion for the given model and input.
func (c *Client) CreateCompletion(ctx context.Context, input *CreateCompletionRequest) (*CreateCompletionResponse, error) {
	if err := input.validate(); err != nil {
		return nil, fmt.Errorf("could not validate input: %s", err)
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("could not marshal request body: %s", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.apiURL+createCompletionEndpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err)
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %s", err)
	}
	defer resp.Body.Close()

	var result CreateCompletionResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode response: %s", err)
	}
	return &result, nil
}

func (c *Client) makeRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	res, err := c.httpCli.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		var errRes ErrorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return nil, fmt.Errorf("could not decode error with status code: %d", res.StatusCode)
		}

		if errRes.Error != nil {
			return nil, fmt.Errorf("error: %s", errRes.Error.Message)
		}

		return nil, fmt.Errorf(
			"received an unexpected error from open ai with status code: %d and message: %s",
			res.StatusCode, errRes.Error.Message,
		)
	}
	return res, nil
}
