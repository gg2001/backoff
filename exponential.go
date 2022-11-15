package backoff

import (
	"context"
	"math"
	"time"
)

// Performs an operation with an exponential backoff
func ExponentialBackoff(min time.Duration, max time.Duration, factor float64) Backoff {
	return ExponentialBackoffWithSuccess(min, max, factor, DefaultSuccess)
}

// Performs an operation with an exponential backoff using a custom success condition
func ExponentialBackoffWithSuccess(min time.Duration, max time.Duration, factor float64, success Success) Backoff {
	// Adjust the params
	if min > max {
		min = max
	}
	minf := float64(min)

	// Return the backoff function
	return func(ctx context.Context, operation Operation, onErr OnErr) error {
		// Set the exit condition
		if success == nil {
			success = DefaultSuccess
		}

		// The operation attempts
		var (
			err      error
			attempts int
			duration time.Duration = min
		)

		// Exit once the duration exceeds the max duration
		for duration <= max {
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

				// Increment the number of attempts
				attempts++

				// Pause in between attempts
				if duration > 0 {
					select {
					case <-time.After(duration):
					case <-ctx.Done():
						return err
					}
				}

				// Update the duration
				duration = time.Duration(minf * math.Pow(factor, float64(attempts)))
			}
		}

		return err
	}
}
