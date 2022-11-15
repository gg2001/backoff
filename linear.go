package backoff

import (
	"context"
	"time"
)

// Performs an operation with a linear backoff
func LinearBackoff(attempts int, duration time.Duration) Backoff {
	return LinearBackoffWithSuccess(attempts, duration, DefaultSuccess)
}

// Performs an operation with a linear backoff using a custom success condition
func LinearBackoffWithSuccess(attempts int, duration time.Duration, success Success) Backoff {
	return func(ctx context.Context, operation Operation, onErr OnErr) error {
		// Set the exit condition
		a := attempts
		if success == nil {
			success = DefaultSuccess
		}

		// The error from the operation
		var err error

		// Exit if we reach 0 attempts
		for a != 0 {
			// Make a requst
			err = operation()

			if success(err) {
				// Return if the operation succeeded
				return err
			} else {
				// Call the error callback
				if onErr != nil {
					onErr(err)
				}

				// Decrement the number of attempts
				a--

				// Pause in between attempts
				if duration > 0 {
					select {
					case <-time.After(duration):
					case <-ctx.Done():
						return err
					}
				}
			}
		}

		return err
	}
}
