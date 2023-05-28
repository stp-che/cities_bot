package common

type DomainError struct {
	msg string
}

func NewDomainError(msg string) *DomainError {
	return &DomainError{msg: msg}
}

func (e DomainError) Error() string {
	return e.msg
}
