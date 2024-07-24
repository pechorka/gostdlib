package ptr

func To[T any](val T) *T {
	return &val
}

func NilOnZero[T comparable](val T) *T {
	var zero T
	if zero == val {
		return nil
	}

	return &val
}

func ZeroOnNil[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}

	return *ptr
}
