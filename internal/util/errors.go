// Package util provides common utility functions and error values.
//
//nolint:revive // util is a common package name for utilities
package util

import "errors"

// ErrNotConfigured indicates that required configuration is missing.
var ErrNotConfigured = errors.New("missing configuration")

// ErrUnexpectedStatus indicates an unexpected HTTP status code was received.
var ErrUnexpectedStatus = errors.New("unexpected status code")
