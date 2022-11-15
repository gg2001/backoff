package backoff

import (
	"context"
	"errors"
	"math"
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	ctx := context.Background()

	min := 10 * time.Millisecond
	max := 100 * time.Millisecond
	factor := float64(2)
	b := ExponentialBackoff(min, max, factor)

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

	// Success after 2 attempts
	counter = 0
	start := time.Now()
	elapsed := time.Duration(0)

	err = b(ctx, func() error {
		counter++

		if counter > 1 {
			elapsed = time.Since(start)
			duration := time.Duration(float64(min) * math.Pow(factor, float64(counter-2)))
			if elapsed < duration {
				t.Fatal("unexpected duration elapsed:", elapsed)
			}
			start = time.Now()
		}

		if counter == 2 {
			return nil
		}
		return errors.New("test")
	}, nil)
	if err != nil {
		t.Fatal("unexpected err")
	}
	if counter != 2 {
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
			duration := time.Duration(float64(min) * math.Pow(factor, float64(counter-2)))
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
	if counter != 4 {
		t.Fatal("invalid counter:", counter)
	}
}
