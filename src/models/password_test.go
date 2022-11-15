package models_test

import (
	"testing"

	"github.com/markwallsgrove/muzz_devops/src/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateRandomPassword(t *testing.T) {
	password, err := models.CreateRandomPassword()
	assert.Nil(t, err)
	assert.Len(t, password, 21)
}

func TestCreatePasswordHash(t *testing.T) {
	tests := []struct {
		password string
	}{
		{password: ""},
		{password: "adfsdfsdfvb"},
		{password: "23434dfdslkjf£$%£$5ldfjhsljf"},
	}

	for _, test := range tests {
		hash, err := models.CreatePasswordHash(test.password)
		assert.Nil(t, err)
		assert.Len(t, hash, 32)
	}

}
