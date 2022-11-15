package backoff

import (
	"context"
)

// The operation to perform in the backoff
type Operation func() error

// Callback when we have an error
type OnErr func(err error)

// The backoff function
type Backoff func(ctx context.Context, operation Operation, onErr OnErr) error

// Returns true if the operation succeeded
type Success func(err error) bool

// The default success condition
func DefaultSuccess(err error) bool {
	return err == nil
}
