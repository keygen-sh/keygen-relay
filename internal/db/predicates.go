package db

// AnyPool is a sentinel value indicating queries should match licenses regardless of pool
var AnyPool = &Pool{ID: 0, Name: "<any>"}

// LicensePredicateFunc is a functional predicate for license queries
type LicensePredicateFunc func(*LicensePredicate)

type LicensePredicate struct {
	// pool can be:
	// - nil: query licenses with null pool_id
	// - *Pool: query licenses with specific pool_id
	// - AnyPool: query all licenses regardless of pool_id
	pool *Pool
}

// WithPool queries licenses for a specific pool
func WithPool(pool *Pool) LicensePredicateFunc {
	return func(pred *LicensePredicate) {
		pred.pool = pool
	}
}

// WithoutPool queries licenses without a pool
func WithoutPool() LicensePredicateFunc {
	return func(predicate *LicensePredicate) {
		predicate.pool = nil
	}
}

// WithAnyPool queries all licenses regardless of pool
func WithAnyPool() LicensePredicateFunc {
	return func(predicate *LicensePredicate) {
		predicate.pool = AnyPool
	}
}

func applyLicensePredicates(fns ...LicensePredicateFunc) *LicensePredicate {
	predicates := &LicensePredicate{
		pool: AnyPool, // default to all licenses
	}

	for _, fn := range fns {
		fn(predicates)
	}

	return predicates
}
