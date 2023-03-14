package form3

import (
	"context"
	"fmt"
)

const accountsBasePath = "organisation/accounts"

func NewAccountsService(client *RestClient) *AccountsService {
	return &AccountsService{
		service{
			client,
		},
	}
}

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
