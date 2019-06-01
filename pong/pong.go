package pong

import (
	"encoding/json"

	"github.com/HotCodeGroup/warscript-tester/games"

	"github.com/pkg/errors"
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
	ticksLeft int
	height    float64
	width     float64
	ball      Movable
	player1   Movable
	player2   Movable
	winner    int
	isEnded   bool
	Error1    string
	Error2    string
}

const (
	fieldHeight = 250
	fieldWidth  = 500
	ballRad     = 10
	image       = "pong"
)

// Init - inits game: sets default params
func (pong *Pong) Init() {
	pong.ticksLeft = 10000
	pong.width = fieldWidth
	pong.height = fieldHeight
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
func (pong *Pong) Images() (im1, im2 string) {
	return image, image
}

type shotInner struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	VX     float64 `json:"vX"`
	VY     float64 `json:"vY"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

type gameShotInner struct {
	Height    float64 `json:"height"`
	Width     float64 `json:"width"`
	TicksLeft int     `json:"ticks_left"`
}

type shot struct {
	Me      shotInner     `json:"me"`
	Enemy   shotInner     `json:"enemy"`
	Ball    shotInner     `json:"ball"`
	Game    gameShotInner `json:"game"`
	Console []string      `json:"console"`
	Error   string        `json:"error"`
}

func (pong *Pong) createShot1() shot {
	return shot{
		Me: shotInner{
			X:      pong.player1.x,
			Y:      pong.player1.y,
			VX:     pong.player1.vX,
			VY:     pong.player1.vY,
			Height: pong.player1.height,
			Width:  pong.player1.width,
		},
		Enemy: shotInner{
			X:      pong.player2.x,
			Y:      pong.player2.y,
			VX:     pong.player2.vX,
			VY:     pong.player2.vY,
			Height: pong.player2.height,
			Width:  pong.player1.width,
		},
		Ball: shotInner{
			X:      pong.ball.x,
			Y:      pong.ball.y,
			VX:     pong.ball.vX,
			VY:     pong.ball.vY,
			Height: pong.ball.height,
			Width:  pong.ball.width,
		},
		Game: gameShotInner{
			Height:    pong.height,
			Width:     pong.width,
			TicksLeft: pong.ticksLeft,
		},
	}
}

func (pong *Pong) createShot2() shot {
	return shot{
		Me: shotInner{
			X:      pong.width - pong.player2.x,
			Y:      pong.height - pong.player2.y,
			VX:     -pong.player2.vX,
			VY:     -pong.player2.vY,
			Height: pong.player2.height,
			Width:  pong.player2.width,
		},
		Enemy: shotInner{
			X:      pong.width - pong.player1.x,
			Y:      pong.height - pong.player1.y,
			VX:     -pong.player1.vX,
			VY:     -pong.player1.vY,
			Height: pong.player1.height,
			Width:  pong.player1.width,
		},
		Ball: shotInner{
			X:      pong.width - pong.ball.x,
			Y:      pong.height - pong.ball.y,
			VX:     -pong.ball.vX,
			VY:     -pong.ball.vY,
			Height: pong.ball.height,
			Width:  pong.ball.width,
		},
		Game: gameShotInner{
			Height:    pong.height,
			Width:     pong.width,
			TicksLeft: pong.ticksLeft,
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

// SaveSnapshots сохранение состояния игры
func (pong *Pong) SaveSnapshots(shot1, shot2 []byte) (gameErr error) {
	var s1, s2 shot
	err1 := json.Unmarshal(shot1, &s1)
	if err1 != nil {
		pong.isEnded = true
		pong.Error1 = games.ErrPlayer1Fail.Msg
		return errors.Wrap(err1, games.ErrPlayer1Fail.Error())
	}

	if s1.Error != "" {
		pong.isEnded = true
		pong.Error1 = s1.Error
		pong.winner = 2
		return
	}

	err2 := json.Unmarshal(shot2, &s2)
	if err2 != nil {
		pong.isEnded = true
		pong.Error2 = games.ErrPlayer2Fail.Msg
		return errors.Wrap(err2, games.ErrPlayer2Fail.Error())
	}

	if s2.Error != "" {
		pong.isEnded = true
		pong.Error2 = s2.Error
		pong.winner = 1
		return
	}

	pong.loadSnapShots(&s1, &s2)
	return nil
}

// GetInfo получение информации об объектах игры
func (pong *Pong) GetInfo() games.Info {
	i := &Info{}
	i.Ball.Diameter = pong.ball.height / pong.height
	i.Racket.Height = pong.player1.height / pong.height
	i.Racket.Width = pong.player1.width / pong.width
	i.Ratio = pong.width / pong.height

	return i
}

// GetState получение текущего состояния игры
func (pong *Pong) GetState() (state games.State, fin bool) {
	return &State{
		Player1: object2D{
			X: pong.player1.x / pong.width,
			Y: pong.player1.y / pong.height,
		},
		Player2: object2D{
			X: pong.player2.x / pong.width,
			Y: pong.player2.y / pong.height,
		},
		Ball: object2D{
			X: pong.ball.x / pong.width,
			Y: pong.ball.y / pong.height,
		},
	}, pong.isEnded
}

// GetResult получение результата игры
func (pong *Pong) GetResult() (result games.Result) {
	if !pong.isEnded {
		return nil
	}

	return &Result{
		Winner: pong.winner,
		Err1:   pong.Error1,
		Err2:   pong.Error2,
		Player1: object2D{
			X: pong.player1.x / pong.width,
			Y: pong.player1.y / pong.height,
		},
		Player2: object2D{
			X: pong.player2.x / pong.width,
			Y: pong.player2.y / pong.height,
		},
		Ball: object2D{
			X: pong.ball.x / pong.width,
			Y: pong.ball.y / pong.height,
		},
	}
}
