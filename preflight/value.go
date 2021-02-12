package preflight

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

// ValueExpectation is an expectation based on a realized value
type ValueExpectation struct {
	*testing.T

	Actual interface{}
}

// ExpectValue returns a new value-based Expectation
func ExpectValue(t *testing.T, actual interface{}) Expectation {
	return &ValueExpectation{
		T:      t,
		Actual: actual,
	}
}

// To returns the current expectation
func (e *ValueExpectation) To() Expectation {
	return e
}

// Be returns the current expectation
func (e *ValueExpectation) Be() Expectation {
	return e
}

// Is returns the current expectation
func (e *ValueExpectation) Is() Expectation {
	return e
}

// Should returns the current expectation
func (e *ValueExpectation) Should() Expectation {
	return e
}

// Not returns a negation of the current expectation
func (e *ValueExpectation) Not() Expectation {
	return &Negation{
		T:       e.T,
		Actual:  e.Actual,
		Inverse: e,
	}
}

// IsNot is equivalent to Is().Not()
func (e *ValueExpectation) IsNot() Expectation {
	return e.Is().Not()
}

// Nil asserts the value is nil
func (e *ValueExpectation) Nil() {
	e.Equals(nil)
}

// True asserts the value is true
func (e *ValueExpectation) True() {
	e.Equals(true)
}

// False asserts the value is false
func (e *ValueExpectation) False() {
	e.Equals(false)
}

// Empty asserts the value has length 0
func (e *ValueExpectation) Empty() {
	e.HasLength(0)
}

// HasLength asserts the value is an array with a given length
func (e *ValueExpectation) HasLength(expected int) {
	if reflect.ValueOf(e.Actual).Len() != expected {
		e.Errorf("%s: len(%s) != %d", e.Name(), e.Actual, expected)
	}
}

// HaveLength is equivalent to HasLength()
func (e *ValueExpectation) HaveLength(expected int) {
	e.HasLength(expected)
}

// Equals asserts equality to an expected value
func (e *ValueExpectation) Equals(expected interface{}) {
	if !equals(e.Actual, expected) {
		e.Errorf("%s: %#v != %#v", e.Name(), e.Actual, expected)
	}
}

// Equal is equivalent to Equals()
func (e *ValueExpectation) Equal(expected interface{}) {
	e.Equals(expected)
}

// EqualTo is equivalent to Equals()
func (e *ValueExpectation) EqualTo(expected interface{}) {
	e.Equals(expected)
}

// Matches asserts the value matches a pattern
func (e *ValueExpectation) Matches(pattern string) {
	actual := fmt.Sprint(e.Actual)
	isMatch, err := regexp.MatchString(pattern, actual)
	if err != nil {
		e.Error("failed to compile regular expression")
	} else if !isMatch {
		e.Errorf("'%s' does not match /%s/", e.Actual, pattern)
	}
}

// Match is equivalent to Matches()
func (e *ValueExpectation) Match(pattern string) {
	e.Matches(pattern)
}
