package form3

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type newRestClientParams struct {
	baseUrl string
}

func NewRestClient(httpClient *http.Client, params newRestClientParams) (*RestClient, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	baseUrl, err := url.Parse(params.baseUrl)
	if err != nil {
		return nil, err
	}

	restClient := &RestClient{
		httpClient: httpClient,
		baseURL:    baseUrl,
	}

	return restClient, nil
}

func (c RestClient) TestHttpCall() error {
	req, err := http.NewRequest("GET", c.baseURL.String()+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(respData))

	return nil
}
