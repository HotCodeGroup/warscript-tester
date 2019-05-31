package pong

import "testing"

var basePlayer = Movable{
	height: 50,
	width:  25,
}
var baseBall = Movable{
	height: 10,
	width:  10,
}

func initPlayer(x, y, vX, vY float64) Movable {
	player := basePlayer
	player.x = x
	player.y = y
	player.vX = vX
	player.vY = vY

	return player
}

func initBall(x, y, vX, vY float64) Movable {
	ball := baseBall
	ball.x = x
	ball.y = y
	ball.vX = vX
	ball.vY = vY

	return ball
}

type collisionTest struct {
	pX  float64
	pY  float64
	pVX float64
	pVY float64
	bX  float64
	bY  float64
	bVX float64
	bVY float64

	isColl   bool
	collSide int
	collX    float64
	collY    float64

	inverse bool
}

var collisionTests = []collisionTest{
	{
		0, 0, 0, 0,
		20, 30, -10, -10,

		true, rightSide, 17.5, 27.5, false,
	},
	{
		0, 0, 0, 0,
		20, -30, -10, 10,

		true, rightSide, 17.5, -27.5, false,
	},
	{
		0, 0, 0, 0,
		20, 30, -1000, -1000,

		true, rightSide, 17.5, 27.5, false,
	},
	{
		0, 0, 0, 0,
		20, -30, -5000, 5000,

		true, rightSide, 17.5, -27.5, false,
	},
	{
		0, 0, 0, 0,
		-20, 30, 10, -10,

		true, leftSide, -17.5, 27.5, false,
	},
	{
		0, 0, 0, 0,
		-20, -30, 10, 10,

		true, leftSide, -17.5, -27.5, false,
	},
	{
		0, 0, 0, 0,
		-20, 30, 1000, -1000,

		true, leftSide, -17.5, 27.5, false,
	},
	{
		0, 0, 0, 0,
		-20, -30, 5000, 5000,

		true, leftSide, -17.5, -27.5, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -10, 0,

		true, rightSide, 17.5, 0, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -1000, 0,

		true, rightSide, 17.5, 0, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -2.5, 30,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, 300,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, 299,

		true, rightSide, 17.5, 29.9, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, -300,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, -299,

		true, rightSide, 17.5, -29.9, false,
	},

	{
		0, 0, 0, 0,
		-20, 0, 10, 0,

		true, leftSide, -17.5, 0, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 1000, 0,

		true, leftSide, -17.5, 0, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 2.5, 30,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, 300,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, 299,

		true, leftSide, -17.5, 29.9, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, -300,

		false, none, 0, 0, false,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, -299,

		true, leftSide, -17.5, -29.9, false,
	},

	// uo-down
	{
		0, 0, 0, 0,
		20, 30, -10, -10,

		true, rightSide, 17.5, 27.5, true,
	},
	{
		0, 0, 0, 0,
		20, -30, -10, 10,

		true, rightSide, 17.5, -27.5, true,
	},
	{
		0, 0, 0, 0,
		20, 30, -1000, -1000,

		true, rightSide, 17.5, 27.5, true,
	},
	{
		0, 0, 0, 0,
		20, -30, -5000, 5000,

		true, rightSide, 17.5, -27.5, true,
	},
	{
		0, 0, 0, 0,
		-20, 30, 10, -10,

		true, leftSide, -17.5, 27.5, true,
	},
	{
		0, 0, 0, 0,
		-20, -30, 10, 10,

		true, leftSide, -17.5, -27.5, true,
	},
	{
		0, 0, 0, 0,
		-20, 30, 1000, -1000,

		true, leftSide, -17.5, 27.5, true,
	},
	{
		0, 0, 0, 0,
		-20, -30, 5000, 5000,

		true, leftSide, -17.5, -27.5, true,
	},

	{
		0, 0, 0, 0,
		20, 0, -10, 0,

		true, rightSide, 17.5, 0, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -1000, 0,

		true, rightSide, 17.5, 0, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -2.5, 30,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, 300,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, 299,

		true, rightSide, 17.5, 29.9, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, -300,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		20, 0, -25, -299,

		true, rightSide, 17.5, -29.9, true,
	},

	{
		0, 0, 0, 0,
		-20, 0, 10, 0,

		true, leftSide, -17.5, 0, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 1000, 0,

		true, leftSide, -17.5, 0, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 2.5, 30,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, 300,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, 299,

		true, leftSide, -17.5, 29.9, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, -300,

		false, none, 0, 0, true,
	},
	{
		0, 0, 0, 0,
		-20, 0, 25, -299,

		true, leftSide, -17.5, -29.9, true,
	},
}

func side(s int) string {
	switch s {
	case upSide:
		return "up  "
	case downSide:
		return "down "
	case rightSide:
		return "right"
	case leftSide:
		return "left "
	case none:
		return "none "
	default:
		return "?????"
	}
}

func TestCollisionParams(t *testing.T) {
	for i, ct := range collisionTests {
		if ct.inverse {
			if ct.collSide == rightSide {
				ct.collSide = upSide
			}
			if ct.collSide == leftSide {
				ct.collSide = downSide
			}
			ct.pVX, ct.pVY = ct.pVY, ct.pVX
			ct.bVX, ct.bVY = ct.bVY, ct.bVX
			ct.pX, ct.pY = ct.pY, ct.pX
			ct.bX, ct.bY = ct.bY, ct.bX
			ct.collX, ct.collY = ct.collY, ct.collX
		}
		player := initPlayer(ct.pX, ct.pY, ct.pVX, ct.pVY)
		ball := initBall(ct.bX, ct.bY, ct.bVX, ct.bVY)
		if ct.inverse {
			player.width, player.height = player.height, player.width
			ball.width, ball.height = ball.height, ball.width
		}
		isColl, collSide, collX, collY := collidePlayerBall(&player, &ball)
		if ct.isColl != isColl || ct.collSide != collSide || ct.collX != collX || ct.collY != collY {
			t.Errorf("test %d\nexpected:\n\tisColl: %v, collSide: %v, "+
				"collX: %v, collY: %+v\ngot:\n\tisColl: %v, "+
				"collSide: %v, collX: %v, collY: %+v\n", i, ct.isColl, side(ct.collSide),
				ct.collX, ct.collY, isColl, side(collSide), collX, collY)
		}
	}

}
