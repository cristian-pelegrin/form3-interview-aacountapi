package form3

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type newRestClientParams struct {
	baseUrl string
}

type body struct {
	Data         any    `json:"data"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func NewRestClient(httpClient HttpClient, params newRestClientParams) (*RestClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseUrl, err := url.ParseRequestURI(params.baseUrl)
	if err != nil {
		return nil, err
	}

	restClient := &RestClient{
		httpClient: httpClient,
		baseURL:    baseUrl,
	}

	return restClient, nil
}

func (c *RestClient) GetRequest(path string) (*RestClientRequest, error) {
	return c.newRequest("GET", path, nil)
}

func (c *RestClient) PostRequest(path string, data any) (*RestClientRequest, error) {
	body := body{Data: data}
	return c.newRequest("POST", path, body)
}

func (c *RestClient) DeleteRequest(path string) (*RestClientRequest, error) {
	return c.newRequest("DELETE", path, nil)
}

func (c *RestClient) newRequest(method string, path string, body any) (*RestClientRequest, error) {
	targetURL, err := url.Parse(
		fmt.Sprintf("%s/%s", c.baseURL.String(), strings.TrimPrefix(path, "/")),
	)
	if err != nil {
		return nil, err
	}

	var payload io.Reader
	if body != nil {
		dataJson, _ := json.Marshal(body)
		payload = bytes.NewBuffer(dataJson)
	}

	req, err := http.NewRequest(method, targetURL.String(), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return &RestClientRequest{req}, nil
}

func (c *RestClient) Do(ctx context.Context, req *http.Request, v any) (*RestClientResponse, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	response := &RestClientResponse{resp}

	if resp.Body == nil {
		return response, nil
	}

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}

	body := &body{}
	if err = json.Unmarshal(respData, body); err != nil {
		return response, err
	}
	if body.ErrorMessage != "" {
		return response, errors.New(body.ErrorMessage)
	}

	if v != nil {
		bodyDataCoded, err := json.Marshal(body.Data)
		if err != nil {
			return response, err
		}
		err = json.Unmarshal(bodyDataCoded, v)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}
