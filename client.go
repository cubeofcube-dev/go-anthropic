package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

type ClientOptions struct {
	ApiKey  string
	BaseUrl string
}

type Client struct {
	apiKey  string
	baseUrl string
	cli     *http.Client
}

const DEFAULT_ANTHROPIC_BASE_URL = "https://api.anthropic.com/v1"

var (
	ANTHROPIC_BASE_URL = getEnv("ANTHROPIC_BASE_URL", DEFAULT_ANTHROPIC_BASE_URL)
	ANTHROPIC_API_KEY  = getEnv("ANTHROPIC_API_KEY", "")
)

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func NewClient(opts ...ClientOptions) *Client {
	for _, opt := range opts {
		if opt.ApiKey != "" {
			ANTHROPIC_API_KEY = opt.ApiKey
		}
		if opt.BaseUrl != "" {
			ANTHROPIC_BASE_URL = opt.BaseUrl
		}
	}
	if ANTHROPIC_API_KEY == "" {
		panic("`ANTHROPIC_API_KEY` is required")
	}
	return &Client{
		apiKey:  ANTHROPIC_API_KEY,
		baseUrl: ANTHROPIC_BASE_URL,
		cli:     &http.Client{},
	}
}

func (c *Client) CreateMessages(mr MessagesRequest) (*MessagesResponse, error) {
	reqData, err := json.Marshal(&mr)
	if err != nil {
		return &MessagesResponse{}, errors.New("[Marshal MessagesRequest] " + err.Error())
	}
	req, err := http.NewRequest("POST", c.baseUrl+"/messages", bytes.NewBuffer(reqData))
	if err != nil {
		return &MessagesResponse{}, errors.New("[HTTP New Request] " + err.Error())
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.cli.Do(req)
	if err != nil {
		return &MessagesResponse{}, errors.New("[HTTP Do Request] " + err.Error())
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)

	// handle error
	if resp.StatusCode != http.StatusOK {
		var errMessages MessagesResponseError
		err = json.Unmarshal(respData, &errMessages)
		if err != nil {
			return &MessagesResponse{}, errors.New("[Unmarshal MessagesResponseError] " + err.Error())
		}
		return &MessagesResponse{}, errors.New("[Anthropic Error] " + errMessages.Error.Message)
	}

	var messages MessagesResponse
	err = json.Unmarshal(respData, &messages)
	if err != nil {
		return &MessagesResponse{}, errors.New("[Unmarshal MessagesResponse] " + err.Error())
	}
	return &messages, nil
}

func (c *Client) CreateMessagesStream(mr MessagesRequest) (*MessagesStream, error) {
	mr.Stream = true

	reqData, err := json.Marshal(&mr)
	if err != nil {
		return &MessagesStream{}, errors.New("[Marshal MessagesRequest] " + err.Error())
	}
	req, err := http.NewRequest("POST", c.baseUrl+"/messages", bytes.NewBuffer(reqData))
	if err != nil {
		return &MessagesStream{}, errors.New("[HTTP New Request] " + err.Error())
	}
	req.Header.Set("X-API-KEY", c.apiKey)
	req.Header.Set("Anthropic-Version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := c.cli.Do(req)
	if err != nil {
		return &MessagesStream{}, errors.New("[HTTP Do Request] " + err.Error())
	}

	// handle error
	if resp.StatusCode != http.StatusOK {
		respData, _ := io.ReadAll(resp.Body)
		var errMessages MessagesResponseError
		err = json.Unmarshal(respData, &errMessages)
		if err != nil {
			return &MessagesStream{}, errors.New("[Unmarshal MessagesResponseError] " + err.Error())
		}
		return &MessagesStream{}, errors.New("[Anthropic Error] " + errMessages.Error.Message)
	}

	return &MessagesStream{
		MessagesStreamReader: &MessagesStreamReader{
			reader:   bufio.NewReader(resp.Body),
			response: resp,
		},
	}, nil
}
