package backoff

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLinearBackoff(t *testing.T) {
	ctx := context.Background()

	attempts := 10
	duration := 10 * time.Millisecond
	b := LinearBackoff(attempts, duration)

	// Success
	counter := 0

	err := b(ctx, func() error {
		counter++

		return nil
	}, nil)
	if err != nil {
		t.Fatal("unexpected err")
	}
	if counter != 1 {
		t.Fatal("invalid counter:", counter)
	}

	// Success after 5 attempts
	counter = 0
	start := time.Now()
	elapsed := time.Duration(0)

	err = b(ctx, func() error {
		counter++

		if counter > 1 {
			elapsed = time.Since(start)
			if elapsed < duration {
				t.Fatal("unexpected duration elapsed:", elapsed)
			}
			start = time.Now()
		}

		if counter == 5 {
			return nil
		}
		return errors.New("test")
	}, nil)
	if err != nil {
		t.Fatal("unexpected err")
	}
	if counter != 5 {
		t.Fatal("invalid counter:", counter)
	}

	// Fail
	counter = 0
	start = time.Now()
	elapsed = time.Duration(0)

	err = b(ctx, func() error {
		counter++

		if counter > 1 {
			elapsed = time.Since(start)
			if elapsed < duration {
				t.Fatal("unexpected duration elapsed:", elapsed)
			}
			start = time.Now()
		}

		return errors.New("test")
	}, nil)
	if err == nil {
		t.Fatal("unexpected nil")
	}
	if counter != 10 {
		t.Fatal("invalid counter:", counter)
	}
}