package test

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func CheckDiff(t *testing.T, expected, actual interface{}) {
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("mismatch (-want +got):\n%s", diff)
	}
}
