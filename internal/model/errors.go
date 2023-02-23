package model

type RetryError struct {
	Err        error
	RetryAfter int
}

func (r RetryError) Error() string {
	return r.Err.Error()
}
