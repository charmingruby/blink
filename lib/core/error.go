package core

func NewNotFoundError(msg string) error {
	return &NotFoundError{
		message: msg,
	}
}

type NotFoundError struct {
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

func NewUnknowClientError(msg string) error {
	return &UnknowClientError{
		message: msg,
	}
}

type UnknowClientError struct {
	message string
}

func (e *UnknowClientError) Error() string {
	return e.message
}
