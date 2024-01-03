package model_test

import (
	"testing"

	"go-dk/model"
)

func TestEmail_IsValid(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"me@example.com", true},
		{"@example.com", false},
		{"me@", false},
		{"@", false},
		{"", false},
	}
	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			if got := model.Email(test.address).IsValid(); got != test.valid {
				t.Errorf("IsValid() = %v, want %v", got, test.valid)
			}
		})
	}
}
