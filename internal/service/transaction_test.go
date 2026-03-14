package service

import (
	"ledgerA/internal/dto"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToRepoFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    dto.TransactionFilters
		wantErr  bool
		wantType string
	}{
		{"basic", dto.TransactionFilters{Type: "all", Page: 1, PerPage: 20}, false, "all"},
		{"income", dto.TransactionFilters{Type: "income"}, false, "income"},
		{"bad from", dto.TransactionFilters{DateFrom: ptr("bad")}, true, ""},
		{"bad to", dto.TransactionFilters{DateTo: ptr("bad")}, true, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := toRepoFilter(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantType, f.Type)
		})
	}
}

func ptr(v string) *string { return &v }
