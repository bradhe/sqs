package sqs

import (
	"testing"
)

func TestChunksAppropriately(t *testing.T) {
	arr := []string{"one", "two", "three", "four", "five"}
	res := chunk(arr, 2)

	if len(res) != 3 {
		t.Fatalf("Expected len(res) = 3, got %d", len(res))
	}
}
