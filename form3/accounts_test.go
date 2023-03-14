package form3

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAccountsService_Create(t *testing.T) {
	mockedHttpClient := mockedHttpClientHandler(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "POST", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/%s", baseFakeUrl, "organisation/accounts"), req.URL.String())

		body := `{"data":{"id":"a","organisation_id":"b","type":"accounts","version":0,"attributes":{"account_classification":"Personal"}}}`
		return mockedResponse(http.StatusCreated, body, nil), nil
	})
	client, err := NewRestClient(mockedHttpClient, NewRestClientParams{BaseUrl: baseFakeUrl})
	service := NewAccountsService(client)

	account := &Account{
		ID:             "a",
		OrganisationID: "b",
		Type:           AcctTypeAccounts,
		Version:        0,
		Attributes: &AccountAttributes{
			AccountClassification: "Personal",
		},
	}

	newAccount, resp, err := service.Create(context.Background(), account)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, resp, "Response should be not nil")
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Response code incorrect")
	assert.NotNil(t, newAccount, "NewAccount should be not nil")
	assert.Equal(t, "a", newAccount.ID, "newAccount.ID incorrect")
	assert.Equal(t, "b", newAccount.OrganisationID, "newAccount.OrganisationID incorrect")
	assert.Equal(t, AcctTypeAccounts, newAccount.Type, "newAccount.Type incorrect")
	assert.Equal(t, 0, newAccount.Version, "newAccount.Version incorrect")
	assert.Equal(t, AcctClassificationPersonal, newAccount.Attributes.AccountClassification, "newAccount.Attributes.AccountClassification incorrect")
}

func TestAccountsService_Get(t *testing.T) {
	mockedHttpClient := mockedHttpClientHandler(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/%s", baseFakeUrl, "organisation/accounts/a1b2c3"), req.URL.String())

		body := `{"data":{"id":"a1b2c3","organisation_id":"b","type":"accounts","version":1}}`
		return mockedResponse(http.StatusOK, body, nil), nil
	})
	client, err := NewRestClient(mockedHttpClient, NewRestClientParams{BaseUrl: baseFakeUrl})
	service := NewAccountsService(client)

	account, resp, err := service.Get(context.Background(), "a1b2c3")

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, resp, "Response should be not nil")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response code incorrect")
	assert.NotNil(t, account, "NewAccount should be not nil")
	assert.Equal(t, "a1b2c3", account.ID, "newAccount.ID incorrect")
	assert.Equal(t, "b", account.OrganisationID, "newAccount.OrganisationID incorrect")
	assert.Equal(t, AcctTypeAccounts, account.Type, "newAccount.Type incorrect")
	assert.Equal(t, 1, account.Version, "newAccount.Version incorrect")
}

func TestAccountsService_Delete(t *testing.T) {
	mockedHttpClient := mockedHttpClientHandler(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "DELETE", req.Method)
		assert.Equal(t, fmt.Sprintf("%s/%s", baseFakeUrl, "organisation/accounts/a1b2c3?version=1"), req.URL.String())

		return mockedResponse(http.StatusNoContent, "", nil), nil
	})
	client, err := NewRestClient(mockedHttpClient, NewRestClientParams{BaseUrl: baseFakeUrl})
	service := NewAccountsService(client)

	resp, err := service.Delete(context.Background(), "a1b2c3", 1)

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, resp, "Response should be not nil")
	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "Response code incorrect")
	assert.Equal(t, resp.Body, http.NoBody)
}
