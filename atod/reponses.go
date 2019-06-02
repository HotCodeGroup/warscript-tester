package atod

import (
	"encoding/json"
)

type obstacleResp struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

type unitResp struct {
	X           float64 `json:"x"`
	Y           float64 `json:"y"`
	Radius      float64 `json:"radius"`
	Health      float64 `json:"health"`
	ViewRange   float64 `json:"view_range"`
	ReloadLeft  float64 `json:"reload_left"`
	ReloadTime  float64 `json:"reload_time"`
	SpecialLeft float64 `json:"special_left"`
	SpecialTime float64 `json:"special_time"`
	UnitType    string  `json:"unit_type"`
}

type flagResp struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type dropzoneResp struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

type projectileResp struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	VX   float64 `json:"vX"`
	VY   float64 `json:"vY"`
}

// State - implements interface games.State
type State struct {
	Obstacles   []obstacleResp   `json:"obstacles"`
	Projectiles []projectileResp `json:"projectiles"`
	P1Units     []unitResp       `json:"p1_units"`
	P2Units     []unitResp       `json:"p2_units"`
	P1Flags     []flagResp       `json:"p1_flags"`
	P2Flags     []flagResp       `json:"p2_flags"`
}

// JSON - returns marshaled json
func (st *State) JSON() []byte {
	j, _ := json.Marshal(st)
	return j
}

// Info - implements interface games.Info
type Info struct {
	Player1Dropzone dropzoneResp `json:"p1_dropzone"`
	Player2Dropzone dropzoneResp `json:"p2_dropzone"`
	Ratio           float64      `json:"ratio"`
}

// JSON - returns marshaled json
func (i *Info) JSON() []byte {
	j, _ := json.Marshal(i)
	return j
}

// Result - implements interface games.Result
type Result struct {
	Obstacles   []obstacleResp   `json:"obstacles"`
	Projectiles []projectileResp `json:"projectiles"`
	P1Units     []unitResp       `json:"p1_units"`
	P2Units     []unitResp       `json:"p2_units"`
	P1Flags     []flagResp       `json:"p1_flags"`
	P2Flags     []flagResp       `json:"p2_flags"`

	Err1    string `json:"error_1"`
	Err2    string `json:"error_2"`
	Winner  int    `json:"winner"`
	Message string `json:"msg,omitempty"`
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
