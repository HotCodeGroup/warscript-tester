package pong

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/games"
)

type object2D struct {
	x float64 `json:"x"`
	y float64 `json:"y"`
}

// State - implements interface games.State
type State struct {
	Player1 object2D `json:"player1"`
	Player2 object2D `json:"player2"`
	Ball    object2D `json:"ball"`
}

// JSON - returns marshaled json
func (st *State) JSON() []byte {
	json, _ := json.Marshal(st)
	return json
}

// Result - implements interface games.Result
type Result struct {
	Player1     object2D `json:"player1"`
	Player2     object2D `json:"player2"`
	Ball        object2D `json:"ball"`
	Winner      int      `json:"winner"`
	Player1Fail *games.GameError
	Player2Fail *games.GameError
	InternalErr *games.GameError
	Message     string `json:"msg,omitempty"`
}

// GetWinner - returns winner of the game
func (res *Result) GetWinner() int {
	return res.Winner
}

// JSON - returns marshaled json
func (res *Result) JSON() []byte {
	json, _ := json.Marshal(res)
	return json
}
