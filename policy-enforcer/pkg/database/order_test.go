package database

import (
	"errors"
	"testing"
)

func TestSanitizeOrder(t *testing.T) {
	tests := []struct {
		name          string
		order         any
		want          interface{}
		expectedError error
	}{
		{
			name:          "happy path",
			order:         " TYPE asc NULLS last , nAmE DeSC   ",
			want:          "type asc nulls last,name desc",
			expectedError: nil,
		},
		{
			name:          "only field", // not supported
			order:         " Type   ",
			want:          "type",
			expectedError: nil,
		},
		{
			name:          "nil",
			order:         nil,
			want:          nil,
			expectedError: nil,
		},
		{
			name:          "empty",
			order:         "",
			want:          nil,
			expectedError: nil,
		},
		{
			name:          "unsupported order type",
			order:         123,
			want:          nil,
			expectedError: ErrInvalidOrder,
		},
		{
			name:          "invalid order",
			order:         "unknown DeSC",
			want:          nil,
			expectedError: ErrInvalidOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sanitizeOrder(tt.order)

			if !errors.Is(err, tt.expectedError) {
				t.Fatalf("Expected error to be %v, got %v", tt.expectedError, err)
			}
			if tt.want != got {
				t.Fatalf("Expected %v, got %v", tt.want, got)
			}
		})
	}
}
