//go:build !integration
// +build !integration

package main

import (
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("tested unit")
}
