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

type NewRestClientParams struct {
	BaseUrl string
}

type body struct {
	Data         any    `json:"data"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// NewRestClient returns a RestClient instance.
// If a httpClient is not provided, a default http.Client will be assigned
func NewRestClient(httpClient HttpClient, params NewRestClientParams) (*RestClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseUrl, err := url.ParseRequestURI(params.BaseUrl)
	if err != nil {
		return nil, err
	}

	restClient := &RestClient{
		httpClient: httpClient,
		baseURL:    baseUrl,
	}

	return restClient, nil
}

// GetRequest this method returns a GET request ready to send to the api.
func (c *RestClient) GetRequest(path string) (*RestClientRequest, error) {
	return c.newRequest("GET", path, nil)
}

// PostRequest this method returns a POST request ready to send to the api.
// The data payload is wraps in a correct body format accepted by the api
func (c *RestClient) PostRequest(path string, data any) (*RestClientRequest, error) {
	body := body{Data: data}
	return c.newRequest("POST", path, body)
}

// DeleteRequest this method returns a DELETE request ready to send to the api.
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

// Do send the request to the API and unwrap the response data in the v target if it is sent as a parameter.
// ctx must not be nil
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

	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if len(respData) == 0 { // means is a http.noBody response
		return response, nil
	}
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
