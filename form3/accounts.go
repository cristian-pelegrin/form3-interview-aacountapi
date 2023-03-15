package form3

import (
	"context"
	"fmt"
)

const accountsBasePath = "organisation/accounts"

// NewAccountsService returns a AccountsService instance.
func NewAccountsService(client *RestClient) *AccountsService {
	return &AccountsService{
		service{
			client,
		},
	}
}

// Create creates a new account and returns it.
func (s *AccountsService) Create(ctx context.Context, data *Account) (*Account, *RestClientResponse, error) {
	req, err := s.client.PostRequest(accountsBasePath, data)
	if err != nil {
		return nil, nil, err
	}

	account := new(Account)
	resp, err := s.client.Do(ctx, req.Request, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, nil
}

// Get retrieves an account by its id
func (s *AccountsService) Get(ctx context.Context, id string) (*Account, *RestClientResponse, error) {
	path := fmt.Sprintf("%s/%s", accountsBasePath, id)

	req, err := s.client.GetRequest(path)
	if err != nil {
		return nil, nil, err
	}

	account := new(Account)
	resp, err := s.client.Do(ctx, req.Request, account)
	if err != nil {
		return nil, resp, err
	}

	return account, resp, nil
}

// Delete deletes an account by its id and version
func (s *AccountsService) Delete(ctx context.Context, id string, version int) (*RestClientResponse, error) {
	url := fmt.Sprintf("%s/%s?version=%d", accountsBasePath, id, version)

	req, err := s.client.DeleteRequest(url)
	if err != nil {
		return nil, nil
	}

	resp, err := s.client.Do(ctx, req.Request, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
