package utils

type Error struct {
	Msg string
}

func (e Error) Error() string {
	return e.Msg
}

func NewError(msg string) error {
	return &Error{msg}
}
