package session

type Repository interface {
	GetByChatID(int64) (*Session, error)
	Save(*Session) error
}
