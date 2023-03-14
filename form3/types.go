package form3

import (
	"net/http"
	"net/url"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RestClient struct {
	httpClient HttpClient
	baseURL    *url.URL
}

type RestClientRequest struct {
	*http.Request
}

type RestClientResponse struct {
	*http.Response
}

type service struct {
	restClient *RestClient
}

type Account struct {
	service
}
