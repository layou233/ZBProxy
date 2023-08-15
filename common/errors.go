package common

type causeError struct {
	inner error
	cause string
}

func (e causeError) Error() string {
	return e.cause + e.inner.Error()
}

func (e causeError) Unwrap() error {
	return e.inner
}

func Cause(cause string, err error) error {
	return causeError{
		inner: err,
		cause: cause,
	}
}

type HasInnerError interface {
	Unwrap() error
}

func Unwrap(err error) error {
	for {
		inner, ok := err.(HasInnerError)
		if !ok {
			break
		}
		innerErr := inner.Unwrap()
		if innerErr == nil {
			break
		}
		err = innerErr
	}
	return err
}
