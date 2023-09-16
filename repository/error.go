package repository

type Error uint

const (
	ErrBookNotFound Error = iota + 1
)

func (e Error) Error() string {
	switch e {
	case ErrBookNotFound:
		return "book not found"
	default:
		return "unknown error"
	}
}
