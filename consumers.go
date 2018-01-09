package gokong

import (
	"encoding/json"
	"fmt"

	"github.com/parnurzeal/gorequest"
)

// ConsumerClient hols the config for the kong cluster
type ConsumerClient struct {
	config *Config
}

// ConsumerRequest is used to request a consumer
type ConsumerRequest struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
}

// Consumer kong dao for consumer
type Consumer struct {
	ID       string `json:"id,omitempty"`
	CustomID string `json:"custom_id,omitempty"`
	Username string `json:"username,omitempty"`
}

// Consumers result type for listing consumers
type Consumers struct {
	Results []*Consumer `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
	Next    string      `json:"next,omitempty"`
}

// ConsumerFilter request object to filter consumers
type ConsumerFilter struct {
	ID       string `url:"id,omitempty"`
	CustomID string `url:"custom_id,omitempty"`
	Username string `url:"username,omitempty"`
	Size     int    `url:"size,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

const ConsumersPath = "/consumers/"

// GetByUsername requests from kong a consumer by their username
func (consumerClient *ConsumerClient) GetByUsername(username string) (*Consumer, error) {
	return consumerClient.GetByID(username)
}

// GetByID requests from kong a consumer by their ID
func (consumerClient *ConsumerClient) GetByID(id string) (*Consumer, error) {
	_, body, errs := gorequest.New().Get(consumerClient.config.HostAddress + ConsumersPath + id).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get consumer, error: %v", errs)
	}

	consumer := &Consumer{}
	err := json.Unmarshal([]byte(body), consumer)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer get response, error: %v", err)
	}

	if consumer.ID == "" {
		return nil, nil
	}

	return consumer, nil
}

// Create a consumer
func (consumerClient *ConsumerClient) Create(consumerRequest *ConsumerRequest) (*Consumer, error) {
	_, body, errs := gorequest.New().Post(consumerClient.config.HostAddress + ConsumersPath).Send(consumerRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create new consumer, error: %v", errs)
	}

	createdConsumer := &Consumer{}
	err := json.Unmarshal([]byte(body), createdConsumer)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer creation response, error: %v", err)
	}

	if createdConsumer.ID == "" {
		return nil, fmt.Errorf("could not create consumer, error: %v", body)
	}

	return createdConsumer, nil
}

// CreateOrUpdate a consumer
func (consumerClient *ConsumerClient) CreateOrUpdate(consumerRequest *ConsumerRequest) (*Consumer, error) {
	var exisitingConsumer *Consumer
	if consumerRequest.CustomID != "" {
		consumer, err := consumerClient.GetByID(consumerRequest.CustomID)
		if err == nil {
			exisitingConsumer = consumer
		}
	} else {
		consumer, err := consumerClient.GetByUsername(consumerRequest.Username)
		if err == nil {
			exisitingConsumer = consumer
		}
	}

	if exisitingConsumer != nil {
		consumerRequest.ID = exisitingConsumer.ID
	}

	_, body, errs := gorequest.New().Put(consumerClient.config.HostAddress + ConsumersPath).Send(consumerRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not create or update new consumer, error: %v", errs)
	}

	createdConsumer := &Consumer{}
	err := json.Unmarshal([]byte(body), createdConsumer)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer creation/update response, error: %v", err)
	}

	if createdConsumer.ID == "" {
		return nil, fmt.Errorf("could not create or update consumer, error: %v", body)
	}

	return createdConsumer, nil
}

// List all consumers
func (consumerClient *ConsumerClient) List() (*Consumers, error) {
	return consumerClient.ListFiltered(nil)
}

// ListFiltered consumers by the given filter criteria
func (consumerClient *ConsumerClient) ListFiltered(filter *ConsumerFilter) (*Consumers, error) {
	address, err := addQueryString(consumerClient.config.HostAddress+ConsumersPath, filter)

	if err != nil {
		return nil, fmt.Errorf("could not build query string for consumer filter, error: %v", err)
	}

	_, body, errs := gorequest.New().Get(address).End()
	if errs != nil {
		return nil, fmt.Errorf("could not get consumers, error: %v", errs)
	}

	consumers := &Consumers{}
	err = json.Unmarshal([]byte(body), consumers)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumers list response, error: %v", err)
	}

	return consumers, nil
}

// DeleteByUsername deletes a consumer with the given username
func (consumerClient *ConsumerClient) DeleteByUsername(username string) error {
	return consumerClient.DeleteByID(username)
}

// DeleteByID deletes a consuimer with the given ID
func (consumerClient *ConsumerClient) DeleteByID(id string) error {
	res, _, errs := gorequest.New().Delete(consumerClient.config.HostAddress + ConsumersPath + id).End()
	if errs != nil {
		return fmt.Errorf("could not delete consumer, result: %v error: %v", res, errs)
	}

	return nil
}

// UpdateByUsername updates a consumer with the given username
func (consumerClient *ConsumerClient) UpdateByUsername(username string, consumerRequest *ConsumerRequest) (*Consumer, error) {
	return consumerClient.UpdateByID(username, consumerRequest)
}

// UpdateByID updates a consumer with the given ID
func (consumerClient *ConsumerClient) UpdateByID(id string, consumerRequest *ConsumerRequest) (*Consumer, error) {
	_, body, errs := gorequest.New().Patch(consumerClient.config.HostAddress + ConsumersPath + id).Send(consumerRequest).End()
	if errs != nil {
		return nil, fmt.Errorf("could not update consumer, error: %v", errs)
	}

	updatedConsumer := &Consumer{}
	err := json.Unmarshal([]byte(body), updatedConsumer)
	if err != nil {
		return nil, fmt.Errorf("could not parse consumer update response, error: %v", err)
	}

	if updatedConsumer.ID == "" {
		return nil, fmt.Errorf("could not update consumer, error: %v", body)
	}

	return updatedConsumer, nil
}
