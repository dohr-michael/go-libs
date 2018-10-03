package errors

type notFoundError struct{}

func (*notFoundError) Error() string { return "Not Found" }

var NotFoundError error = &notFoundError{}
