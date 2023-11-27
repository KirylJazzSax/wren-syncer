package utils

func Must[T any](a T, err error) T {
	if err != nil {
		panic(err)
	}
	return a
}

func MustVoid(err error) {
	if err != nil {
		panic(err)
	}
}
