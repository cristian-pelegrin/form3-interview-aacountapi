package form3

import (
	"os"
	"testing"
)

func TestClient_HelloWorld(t *testing.T) {
	baseUrl := os.Getenv("API_URL")
	client, err := NewRestClient(nil, newRestClientParams{baseUrl})
	if err != nil {
		panic(err)
	}

	if err := client.TestHttpCall(); err != nil {
		panic(err)
	}
}
