package errors

import "reflect"

// isNonRecursive is a copy of errors.Is, without Unwrapping err.
func isNonRecursive(err, target error) bool {
	if target == nil {
		return err == target
	}

	isComparable := reflect.TypeOf(target).Comparable()
	if isComparable && err == target {
		return true
	}

	if x, ok := err.(Iser); ok {
		return x.Is(target)
	}

	return false
}

// asNonRecursive is a copy of errors.As, without Unwrapping err.
func asNonRecursive(err error, target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	targetType := typ.Elem()
	if targetType.Kind() != reflect.Interface && !targetType.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}

	if reflect.TypeOf(err).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(err))
		return true
	}
	if x, ok := err.(Aser); ok && x.As(target) {
		return true
	}

	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
