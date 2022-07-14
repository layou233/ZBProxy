package common

func Must[T any](s T, err error) T {
	if err != nil {
		panic(err)
	}
	return s
}

func GetSecond[T any](_ any, r T) T {
	return r
}
