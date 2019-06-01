package pong

import (
	"encoding/json"
)

type object2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// State - implements interface games.State
type State struct {
	Player1 object2D `json:"player_1"`
	Player2 object2D `json:"player_2"`
	Ball    object2D `json:"ball"`
	Console []string `json:"console"`
}

// JSON - returns marshaled json
func (st *State) JSON() []byte {
	j, _ := json.Marshal(st)
	return j
}

// Info - implements interface games.Info
type Info struct {
	Ball struct {
		Diameter float64 `json:"diameter"`
	} `json:"ball"`

	Racket struct {
		Height float64 `json:"h"`
		Width  float64 `json:"w"`
	} `json:"racket"`

	Ratio float64 `json:"ratio"`
}

// JSON - returns marshaled json
func (i *Info) JSON() []byte {
	j, _ := json.Marshal(i)
	return j
}

// Result - implements interface games.Result
type Result struct {
	Player1 object2D `json:"player_1"`
	Player2 object2D `json:"player_2"`
	Ball    object2D `json:"ball"`
	Winner  int      `json:"winner"`
	Err1    string   `json:"error_1"`
	Err2    string   `json:"error_2"`
	Message string   `json:"msg,omitempty"`
}

// GetWinner - returns winner of the game
func (res *Result) GetWinner() int {
	return res.Winner
}

func (res *Result) Error1() string {
	return res.Err1
}

func (res *Result) Error2() string {
	return res.Err2
}

// JSON - returns marshaled json
func (res *Result) JSON() []byte {
	j, _ := json.Marshal(res)
	return j
}
