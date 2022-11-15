//go:build !integration
// +build !integration

package models_test

import (
	"testing"

	"github.com/markwallsgrove/muzz_devops/src/models"
	"github.com/stretchr/testify/assert"
)

func TestGender(t *testing.T) {
	tests := []struct {
		gender models.Gender
		str    string
	}{
		{gender: models.Female, str: "Female"},
		{gender: models.UnknownGender, str: "Unknown"},
		{gender: models.Male, str: "Male"},
	}

	for _, test := range tests {
		assert.Equal(t, test.gender.String(), test.str)
	}
}
