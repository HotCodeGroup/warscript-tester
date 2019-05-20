package pong

import (
	"math"
)

const speedLim = float64(10)

func (pong *Pong) loadSnapShots(s1, s2 *shot) {
	// loading speed of player1
	if s1.Me.VX*s1.Me.VX+s1.Me.VY*s1.Me.VY > speedLim*speedLim {
		coef := math.Sqrt((speedLim * speedLim) / (s1.Me.VX*s1.Me.VX + s1.Me.VY*s1.Me.VY))
		s1.Me.VX *= coef
		s1.Me.VY *= coef
	}
	// loading speed of player2
	if s2.Me.VX*s2.Me.VX+s2.Me.VY*s2.Me.VY > speedLim*speedLim {
		coef := math.Sqrt((speedLim * speedLim) / (s2.Me.VX*s2.Me.VX + s2.Me.VY*s2.Me.VY))
		s2.Me.VX *= coef
		s2.Me.VY *= coef
	}

	pong.player1.vX = s1.Me.VX
	pong.player1.vY = s1.Me.VY
	pong.player2.vX = -s2.Me.VX
	pong.player2.vY = -s2.Me.VY
	pong.ticksLeft--
	res := pong.tick()
	if res == p1Win {
		pong.winner = 1
		pong.isEnded = true
	}
	if res == p2Win {
		pong.winner = 2
		pong.isEnded = true
	}
	if pong.ticksLeft == 0 {
		pong.winner = 0
		pong.isEnded = true
	}
}

const (
	none = iota
	upSide
	downSide
	rightSide
	leftSide
)

type point struct {
	x float64
	y float64
}

type line struct {
	beg point
	end point
}

func vMult(ax float64, ay float64, bx float64, by float64) float64 {
	return ax*by - bx*ay
}

func intersection(l1 line, l2 line) (intersect bool, x float64, y float64) {
	v1 := vMult(l2.end.x-l2.beg.x, l2.end.y-l2.beg.y, l1.beg.x-l2.beg.x, l1.beg.y-l2.beg.y)
	v2 := vMult(l2.end.x-l2.beg.x, l2.end.y-l2.beg.y, l1.end.x-l2.beg.x, l1.end.y-l2.beg.y)
	v3 := vMult(l1.end.x-l1.beg.x, l1.end.y-l1.beg.y, l2.beg.x-l1.beg.x, l2.beg.y-l1.beg.y)
	v4 := vMult(l1.end.x-l1.beg.x, l1.end.y-l1.beg.y, l2.end.x-l1.beg.x, l2.end.y-l1.beg.y)
	// fmt.Printf("l1: %+v, l2: %+v\n", l1, l2)
	// fmt.Printf("v1: %v, v2: %v, v3: %v, v4: %v\n\n", v1, v2, v3, v4)
	intersect = ((v1*v2) < 0 && (v3*v4) < 0)
	if intersect {
		A1 := l1.end.y - l1.beg.y
		B1 := l1.beg.x - l1.end.x
		C1 := -l1.beg.x*(l1.end.y-l1.beg.y) + l1.beg.y*(l1.end.x-l1.beg.x)

		A2 := l2.end.y - l2.beg.y
		B2 := l2.beg.x - l2.end.x
		C2 := -l2.beg.x*(l2.end.y-l2.beg.y) + l2.beg.y*(l2.end.x-l2.beg.x)

		d := (A1*B2 - B1*A2)
		dx := (-C1*B2 + B1*C2)
		dy := (-A1*C2 + C1*A2)
		x = (dx / d)
		y = (dy / d)
		return
	}
	return
}

const epsilonMove = float64(0.01)

