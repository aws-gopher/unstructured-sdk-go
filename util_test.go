package unstructured

import (
	"fmt"
	"slices"
)

func eq[T comparable](name string, got, want T) error {
	if want != got {
		return fmt.Errorf("expected %s to be %v, got %v", name, want, got)
	}

	return nil
}

func equal[T interface{ Equal(T) bool }](name string, got, want T) error {
	if !want.Equal(got) {
		return fmt.Errorf("expected %s to be %v, got %v", name, want, got)
	}

	return nil
}

func eqs[T comparable](name string, got, want []T) error {
	if !slices.Equal(got, want) {
		return fmt.Errorf("expected %s to be %v, got %v", name, want, got)
	}

	return nil
}
