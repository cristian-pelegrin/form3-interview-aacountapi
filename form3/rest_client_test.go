package form3

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

const baseFakeUrl = "https://www.fake-api.com/v1"

type mockedHttpClientHandler func(req *http.Request) (*http.Response, error)

func (handler mockedHttpClientHandler) Do(req *http.Request) (*http.Response, error) {
	return handler(req)
}

func mockedResponse(statusCode int, body string, header http.Header) *http.Response {
	if header == nil {
		header = make(http.Header)
	}

	response := &http.Response{
		StatusCode: statusCode,
		Header:     header,
	}
	if body != "" {
		response.Body = io.NopCloser(bytes.NewBufferString(body))
	}

	return response
}

type testRequestExpected struct {
	method string
	path   string
	body   string
}

func testRequest(t *testing.T, r *http.Request, expected testRequestExpected) {
	t.Helper()

	assert.Equal(t, expected.method, r.Method, "Request method expected eror")

	expectedUrl := fmt.Sprintf("%s/%s", baseFakeUrl, strings.TrimPrefix(expected.path, "/"))
	assert.Equal(t, expectedUrl, r.URL.String(), "Request URL expected error")

	textBody := ""
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Error reading Request Body: %v", err)
		}
		textBody = string(body)
	}

	assert.Equal(t, expected.body, textBody, "Request body expected error")
}

func TestRestClient_GET(t *testing.T) {
	c, _ := NewRestClient(nil, newRestClientParams{baseUrl: baseFakeUrl})
	testPath := "/test-get-path"
	req, _ := c.GetRequest(testPath)
	testRequest(t, req.Request, testRequestExpected{method: "GET", path: testPath, body: ""})
}

func TestRestClient_POST(t *testing.T) {
	c, _ := NewRestClient(nil, newRestClientParams{baseUrl: baseFakeUrl})
	testPath := "test-post-path"
	testData := struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}{Id: "1", Name: "cristian"}
	req, _ := c.PostRequest(testPath, testData)
	testRequest(t, req.Request, testRequestExpected{method: "POST", path: testPath, body: `{"data":{"id":"1","name":"cristian"}}`})
}

func TestRestClient_DELETE(t *testing.T) {
	c, _ := NewRestClient(nil, newRestClientParams{baseUrl: baseFakeUrl})
	testPath := "/test-delete-path/123456789"
	req, _ := c.DeleteRequest(testPath)
	testRequest(t, req.Request, testRequestExpected{method: "DELETE", path: testPath, body: ""})
}

func TestRestClient_Do_successResponse(t *testing.T) {
	type testResponse struct {
		A int `json:"a"`
		B struct {
			C string `json:"c"`
		} `json:"b"`
	}

	mockedHttpClient := mockedHttpClientHandler(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, fmt.Sprintf("%s/%s", baseFakeUrl, "foo"), req.URL.String())

		return mockedResponse(
			http.StatusOK,
			`{ "data": { "a": 1, "b": { "c": "bar" } } }`,
			nil,
		), nil
	})
	client, err := NewRestClient(mockedHttpClient, newRestClientParams{baseUrl: baseFakeUrl})
	if err != nil {
		t.Fatalf("Error creting RestClient: %v", err)
	}

	req, err := client.GetRequest("foo")
	if err != nil {
		t.Fatalf("Error creating GET request: %v", err)
	}

	var v = new(testResponse)
	resp, err := client.Do(context.Background(), req.Request, v)
	if err != nil {
		t.Fatalf("Error calling RestClient Do: %v", err)
	}

	assert.Nil(t, err, "RestClient Error is not nil")
	assert.NotNil(t, resp, "RestClient response is nil")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Error in response status expected")
	assert.Equal(t, 1, v.A, "Error expected response field")
	assert.Equal(t, "bar", v.B.C, "Error expected response field")
}

func TestRestClient_Do_errorResponse(t *testing.T) {
	type testResponse struct {
		A int `json:"a"`
	}

	mockedHttpClient := mockedHttpClientHandler(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, fmt.Sprintf("%s/%s", baseFakeUrl, "foo"), req.URL.String())

		return mockedResponse(
			http.StatusNotFound,
			`{ "error_message": "record 123 does not exist" }`,
			nil,
		), nil
	})
	client, err := NewRestClient(mockedHttpClient, newRestClientParams{baseUrl: baseFakeUrl})
	if err != nil {
		t.Fatalf("Error creting RestClient: %v", err)
	}

	req, err := client.GetRequest("foo")
	if err != nil {
		t.Fatalf("Error creating GET request: %v", err)
	}

	var v = new(testResponse)
	resp, err := client.Do(context.Background(), req.Request, v)

	assert.NotNil(t, resp, "RestClient response should be nil")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Error in response status expected")
	assert.NotNil(t, err, "RestClient error should be not nil")
	assert.Equal(t, "record 123 does not exist", err.Error(), "Error message expected")
	assert.Empty(t, 0, v, "Target response object should be empty")
}
