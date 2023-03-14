package form3

import (
	"net/http"
	"net/url"
	"time"
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
	client *RestClient
}

type AccountsService struct {
	service
}

type AccountType string

const (
	AcctTypeAccounts AccountType = "accounts"
)

type Account struct {
	ID             string             `json:"id"`
	OrganisationID string             `json:"organisation_id"`
	Type           AccountType        `json:"type"`
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	Version        int                `json:"version"`
	CreatedOn      *time.Time         `json:"created_on,omitempty"`
	ModifiedOn     *time.Time         `json:"modified_on,omitempty"`
}

type AccountClassification string

const (
	AcctClassificationPersonal AccountClassification = "Personal"
	AcctClassificationBusiness AccountClassification = "Business"
)

type BankIDCode string

// in real production code We should add here a const for each possible bank ID listed in the doc
const (
	BankIDCodeBelgium BankIDCode = "BE"
	BankIDCodeEstonia BankIDCode = "EE"
	BankIDCodeFrance  BankIDCode = "FR"
)

type BaseCurrency string

const (
	BaseCurrencyEur BaseCurrency = "EUR"
	BaseCurrencyUsd BaseCurrency = "USD"
)

type AccountStatus string

const (
	AcctStatusPending   AccountStatus = "pending"
	AcctStatusConfirmed AccountStatus = "confirmed"
	AcctStatusClosed    AccountStatus = "closed"
)

type CountryCode string

// in real production code We should add here a const for each possible country listed in the doc
const (
	CountryCodeBelgium CountryCode = "BE"
	CountryCodeEstonia CountryCode = "EE"
	CountryCodeFrance  CountryCode = "FR"
)

type AccountAttributes struct {
	AccountClassification   AccountClassification `json:"account_classification,omitempty"`
	AccountMatchingOptOut   bool                  `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string                `json:"account_number,omitempty"`
	AlternativeNames        []string              `json:"alternative_names,omitempty"`
	BankID                  string                `json:"bank_id,omitempty"`
	BankIDCode              BankIDCode            `json:"bank_id_code,omitempty"`
	BaseCurrency            BaseCurrency          `json:"base_currency,omitempty"`
	Bic                     string                `json:"bic,omitempty"`
	Country                 CountryCode           `json:"country,omitempty"`
	Iban                    string                `json:"iban,omitempty"`
	JointAccount            bool                  `json:"joint_account,omitempty"`
	Name                    []string              `json:"name,omitempty"`
	SecondaryIdentification string                `json:"secondary_identification,omitempty"`
	Status                  AccountStatus         `json:"status,omitempty"`
	Switched                bool                  `json:"switched,omitempty"`
}
