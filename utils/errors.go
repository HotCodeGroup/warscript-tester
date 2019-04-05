package utils

// GameError error ocсured during the game
type GameError struct {
	msg string
}

func (e *GameError) Error() string {
	return e.msg
}
