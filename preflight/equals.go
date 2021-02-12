package preflight

import "reflect"

// equals returns true if the arguments are equal
func equals(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Slice:
		return listsEqual(a, b)
	case reflect.Array:
		return listsEqual(a, b)
	default:
		return a == b
	}
}

// listsEqual returns true if every item in two slices are equal
func listsEqual(a, b interface{}) bool {
	x := reflect.ValueOf(a)
	y := reflect.ValueOf(b)

	if x.Len() != y.Len() {
		return false
	}

	for i := 0; i < x.Len(); i++ {
		if x.Index(i).Interface() != y.Index(i).Interface() {
			return false
		}
	}
	return true
}
