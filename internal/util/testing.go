package util

import "fmt"

type TestingT interface {
	Fatalf(format string, args ...any)
}

type MockTestingT struct {
	fatalCalled bool
}

func (t *MockTestingT) Fatalf(format string, args ...any) {
	t.fatalCalled = true
	panic(fmt.Sprintf(format, args...))
}
