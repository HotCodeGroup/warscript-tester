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
	collisionTest{
		0, 0, 0, 0,
		20, 30, -10, -10,

		true, right, 17.5, 27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, -30, -10, 10,

		true, right, 17.5, -27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 30, -1000, -1000,

		true, right, 17.5, 27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, -30, -5000, 5000,

		true, right, 17.5, -27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 30, 10, -10,

		true, left, -17.5, 27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, -30, 10, 10,

		true, left, -17.5, -27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 30, 1000, -1000,

		true, left, -17.5, 27.5, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, -30, 5000, 5000,

		true, left, -17.5, -27.5, false,
	},

	collisionTest{
		0, 0, 0, 0,
		20, 0, -10, 0,

		true, right, 17.5, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -1000, 0,

		true, right, 17.5, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -2.5, 30,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, 300,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, 299,

		true, right, 17.5, 29.9, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, -300,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, -299,

		true, right, 17.5, -29.9, false,
	},

	collisionTest{
		0, 0, 0, 0,
		-20, 0, 10, 0,

		true, left, -17.5, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 1000, 0,

		true, left, -17.5, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 2.5, 30,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, 300,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, 299,

		true, left, -17.5, 29.9, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, -300,

		false, none, 0, 0, false,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, -299,

		true, left, -17.5, -29.9, false,
	},

	//uo-down
	collisionTest{
		0, 0, 0, 0,
		20, 30, -10, -10,

		true, right, 17.5, 27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, -30, -10, 10,

		true, right, 17.5, -27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 30, -1000, -1000,

		true, right, 17.5, 27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, -30, -5000, 5000,

		true, right, 17.5, -27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 30, 10, -10,

		true, left, -17.5, 27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, -30, 10, 10,

		true, left, -17.5, -27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 30, 1000, -1000,

		true, left, -17.5, 27.5, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, -30, 5000, 5000,

		true, left, -17.5, -27.5, true,
	},

	collisionTest{
		0, 0, 0, 0,
		20, 0, -10, 0,

		true, right, 17.5, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -1000, 0,

		true, right, 17.5, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -2.5, 30,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, 300,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, 299,

		true, right, 17.5, 29.9, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, -300,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		20, 0, -25, -299,

		true, right, 17.5, -29.9, true,
	},

	collisionTest{
		0, 0, 0, 0,
		-20, 0, 10, 0,

		true, left, -17.5, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 1000, 0,

		true, left, -17.5, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 2.5, 30,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, 300,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, 299,

		true, left, -17.5, 29.9, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, -300,

		false, none, 0, 0, true,
	},
	collisionTest{
		0, 0, 0, 0,
		-20, 0, 25, -299,

		true, left, -17.5, -29.9, true,
	},
}

func side(s int) string {
	switch s {
	case up:
		return "up  "
	case down:
		return "down "
	case right:
		return "right"
	case left:
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
			if ct.collSide == right {
				ct.collSide = up
			}
			if ct.collSide == left {
				ct.collSide = down
			}
			t := ct.pVX
			ct.pVX = ct.pVY
			ct.pVY = t

			t = ct.bVX
			ct.bVX = ct.bVY
			ct.bVY = t

			t = ct.pX
			ct.pX = ct.pY
			ct.pY = t

			t = ct.bX
			ct.bX = ct.bY
			ct.bY = t

			t = ct.collX
			ct.collX = ct.collY
			ct.collY = t
		}
		player := initPlayer(ct.pX, ct.pY, ct.pVX, ct.pVY)
		ball := initBall(ct.bX, ct.bY, ct.bVX, ct.bVY)
		if ct.inverse {
			t := player.width
			player.width = player.height
			player.height = t

			t = ball.width
			ball.width = ball.height
			ball.height = t
		}
		isColl, collSide, collX, collY := collisionParams(&player, &ball)
		if ct.isColl != isColl || ct.collSide != collSide || ct.collX != collX || ct.collY != collY {
			t.Errorf("test %d\nexpected:\n\tisColl: %v, collSide: %v, collX: %v, collY: %+v\ngot:\n\tisColl: %v, collSide: %v, collX: %v, collY: %+v\n", i, ct.isColl, side(ct.collSide), ct.collX, ct.collY, isColl, side(collSide), collX, collY)
		}
	}

}
