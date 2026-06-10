package util

// Asserts that err is nil and returns the result.
// If err is not nil, a panic occurs.
func Must[T any](result T, err error) T {
	if err != nil {
		panic("must failed: " + err.Error())
	}

	return result
}
