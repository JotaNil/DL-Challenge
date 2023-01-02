package common

import (
	"testing"
)

type TestCase struct {
	Scenario string
	TestFn   func(t *testing.T)
}
