package netproto

type Error struct {
	Message string
}

var _ error = (*Error)(nil)

func NewError(message string) *Error {
	return &Error{
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Is(target error) bool {
	return e.Message == target.Error()
}