func collidePlayerBall(player *Movable, ball *Movable) (isColliding bool, collisionSide int, collisionPointX, collisionPointY float64) {
	// translate to player's fixed system
	ball.vX -= player.vX
	ball.vY -= player.vY
	// fmt.Printf("ball: %+v\nplayer: %+v\n\n", *ball, *player)
	collisionSide = none
	defer func() {
		ball.vX += player.vX
		ball.vY += player.vY
		if isColliding {
			ball.x = collisionPointX
			ball.y = collisionPointY
			if collisionSide == rightSide {
				if ball.vX < 0 {
					ball.vX = -ball.vX
				}
				ball.x += epsilonMove
				return
			}
			if collisionSide == leftSide {
				if ball.vX > 0 {
					ball.vX = -ball.vX
				}
				ball.x -= epsilonMove
				return
			}
			if collisionSide == upSide {
				if ball.vY < 0 {
					ball.vY = -ball.vY
				}
				ball.y += epsilonMove
				return
			}
			if collisionSide == downSide {
				if ball.vY > 0 {
					ball.vY = -ball.vY
				}
				ball.y -= epsilonMove
				return
			}
		}
	}()

	pRight := player.x + player.width/2
	pLeft := player.x - player.width/2
	pUp := player.y + player.height/2
	pDown := player.y - player.height/2

	bRight := ball.x + ball.width/2
	bLeft := ball.x - ball.width/2
	bUp := ball.y + ball.height/2
	bDown := ball.y - ball.height/2

	pLeftLine := line{point{pLeft, pUp}, point{pLeft, pDown}}
	pRightLine := line{point{pRight, pUp}, point{pRight, pDown}}
	pUpLine := line{point{pLeft, pUp}, point{pRight, pUp}}
	pDownLine := line{point{pLeft, pDown}, point{pRight, pDown}}

	bRightDownLine := line{point{bRight, bDown}, point{bRight + ball.vX, bDown + ball.vY}}
	bLeftDownLine := line{point{bLeft, bDown}, point{bLeft + ball.vX, bDown + ball.vY}}
	bRightUpLine := line{point{bRight, bUp}, point{bRight + ball.vX, bUp + ball.vY}}
	bLeftUpLine := line{point{bLeft, bUp}, point{bLeft + ball.vX, bUp + ball.vY}}

	// collision detection
	if bRight <= pLeft {
		if isColliding, collisionPointX, collisionPointY = intersection(pLeftLine, bRightDownLine); isColliding {
			collisionPointX -= ball.width / 2
			collisionPointY += ball.height / 2

			collisionSide = leftSide
			return
		}
		if isColliding, collisionPointX, collisionPointY = intersection(pLeftLine, bRightUpLine); isColliding {
			collisionPointX -= ball.width / 2
			collisionPointY -= ball.height / 2

			collisionSide = leftSide
			return
		}
	}
	if pRight <= bLeft {
		if isColliding, collisionPointX, collisionPointY = intersection(pRightLine, bLeftDownLine); isColliding {
			collisionPointX += ball.width / 2
			collisionPointY += ball.height / 2

			collisionSide = rightSide
			return
		}
		if isColliding, collisionPointX, collisionPointY = intersection(pRightLine, bLeftUpLine); isColliding {
			collisionPointX += ball.width / 2
			collisionPointY -= ball.height / 2

			collisionSide = rightSide
			return
		}
	}

	if pUp <= bDown {
		if isColliding, collisionPointX, collisionPointY = intersection(pUpLine, bLeftDownLine); isColliding {
			collisionPointX += ball.width / 2
			collisionPointY += ball.height / 2

			collisionSide = upSide
			return
		}
		if isColliding, collisionPointX, collisionPointY = intersection(pUpLine, bRightDownLine); isColliding {
			collisionPointX -= ball.width / 2
			collisionPointY += ball.height / 2

			collisionSide = upSide
			return
		}
	}

	if bUp <= pDown {
		if isColliding, collisionPointX, collisionPointY = intersection(pDownLine, bLeftUpLine); isColliding {
			collisionPointX += ball.width / 2
			collisionPointY -= ball.height / 2

			collisionSide = downSide
			return
		}
		if isColliding, collisionPointX, collisionPointY = intersection(pDownLine, bRightUpLine); isColliding {
			collisionPointX -= ball.width / 2
			collisionPointY -= ball.height / 2

			collisionSide = downSide
			return
		}
	}
	return
}

func movePlayer(player *Movable, up, down, left, right float64) {
	player.x += player.vX
	player.y += player.vY

	// controls player not to cross bounds
	// on x && width
	if player.x-player.width/2 < left {
		player.x = left + player.width/2
	}
	if player.x+player.width/2 > right {
		player.x = right - player.width/2
	}
	// on y && height
	if player.y-player.height/2 < down {
		player.y = down + player.height/2
	}
	if player.y+player.height/2 > up {
		player.y = up - player.height/2
	}
}

