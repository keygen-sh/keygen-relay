package db

import (
	"testing"
)

func TestLicensePredicates(t *testing.T) {
	t.Run("default predicates", func(t *testing.T) {
		predicates := applyLicensePredicates()
		if predicates.pool != AnyPool {
			t.Error("expected default predicates to use AnyPool")
		}
	})

	t.Run("WithPool", func(t *testing.T) {
		pool := &Pool{ID: 1, Name: "test"}

		predicates := applyLicensePredicates(WithPool(pool))
		if predicates.pool != pool {
			t.Errorf("expected pool to be %v, got %v", pool, predicates.pool)
		}

		predicates = applyLicensePredicates(WithPool(nil))
		if predicates.pool != nil {
			t.Errorf("expected pool to be %v, got %v", nil, predicates.pool)
		}
	})

	t.Run("WithoutPool", func(t *testing.T) {
		predicates := applyLicensePredicates(WithoutPool())
		if predicates.pool != nil {
			t.Errorf("expected pool to be nil, got %v", predicates.pool)
		}
	})

	t.Run("WithAnyPool", func(t *testing.T) {
		predicates := applyLicensePredicates(WithAnyPool())
		if predicates.pool != AnyPool {
			t.Error("expected pool to be AnyPool")
		}
	})

	t.Run("last option wins", func(t *testing.T) {
		pool := &Pool{ID: 1, Name: "test"}

		predicates := applyLicensePredicates(WithPool(pool), WithoutPool())
		if predicates.pool != nil {
			t.Errorf("expected pool to be nil (last option), got %v", predicates.pool)
		}
	})
}
