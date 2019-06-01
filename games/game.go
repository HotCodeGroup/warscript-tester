package games

// State - state of the game
type State interface {
	JSON() []byte
}

// Info represents information about game objects
type Info interface {
	JSON() []byte
}

// Result - result of the game
type Result interface {
	State
	GetWinner() int
	Error1() string
	Error2() string
}

// Game interface for working with different games
type Game interface {
	Init()
	Images() (string, string)
	Snapshots() (shot1, shot2 []byte)
	SaveSnapshots(shot1, shot2 []byte) error
	GetInfo() (info Info)
	GetState() (state State, fin bool)
	GetResult() (result Result)
	GetLogs() ([]string, []string)
}

// GameError ошибка, возникшая при проверки игры
type GameError struct {
	Msg string `json:"message"`
}

func (e *GameError) Error() string {
	return e.Msg
}

var (
	// ErrPlayer1Fail ошибка в ответе игрока 1
	ErrPlayer1Fail = &GameError{
		Msg: "player1 response was incorrect",
	}
	// ErrPlayer2Fail ошибка в ответе игрока 2
	ErrPlayer2Fail = &GameError{
		Msg: "player2 response was incorrect",
	}
	// ErrInternal внутренняя ошибка сети
	ErrInternal = &GameError{
		Msg: "internal failure",
	}
)
