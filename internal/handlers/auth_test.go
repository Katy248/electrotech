package handlers

import (
	"context"
	"testing"
)

//nolint:containedctx
func TestGetUserID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		ctxToCheck    context.Context
		expectedID    int64
		valueExpected bool
	}{
		{
			name:          "basic/ok",
			ctxToCheck:    context.WithValue(t.Context(), userIDKey, int64(12)),
			expectedID:    12,
			valueExpected: true,
		},
		{
			name:          "basic/bad-int",
			ctxToCheck:    context.WithValue(t.Context(), userIDKey, 12),
			expectedID:    12,
			valueExpected: false,
		},
		{
			name:          "no-value",
			ctxToCheck:    t.Context(),
			valueExpected: false,
			expectedID:    0,
		},
		{
			name:          "bad-value",
			ctxToCheck:    context.WithValue(t.Context(), userIDKey, "12"),
			valueExpected: false,
			expectedID:    0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, exists := GetUserID(tt.ctxToCheck)
			if !tt.valueExpected && exists {
				t.Errorf("Value not expected, but exists")
			}

			if tt.valueExpected && !exists {
				t.Errorf("Value expected, but not exists")
			}

			if tt.valueExpected && tt.expectedID != id {
				t.Errorf("Expected %d, but returned %d", tt.expectedID, id)
			}
		})
	}
}
