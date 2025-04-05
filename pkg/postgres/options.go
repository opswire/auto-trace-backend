package postgres

type Option func(postgres *Postgres)

func WithMaxPoolSize(size int) Option {
	return func(postgres *Postgres) {
		postgres.maxPoolSize = size
	}
}

func WithConnAttempts(attempts int) Option {
	return func(postgres *Postgres) {
		postgres.connAttempts = attempts
	}
}

func WithConnTimeout(attempts int) Option {
	return func(postgres *Postgres) {
		postgres.connAttempts = attempts
	}
}
