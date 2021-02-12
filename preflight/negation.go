package preflight

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

// Negation is a negated expectation
type Negation struct {
	*testing.T

	Actual  interface{}
	Inverse Expectation
}

// To returns the current expectation
func (not *Negation) To() Expectation {
	return not
}

// Be returns the current expectation
func (not *Negation) Be() Expectation {
	return not
}

// Is returns the current expectation
func (not *Negation) Is() Expectation {
	return not
}

// Should returns the current expectation
func (not *Negation) Should() Expectation {
	return not
}

// Not returns a negation of the current expectation
func (not *Negation) Not() Expectation {
	return not.Inverse
}

// IsNot is equivalent to Is().Not()
func (not *Negation) IsNot() Expectation {
	return not.Inverse.Is().Not()
}

// Nil asserts the value is not nil
func (not *Negation) Nil() {
	not.Equals(nil)
}

// True asserts the value is true
func (not *Negation) True() {
	not.Equals(true)
}

// False asserts the value is false
func (not *Negation) False() {
	not.Equals(false)
}

// Empty asserts the value is an array with length != 0
func (not *Negation) Empty() {
	not.HasLength(0)
}

// HasLength asserts the value is an array with length != given
func (not *Negation) HasLength(given int) {
	if reflect.ValueOf(not.Actual).Len() == given {
		not.Fail()
	}
}

// HaveLength is equivalent to HasLength()
func (not *Negation) HaveLength(given int) {
	not.HasLength(given)
}

// Equals asserts inequality to a given value
func (not *Negation) Equals(given interface{}) {
	if not.Actual == given {
		not.Fail()
	}
}

// Equal is equivalent to Equals()
func (not *Negation) Equal(given interface{}) {
	not.Equals(given)
}

// EqualTo is equivalent to Equals()
func (not *Negation) EqualTo(given interface{}) {
	not.Equals(given)
}

// Matches asserts the value does not match a pattern
func (not *Negation) Matches(pattern string) {
	actual := fmt.Sprint(not.Actual)
	isMatch, err := regexp.MatchString(pattern, actual)
	if err != nil {
		not.Error("failed to compile regular expression")
	} else if isMatch {
		not.Fail()
	}
}

// Match is equivalent to Matches()
func (not *Negation) Match(pattern string) {
	not.Matches(pattern)
}
