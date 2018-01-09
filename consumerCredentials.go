package gokong

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

// ConsumerCredentialClient hols the config for the kong cluster
type ConsumerCredentialClient struct {
	config *Config
}

// ConsumerJWTCredentialRequest request object used to request a consumer
type ConsumerJWTCredentialRequest struct {
	CreatedAt    int    `json:"created_at,omitempty"`
	ID           string `json:"id,omitempty"`
	ConsumerID   string `json:"consumer_id"`
	Key          string `json:"key,omitempty"`
	Secret       string `json:"secret,omitempty"`
	RSAPublicKey string `json:"rsa_public_key,omitempty"`
	Algorithm    string `json:"algorithm,omitempty"`
}

// ConsumerJWTCredential kong dao for consumer JWT credential
type ConsumerJWTCredential struct {
	CreatedAt    int    `json:"created_at"`
	ID           string `json:"id"`
	ConsumerID   string `json:"consumer_id"`
	Key          string `json:"key,omitempty"`
	Secret       string `json:"secret,omitempty"`
	RSAPublicKey string `json:"rsa_public_key,omitempty"`
	Algorithm    string `json:"algorithm,omitempty"`
}

// ConsumerJWTCredentials result type for listing consumers
type ConsumerJWTCredentials struct {
	Results []*ConsumerJWTCredential `json:"data,omitempty"`
	Total   int                      `json:"total,omitempty"`
	Next    string                   `json:"next,omitempty"`
}

const ConsumerJWTCredentialPath = "/jwt/"
const JWTCredentialPath = "/jwts/"

// CreateJWTCredential for a consumer
func (client *ConsumerCredentialClient) CreateJWTCredential(credRequest *ConsumerJWTCredentialRequest) (*ConsumerJWTCredential, error) {
	address := client.config.HostAddress + ConsumersPath + credRequest.ConsumerID + ConsumerJWTCredentialPath
	_, body, errs := gorequest.New().Post(address).Send(credRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new consumer JWT credential, error: %v", errs)
	}

	createdCred := &ConsumerJWTCredential{}
	err := json.Unmarshal([]byte(body), createdCred)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer JWT credential creation response, error: %v", err)
	}

	if createdCred.ID == "" {
		return nil, fmt.Errorf("could not create consumer JWT credential, error: %v", body)
	}

	return createdCred, nil
}

// GetJWTByID requests from kong a consumer JWT credential by ID
func (client *ConsumerCredentialClient) GetJWTByID(consumerID, id string) (*ConsumerJWTCredential, error) {
	_, body, errs := gorequest.New().Get(client.config.HostAddress + ConsumersPath + consumerID + ConsumerJWTCredentialPath + id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get consumer JWT credential, error: %v", errs)
	}

	cred := &ConsumerJWTCredential{}
	err := json.Unmarshal([]byte(body), cred)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer JWT credential get response, error: %v", err)
	}

	if cred.ID == "" {
		return nil, nil
	}

	return cred, nil
}

// DeleteJWTByID deletes from kong a consumer JWT credential by ID
func (client *ConsumerCredentialClient) DeleteJWTByID(consumerID, id string) error {
	res, _, errs := gorequest.New().Delete(client.config.HostAddress + ConsumersPath + consumerID + ConsumerJWTCredentialPath + id).End()
	if errs != nil {
		return fmt.Errorf("could not delete consumer JWT credential, result: %v error: %v", res, errs)
	}

	switch res.StatusCode {
	case http.StatusNoContent: // success
		fallthrough
	case http.StatusNotFound: // didn't exsist so consider it deleted
		return nil
	default:
		return fmt.Errorf("could not delete consumer JWT credential, result: %+v", res)
	}
}

// ListJWTOnConsumer returns a list of all of the jwt credentials on a given consumer
func (client *ConsumerCredentialClient) ListJWTOnConsumer(consumerID string) (*ConsumerJWTCredentials, error) {
	_, body, errs := gorequest.New().Get(client.config.HostAddress + ConsumersPath + consumerID + ConsumerJWTCredentialPath).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get list of consumer JWT credentials, error: %v", errs)
	}

	cred := &ConsumerJWTCredentials{}
	err := json.Unmarshal([]byte(body), cred)
	if err != nil {
		return nil, fmt.Errorf("could not parse list consumer JWT credential get response, error: %v", err)
	}

	return cred, nil
}
