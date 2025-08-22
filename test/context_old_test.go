//go:build !go1.24 && integration

package test

import (
	"context"
	"crypto/rand"
	"testing"
)

// textContext mimics [*testing.T.Context], which was added in go1.24.
func testContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	return ctx
}

func randText() string {
	const base32 = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	src := make([]byte, 26)
	rand.Read(src)
	for i := range src {
		src[i] = base32[src[i]%32]
	}
	return string(src)
}
