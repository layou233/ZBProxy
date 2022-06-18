package buf

type hasInnerError interface {
	// Unwrap returns the underlying error of this one.
	Unwrap() error
}

// cause returns the root cause of this error.
func cause(err error) error {
	if err == nil {
		return nil
	}
L:
	for {
		switch inner := err.(type) {
		case hasInnerError:
			if inner.Unwrap() == nil {
				break L
			}
			err = inner.Unwrap()
		default:
			break L
		}
	}
	return err
}
