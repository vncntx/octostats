package preflight

// Transformation is a function that modifies the thing under test
type Transformation func(interface{}) interface{}

// Expectation represents an expectation about a value or behavior
type Expectation interface {
	// Sugar
	To() Expectation
	Be() Expectation
	Is() Expectation
	Should() Expectation

	// Negation
	Not() Expectation
	IsNot() Expectation

	// Assertions
	Nil()
	True()
	False()
	Empty()
	HasLength(expected int)
	HaveLength(expected int)
	Equal(expected interface{})
	Equals(expected interface{})
	EqualTo(expected interface{})
	Match(pattern string)
	Matches(pattern string)
}
