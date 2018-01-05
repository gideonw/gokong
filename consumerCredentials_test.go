package gokong

import (
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ConsumerCredentialGetJWTByID(t *testing.T) {
	consumerRequest := &ConsumerRequest{
		Username: "username-" + uuid.NewV4().String(),
		CustomID: "test-" + uuid.NewV4().String(),
	}

	client := NewClient(NewDefaultConfig())
	createdConsumer, err := client.Consumers().Create(consumerRequest)

	assert.Nil(t, err)
	assert.NotNil(t, createdConsumer)

	consumerCredRequest := &ConsumerJWTCredentialRequest{
		ConsumerID: createdConsumer.ID,
	}

	createdCred, err := client.ConsumerCredentials().CreateJWTCredential(consumerCredRequest)

	assert.Nil(t, err)
	assert.NotNil(t, createdCred)

	result, err := client.ConsumerCredentials().GetJWTByID(createdConsumer.Username, createdCred.ID)

	assert.Equal(t, createdCred.ConsumerID, result.ConsumerID)
}

func Test_ConsumerCredentialCreateJWTCredential(t *testing.T) {
	consumerRequest := &ConsumerRequest{
		Username: "username-" + uuid.NewV4().String(),
		CustomID: "test-" + uuid.NewV4().String(),
	}

	client := NewClient(NewDefaultConfig())
	createdConsumer, err := client.Consumers().Create(consumerRequest)

	assert.Nil(t, err)
	assert.NotNil(t, createdConsumer)

	consumerCredRequest := &ConsumerJWTCredentialRequest{
		ConsumerID:   createdConsumer.ID,
		Algorithm:    "RS256",
		RSAPublicKey: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoEG4nvpj/EPGQsAunHuK\nGa/WjDHBaLrTjYH6Z9jxSCINCeHLmn5J4J6zBHlmxK5bpkEsSRgj60LxZZVGvBdt\nKIid+KKEyqTuIWJxAgUd1S/4otUlmGRi0OVji2YgZgRPwRgMv+Ys749jJrhGaN7e\njX6PmAMyP2QGWYGp3tMXDG3mxJJG2V7jFqrwFMgA2O3xRZVeVFyJyhJuFrvn25fI\n+HkPM0+pOEjgTfSX5nMjVyOj8WH18bBfkYxyrPsEmaYC658U5aK1LRp4tuU00JAC\ngYb2DRGnfh6d1/fPtU99Zz3jsU5lmNe5i5zjddX5EjHfeCKhfjlO6Y1bMN7MleG9\nuwIDAQAB\n-----END PUBLIC KEY-----",
	}

	createdCred, err := client.ConsumerCredentials().CreateJWTCredential(consumerCredRequest)

	assert.Nil(t, err)
	assert.NotNil(t, createdCred)

	result, err := client.ConsumerCredentials().GetJWTByID(createdConsumer.Username, createdCred.ID)

	assert.Equal(t, createdCred.ConsumerID, result.ConsumerID)
	assert.Equal(t, createdCred.RSAPublicKey, result.RSAPublicKey)
}
