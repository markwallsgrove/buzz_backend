//go:build !integration
// +build !integration

package domain_test

import (
	"testing"

	"github.com/markwallsgrove/muzz_devops/src/models/domain"
	"github.com/stretchr/testify/assert"
)

func TestGender(t *testing.T) {
	tests := []struct {
		gender domain.Gender
		str    string
	}{
		{gender: domain.Female, str: "Female"},
		{gender: domain.UnknownGender, str: "Unknown"},
		{gender: domain.Male, str: "Male"},
	}

	for _, test := range tests {
		assert.Equal(t, test.gender.String(), test.str)
	}
}
