package errors

import "reflect"

// AsGeneric is an experimental generic equivalent of As.
//
// It allows using a shorter and safer syntax:
//
//     // With As
//     var httpError *HTTPError
//     if errors.As(err, &httpError) {
//         fmt.Println(httpError.code)
//     }
//
//     // With AsGeneric
//     if httpError, ok := errors.AsGeneric[*HTTPError](err); ok {
//         fmt.Println(httpError.code)
//     }
func AsGeneric[T any](err error) (T, bool) {
	if err == nil {
		var zero T
		return zero, false
	}

	target := new(T)
	targetVal := reflect.ValueOf(target).Elem()
	targetType := targetVal.Type()

	for err != nil {
		if reflect.TypeOf(err).AssignableTo(targetType) {
			targetVal.Set(reflect.ValueOf(err))
			return *target, true
		}

		if x, ok := err.(interface{ As(any) bool }); ok && x.As(target) {
			return *target, true
		}

		err = Unwrap(err)
	}

	return *target, false
}
