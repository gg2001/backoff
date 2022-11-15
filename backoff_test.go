package backoff

import (
	"errors"
	"testing"
)

func TestDefaultSuccess(t *testing.T) {
	if !DefaultSuccess(nil) {
		t.Fatal("unexpected err")
	}

	if DefaultSuccess(errors.New("test")) {
		t.Fatal("unexpected success")
	}
}
