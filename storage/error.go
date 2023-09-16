package storage

type Error uint

const (
	ErrNotFound Error = iota + 1
)

func (e Error) Error() string {
	switch e {
	case ErrNotFound:
		return "not found"
	default:
		return "unknown error"
	}
}
