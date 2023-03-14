package integration

import (
	"context"
	"fmt"
	"form3-interview-accountapi/form3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestAccountsService_Create(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}
	account := generateTestAccount()

	newAccount, resp, err := service.Create(context.Background(), account)

	assert.Nil(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// assert sent values
	assert.NotNil(t, newAccount)
	assert.NotNil(t, newAccount.CreatedOn)
	assert.NotNil(t, newAccount.ModifiedOn)
	assert.Equal(t, account.ID, newAccount.ID)
	assert.Equal(t, account.OrganisationID, newAccount.OrganisationID)
	assert.Equal(t, account.Attributes.AccountClassification, newAccount.Attributes.AccountClassification)
	assert.Equal(t, account.Attributes.AlternativeNames, newAccount.Attributes.AlternativeNames)
	assert.Equal(t, account.Attributes.BankID, newAccount.Attributes.BankID)
	assert.Equal(t, account.Attributes.BankIDCode, newAccount.Attributes.BankIDCode)
	assert.Equal(t, account.Attributes.BaseCurrency, newAccount.Attributes.BaseCurrency)
	assert.Equal(t, account.Attributes.Country, newAccount.Attributes.Country)
	assert.Equal(t, account.Attributes.Name, newAccount.Attributes.Name)
	assert.Equal(t, account.Attributes.Status, newAccount.Attributes.Status)

	// assert default values expected
	assert.Equal(t, 0, newAccount.Version)
	assert.Equal(t, false, newAccount.Attributes.AccountMatchingOptOut)
	assert.Equal(t, false, newAccount.Attributes.JointAccount)
	assert.Equal(t, false, newAccount.Attributes.Switched)
	assert.Empty(t, newAccount.Attributes.AccountNumber)
	assert.Empty(t, newAccount.Attributes.Bic)
	assert.Empty(t, newAccount.Attributes.Iban)
	assert.Empty(t, newAccount.Attributes.SecondaryIdentification)
}

func TestAccountsService_Create_ErrorResponse(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	account := generateTestAccount()
	account.ID = "invalid-UUID"

	newAccount, resp, err := service.Create(context.Background(), account)

	assert.Nil(t, newAccount)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id in body must be of type uuid")
}

func TestAccountsService_Get(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	account := generateTestAccount()
	ctx := context.Background()
	newAccount, _, _ := service.Create(ctx, account)

	retrievedAcct, resp, err := service.Get(ctx, newAccount.ID)

	assert.Nil(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.Equal(t, newAccount.ID, retrievedAcct.ID)
	assert.Equal(t, newAccount.CreatedOn, retrievedAcct.CreatedOn)
	assert.Equal(t, newAccount.ModifiedOn, retrievedAcct.ModifiedOn)
}

func TestAccountsService_Get_ErrorResponse(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	retrievedAcct, resp, err := service.Get(context.Background(), "invalid-UUID")

	assert.Nil(t, retrievedAcct)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id is not a valid uuid")
}

func TestAccountsService_Get_ErrorNotFound(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	notStoredId := uuid.New().String()
	retrievedAcct, resp, err := service.Get(context.Background(), notStoredId)

	assert.Nil(t, retrievedAcct)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("record %s does not exist", notStoredId))
}

func TestAccountsService_Delete(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	account := generateTestAccount()
	ctx := context.Background()
	_, _, err = service.Create(ctx, account)
	assert.Nil(t, err)

	resp, err := service.Delete(ctx, account.ID, 0)

	assert.Nil(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func TestAccountsService_Delete_ErrorNotFound(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	resp, err := service.Delete(context.Background(), uuid.New().String(), 0)

	assert.Nil(t, err)

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAccountsService_Delete_ErrorInvalidVersion(t *testing.T) {
	service, err := getNewAccountsService()
	if err != nil {
		t.Fatalf("Error creating AccountsService: %v", err)
	}

	account := generateTestAccount()
	ctx := context.Background()
	_, _, err = service.Create(ctx, account)
	assert.Nil(t, err)

	resp, err := service.Delete(ctx, account.ID, 3)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "invalid version")

	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func getNewAccountsService() (*form3.AccountsService, error) {
	baseUrl := os.Getenv("API_URL")
	client, err := form3.NewRestClient(nil, form3.NewRestClientParams{BaseUrl: baseUrl})
	if err != nil {
		return nil, err
	}

	return form3.NewAccountsService(client), nil
}

func generateTestAccount() *form3.Account {
	return &form3.Account{
		ID:             uuid.New().String(),
		OrganisationID: uuid.New().String(),
		Type:           form3.AcctTypeAccounts,
		Attributes: &form3.AccountAttributes{
			AccountClassification: form3.AcctClassificationPersonal,
			AlternativeNames:      []string{"foo", "bar"},
			BankID:                "ZXE",
			BankIDCode:            form3.BankIDCodeBelgium,
			BaseCurrency:          form3.BaseCurrencyEur,
			Country:               form3.CountryCodeBelgium,
			Name:                  []string{"cristian", "pelegrin"},
			Status:                form3.AcctStatusPending,
		},
	}
}
