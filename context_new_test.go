//go:build go1.24

package unstructured

import (
	"context"
	"testing"
)

// textContext mimics [*testing.T.Context], which was added in go1.24.
func testContext(t *testing.T) context.Context {
	return t.Context()
}
