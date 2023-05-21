package telegram

type UserError struct {
	Msg string
}

func (e UserError) Error() string {
	return e.Msg
}

var (
	ErrGameNameNotSpecified = UserError{"game name not specified"}
	ErrGameDoesNotExist     = UserError{"game does not exist"}
	ErrGameAlreadyStarted   = UserError{"game already started"}
	ErrGameNotStarted       = UserError{"game not started"}
)
