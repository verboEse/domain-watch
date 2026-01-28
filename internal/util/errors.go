// Package util provides common utility functions and error values.
package util

import "errors"

// ErrNotConfigured indicates that required configuration is missing.
var ErrNotConfigured = errors.New("missing configuration")

// ErrUnexpectedStatus indicates an unexpected HTTP status code was received.
var ErrUnexpectedStatus = errors.New("unexpected status code")
