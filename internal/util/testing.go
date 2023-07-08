package util

import "fmt"

// TestingT interface to use only Fatalf from testing.T
type TestingT interface {
	Fatalf(format string, args ...any)
}

// MockTestingT mock for testing.T
type MockTestingT struct {
	fatalCalled bool
}

// Fatalf mock function of testing.T.Fatalf
func (t *MockTestingT) Fatalf(format string, args ...any) {
	t.fatalCalled = true
	panic(fmt.Sprintf(format, args...))
}
