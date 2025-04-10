package types_test

import (
	"fmt"
	"testing"

	"github.com/ArtificialLegacy/imgscal/pkg/types"
)

func Test_ParseEscaped(t *testing.T) {
	teststr := "test, {escaped_1}, test, {escaped_2}"
	resultstr := "test, %s, test, %s"
	finalstr := "test, escaped_1, test, escaped_2"

	result, escaped := types.ParseEscaped(teststr)

	if len(escaped) != 2 {
		t.Fatalf("expected 2 results, got %d", len(escaped))
	}

	if result != resultstr {
		t.Fatalf("wrong result: expected %s, got %s", resultstr, result)
	}

	refmt := fmt.Sprintf(result, anyArray(escaped)...)
	if refmt != finalstr {
		t.Fatalf("strings mismatched: expected %s, got %s", teststr, refmt)
	}
}

func anyArray[T any](a []T) []any {
	result := make([]any, len(a))

	for k, v := range a {
		result[k] = v
	}

	return result
}
