package errors

type notFoundError struct{}

func (*notFoundError) Error() string { return "Not Found" }

var NotFoundError error = &notFoundError{}

type illegalEntityError struct{ message string }

func (c *illegalEntityError) Error() string { return c.message }

func IllegalEntityError(message string) error { return &illegalEntityError{message} }
