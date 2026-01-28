package config

import "context"

type ctxKey uint8

const configKey ctxKey = iota

// NewContext returns a new context with the given config attached.
func NewContext(ctx context.Context, conf *Config) context.Context {
	return context.WithValue(ctx, configKey, conf)
}

// FromContext retrieves the config from the context.
func FromContext(ctx context.Context) (*Config, bool) {
	conf, ok := ctx.Value(configKey).(*Config)
	return conf, ok
}
