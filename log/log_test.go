package log

import (
	"errors"
	"testing"
)

func TestDEBUG(t *testing.T) {
	err := errors.New("test error")
	DEBUG("%v", err)

	DEBUG("%v", T{"etst", 10})
}

type T struct {
	name string
	age  int
}
