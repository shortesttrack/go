package errors

import (
	"testing"
	serrors "errors"
)

func TestEqual(t *testing.T) {
	e1 := NotFound
	e2 := New(NotFound.Error())
	e3 := serrors.New(NotFound.Error())

	if !Equal(e1, e2) {
		t.Error(`errors are equal`)
	}

	if !Equal(e1, e3) {
		t.Error(`errors are equal`)
	}
}
