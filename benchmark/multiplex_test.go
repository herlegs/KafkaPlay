package benchmark

import "testing"

func TestWithoutMultiplex(t *testing.T) {
	WithoutMultiplex()
}

func TestWithMultiplex(t *testing.T) {
	WithMultiplex()
}