package db

import (
	"testing"
)

func TestAnyPool(t *testing.T) {
	// AnyPool should have ID=0 and Name="<any>"
	if AnyPool.ID != 0 {
		t.Errorf("expected AnyPool.ID to be 0, got %d", AnyPool.ID)
	}
	if AnyPool.Name != "<any>" {
		t.Errorf("expected AnyPool.Name to be \"<any>\", got %q", AnyPool.Name)
	}
}

func TestIsAnyPool(t *testing.T) {
	tests := []struct {
		name     string
		pool     *Pool
		expected bool
	}{
		{
			name:     "nil pool",
			pool:     nil,
			expected: false,
		},
		{
			name:     "AnyPool",
			pool:     AnyPool,
			expected: true,
		},
		{
			name:     "regular pool",
			pool:     &Pool{ID: 1, Name: "prod"},
			expected: false,
		},
		{
			name:     "pool with same ID but different name",
			pool:     &Pool{ID: 0, Name: "prod"},
			expected: false,
		},
		{
			name:     "pool with different ID but same name",
			pool:     &Pool{ID: 1, Name: "<any>"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if actual := tt.pool == AnyPool; actual != tt.expected {
				t.Errorf("%v == AnyPool = %v, expected %v", tt.pool, actual, tt.expected)
			}
		})
	}
}

func TestLicensePredicates(t *testing.T) {
	t.Run("default predicates", func(t *testing.T) {
		opts := applyLicensePredicates()
		if opts.pool != AnyPool {
			t.Error("expected default predicates to use AnyPool")
		}
	})

	t.Run("WithPool", func(t *testing.T) {
		pool := &Pool{ID: 1, Name: "test"}

		opts := applyLicensePredicates(WithPool(pool))
		if opts.pool != pool {
			t.Errorf("expected pool to be %v, got %v", pool, opts.pool)
		}

		opts = applyLicensePredicates(WithPool(nil))
		if opts.pool != nil {
			t.Errorf("expected pool to be %v, got %v", nil, opts.pool)
		}
	})

	t.Run("WithoutPool", func(t *testing.T) {
		opts := applyLicensePredicates(WithoutPool())
		if opts.pool != nil {
			t.Errorf("expected pool to be nil, got %v", opts.pool)
		}
	})

	t.Run("WithAnyPool", func(t *testing.T) {
		opts := applyLicensePredicates(WithAnyPool())
		if opts.pool != AnyPool {
			t.Error("expected pool to be AnyPool")
		}
	})

	t.Run("last option wins", func(t *testing.T) {
		pool := &Pool{ID: 1, Name: "test"}

		opts := applyLicensePredicates(WithPool(pool), WithoutPool())
		if opts.pool != nil {
			t.Errorf("expected pool to be nil (last option), got %v", opts.pool)
		}
	})
}
