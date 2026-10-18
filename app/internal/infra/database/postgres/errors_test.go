package postgres

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/marcoscaco/pismo-transactions-routine-service/internal/domain"
)

func TestMapPostgresError(t *testing.T) {
	tests := []struct {
		name     string
		postgres error
		want     error
	}{
		{
			name:     "unknown operation type foreign key",
			postgres: &pq.Error{Code: "23503", Constraint: "fk_operation_type"},
			want:     domain.ErrInvalidOperationTypeID,
		},
		{
			name:     "unknown account foreign key",
			postgres: &pq.Error{Code: "23503", Constraint: "fk_account"},
			want:     domain.ErrAccountNotFound,
		},
		{
			name:     "unrelated database error",
			postgres: &pq.Error{Code: "23505", Constraint: "some_unique_constraint"},
		},
		{
			name:     "unmapped foreign key",
			postgres: &pq.Error{Code: "23503", Constraint: "other_foreign_key"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapPostgresError(fmt.Errorf("insert transaction: %w", tt.postgres))
			if tt.want != nil {
				if !errors.Is(got, tt.want) {
					t.Fatalf("mapPostgresError() = %v, want %v", got, tt.want)
				}
				return
			}
			if got == nil {
				t.Fatal("mapPostgresError() = nil, want original error")
			}
			var gotPostgresError *pq.Error
			if !errors.As(got, &gotPostgresError) || gotPostgresError != tt.postgres {
				t.Fatalf("mapPostgresError() = %v, want original PostgreSQL error", got)
			}
		})
	}
}
