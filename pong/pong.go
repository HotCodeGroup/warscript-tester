package pong

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/games"
)

// Movable - 2D object with position, speed, height and width
type Movable struct {
	x      float64
	y      float64
	vX     float64
	vY     float64
	height float64
	width  float64
}

// Pong - pong game
type Pong struct {
	height  float64
	width   float64
	ball    Movable
	player1 Movable
	player2 Movable

	winner  int
	isEnded bool
}

const (
	fieldHeight  = 250
	fieldWidth   = 500
	playerHeight = 50
	playerWidth  = 25
	ballRad      = 10
	image        = "pong"
)

// Init - inits game: sets default params
func (pong *Pong) Init() {
	pong.ball = Movable{
		height: ballRad,
		width:  ballRad,
		x:      fieldWidth / 2,
		y:      fieldHeight / 2,
		vX:     2,
		vY:     2,
	}

	pong.player1 = Movable{
		height: fieldHeight / 5,
		width:  fieldWidth / 20,
		x:      fieldWidth / 10,
		y:      fieldHeight / 2,
		vX:     0,
		vY:     0,
	}

	pong.player2 = Movable{
		height: fieldHeight / 5,
		width:  fieldWidth / 20,
		x:      fieldWidth - fieldWidth/10,
		y:      fieldHeight / 2,
		vX:     0,
		vY:     0,
	}
}

// Images - returns names of images
func (pong *Pong) Images() (string, string) {
	return image, image
}

type shotInner struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	VX float64 `json:"vX"`
	VY float64 `json:"vY"`
}

type shot struct {
	Me    shotInner `json:"me"`
	Enemy shotInner `json:"enemy"`
	Ball  shotInner `json:"ball"`
}

func (pong *Pong) createShot1() shot {
	return shot{
		Me: shotInner{
			X:  pong.player1.x,
			Y:  pong.player1.y,
			VX: pong.player1.vX,
			VY: pong.player1.vY,
		},
		Enemy: shotInner{
			X:  pong.player2.x,
			Y:  pong.player2.y,
			VX: pong.player2.vX,
			VY: pong.player2.vY,
		},
		Ball: shotInner{
			X:  pong.ball.x,
			Y:  pong.ball.y,
			VX: pong.ball.vX,
			VY: pong.ball.vY,
		},
	}
}

func (pong *Pong) createShot2() shot {
	return shot{
		Me: shotInner{
			X:  pong.width - pong.player2.x,
			Y:  pong.height - pong.player2.y,
			VX: -pong.player2.vX,
			VY: -pong.player2.vY,
		},
		Enemy: shotInner{
			X:  pong.width - pong.player1.x,
			Y:  pong.height - pong.player1.y,
			VX: -pong.player1.vX,
			VY: -pong.player1.vY,
		},
		Ball: shotInner{
			X:  pong.width - pong.ball.x,
			Y:  pong.height - pong.ball.y,
			VX: -pong.ball.vX,
			VY: -pong.ball.vY,
		},
	}
}

// Snapshots - returns encoded json struct
// that can be send to test server
func (pong *Pong) Snapshots() (shot1, shot2 []byte) {
	shot1, _ = json.Marshal(pong.createShot1())
	shot2, _ = json.Marshal(pong.createShot2())
	return
}

func (pong *Pong) SaveSnapshots(shot1, shot2 []byte) (gameErr error) {
	var s1, s2 shot
	err1 := json.Unmarshal(shot1, s1)
	if err1 != nil {
		err1 = games.ErrPlayer1Fail
	}
	err2 := json.Unmarshal(shot2, s2)
	if err2 != nil {
		err2 = games.ErrPlayer2Fail
	}
	if err1 != nil || err2 != nil {

	}
}

func (pong *Pong) GetState() (state State, fin bool) {}
func (pong *Pong) GetResult() (result Result)        {}
