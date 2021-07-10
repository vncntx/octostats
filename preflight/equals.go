package preflight

import "reflect"

// equals returns true if the arguments are equal
func equals(a, b interface{}) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	switch reflect.TypeOf(a).Kind() {
	case reflect.Invalid:
		return false
	case reflect.Slice:
		return listsEqual(a, b)
	case reflect.Array:
		return listsEqual(a, b)
	case reflect.Bool, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.Float32, reflect.Float64, reflect.Func, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Interface, reflect.Map, reflect.Ptr, reflect.String, reflect.Struct, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.UnsafePointer:
		return a == b
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
