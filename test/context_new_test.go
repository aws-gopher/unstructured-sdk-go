//go:build go1.24 && integration

package test

import (
	"context"
	"crypto/rand"
	"testing"
)

// textContext mimics [*testing.T.Context], which was added in go1.24.
func testContext(t *testing.T) context.Context {
	return t.Context()
}

func randText() string {
	return rand.Text()
}
