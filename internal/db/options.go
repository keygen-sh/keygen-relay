package db

type QueryOption func(*queryOptions)

type queryOptions struct {
	pool        *Pool
	withoutPool bool
}

func WithPool(pool *Pool) QueryOption {
	return func(opts *queryOptions) {
		opts.pool = pool
		opts.withoutPool = false
	}
}

func WithoutPool() QueryOption {
	return func(opts *queryOptions) {
		opts.pool = nil
		opts.withoutPool = true
	}
}

func applyOptions(opts ...QueryOption) *queryOptions {
	options := &queryOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}