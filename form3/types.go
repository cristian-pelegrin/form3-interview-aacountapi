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

type Account struct {
	ID             string             `json:"id"`
	OrganisationID string             `json:"organisation_id"`
	Type           string             `json:"type"`
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	Version        int                `json:"version"`
	CreatedOn      *time.Time         `json:"created_on,omitempty"`
	ModifiedOn     *time.Time         `json:"modified_on,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   string   `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}