func movePlayerWithBall(player *Movable, ball *Movable, up, down, left, right float64, collSide int, collPX, collPY float64) {
	player.x += player.vX
	player.y += player.vY
	ball.x += player.vX
	ball.y += player.vY

	// controls player not to cross bounds
	// on x && width
	if player.x-player.width/2 < left {
		player.x = left + player.width/2
	}
	if player.x+player.width/2 > right {
		player.x = right - player.width/2
	}
	// on y && height
	if player.y-player.height/2 < down {
		player.y = down + player.height/2
	}
	if player.y+player.height/2 > up {
		player.y = up - player.height/2
	}

	if up < ball.y+ball.height/2 && collSide == upSide {
		ball.y = up - ball.height/2 - epsilonMove
		player.y = ball.y - ball.height/2 - epsilonMove - player.height/2

		fullBallV := math.Sqrt(ball.vX*ball.vX + ball.vY*ball.vY)
		ball.vY = 0
		if ball.vX < 0 {
			ball.vX = -fullBallV
		} else {
			ball.vX = fullBallV
		}
	}
	if down > ball.y-ball.height/2 && collSide == downSide {
		ball.y = down + ball.height/2 + epsilonMove
		player.y = ball.y + ball.height/2 + epsilonMove + player.height/2

		fullBallV := math.Sqrt(ball.vX*ball.vX + ball.vY*ball.vY)
		ball.vY = 0
		if ball.vX < 0 {
			ball.vX = -fullBallV
		} else {
			ball.vX = fullBallV
		}
	}

	if math.Abs(player.vX) < math.Abs(ball.vX) {
		ball.x += ball.vX - player.vX
	}
	if math.Abs(player.vY) < math.Abs(ball.vY) {
		ball.y += ball.vY - player.vY
	}
}

const (
	noWinner = iota
	p1Win
	p2Win
)

func fixBallPos(ball *Movable, height float64) bool {
	if ball.y-ball.height < 0 {
		ball.y = ball.height
		ball.vY = -ball.vY
		return true
	}
	if ball.y+ball.height > height {
		ball.y = height - ball.height
		ball.vY = -ball.vY
		return true
	}
	return false
}

func winnerCheck(ball *Movable, width float64) int {
	if (ball.x-ball.width)/2 < 0 {
		return p2Win
	}
	if (ball.x+ball.width)/2 > width {
		return p1Win
	}
	return noWinner
}

func (pong *Pong) tick() int {
	collide1, collSide1, collPX1, collPY1 := collidePlayerBall(&pong.player1, &pong.ball)
	collide2, collSide2, collPX2, collPY2 := collidePlayerBall(&pong.player2, &pong.ball)

	if collide1 {
		movePlayerWithBall(&pong.player1, &pong.ball, pong.height, 0, 0, pong.width/3, collSide1, collPX1, collPY1)
		movePlayer(&pong.player2, pong.height, 0, pong.width*2/3, pong.width)
	} else if collide2 {
		movePlayer(&pong.player1, pong.height, 0, 0, pong.width/3)
		movePlayerWithBall(&pong.player2, &pong.ball, pong.height, 0, pong.width*2/3, pong.width, collSide2, collPX2, collPY2)
	} else {
		movePlayer(&pong.player1, pong.height, 0, 0, pong.width/3)
		movePlayer(&pong.player2, pong.height, 0, pong.width*2/3, pong.width)
		pong.ball.y += pong.ball.vY
		if fixBallPos(&pong.ball, pong.height) {
			collide1, collSide1, collPX1, collPY1 := collidePlayerBall(&pong.player1, &pong.ball)
			collide2, collSide2, collPX2, collPY2 := collidePlayerBall(&pong.player2, &pong.ball)
			if collide1 {
				movePlayerWithBall(&pong.player1, &pong.ball, pong.height, 0, 0, pong.width/3, collSide1, collPX1, collPY1)
				movePlayer(&pong.player2, pong.height, 0, pong.width*2/3, pong.width)
			} else if collide2 {
				movePlayer(&pong.player1, pong.height, 0, 0, pong.width/3)
				movePlayerWithBall(&pong.player2, &pong.ball, pong.height,
					0, pong.width*2/3, pong.width, collSide2, collPX2, collPY2)
			}
		} else {
			pong.ball.x += pong.ball.vX
		}
	}

	return winnerCheck(&pong.ball, pong.width)
}
