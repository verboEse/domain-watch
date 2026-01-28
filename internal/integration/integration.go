package integration

import (
	"context"

	"gabe565.com/domain-watch/internal/config"
)

// Integration defines the interface for notification integrations.
type Integration interface {
	Name() string
	Setup(ctx context.Context, conf *config.Config) error
	Send(ctx context.Context, text string) error
}
