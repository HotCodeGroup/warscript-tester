package games

// State - state of the game
type State interface {
	JSON() []byte
}

// Result - result of the game
type Result interface {
	State
	GetWinner() int
}

// Game interface for working with different games
type Game interface {
	Init()
	Images() (string, string)
	Snapshots() (shot1, shot2 []byte)
	SaveSnapshots(shot1, shot2 []byte) error
	GetState() (state State, fin bool)
	GetResult() (result Result)
}

type GameError struct {
	msg string
}

func (e *GameError) Error() string {
	return e.msg
}

var (
	ErrPlayer1Fail = &GameError{
		msg: "player1 response was incorrect",
	}
	ErrPlayer2Fail = &GameError{
		msg: "player2 response was incorrect",
	}
	ErrInternal = &GameError{
		msg: "internal failure",
	}
)
