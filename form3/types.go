package form3

import (
	"net/http"
	"net/url"
)

type RestClient struct {
	httpClient *http.Client
	baseURL    *url.URL
}

type service struct {
	restClient *RestClient
}

type Account struct {
	service
}
