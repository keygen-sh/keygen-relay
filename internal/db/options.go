package db

// AnyPool is a sentinel value indicating queries should match licenses regardless of pool
var AnyPool = &Pool{ID: 0, Name: "<any>"}

// LicenseOption is a functional option for license queries
type LicenseOption func(*licenseOptions)

type licenseOptions struct {
	// pool can be:
	// - nil: query licenses with null pool_id
	// - *Pool: query licenses with specific pool_id
	// - AnyPool: query all licenses regardless of pool_id
	pool *Pool
}

// WithPool queries licenses for a specific pool
func WithPool(pool *Pool) LicenseOption {
	return func(opts *licenseOptions) {
		opts.pool = pool
	}
}

// WithoutPool queries licenses without a pool
func WithoutPool() LicenseOption {
	return func(opts *licenseOptions) {
		opts.pool = nil
	}
}

// WithAnyPool queries all licenses regardless of pool
func WithAnyPool() LicenseOption {
	return func(opts *licenseOptions) {
		opts.pool = AnyPool
	}
}

func applyLicenseOptions(opts ...LicenseOption) *licenseOptions {
	options := &licenseOptions{
		pool: AnyPool, // default to all licenses
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// isAnyPool checks if the pool is the AnyPool sentinel
func isAnyPool(pool *Pool) bool {
	return pool != nil && pool.ID == 0 && pool.Name == "<any>"
}
