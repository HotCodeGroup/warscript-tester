package pong

import "math"

const speedLim = float64(10)

func (pong *Pong) loadSnapShots(s1 shot, s2 shot) {
	//loading speed of player1
	if s1.Me.VX*s1.Me.VX+s1.Me.VY*s1.Me.VY > speedLim*speedLim {
		coef := math.Sqrt((speedLim * speedLim) / (s1.Me.VX*s1.Me.VX + s1.Me.VY*s1.Me.VY))
		s1.Me.VX *= coef
		s1.Me.VY *= coef
	}
	//loading speed of player2
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
	if res == p1Win {
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
	up
	down
	right
	left
)

func collisionParams(player *Movable, ball *Movable) (isColliding bool, collisionSide int, collisionPointX, collisionPointY float64) {
	// translate to player's fixed system
	ball.vX -= player.vX
	ball.vY -= player.vY
	collisionSide = none
	defer func() {
		ball.vX += player.vX
		ball.vY += player.vY
		if collisionSide == right && ball.vX < 0 {
			ball.vX = -ball.vX
			return
		}
		if collisionSide == left && ball.vX > 0 {
			ball.vX = -ball.vX
			return
		}
		if collisionSide == up && ball.vY < 0 {
			ball.vY = -ball.vY
			return
		}
		if collisionSide == down && ball.vY > 0 {
			ball.vY = -ball.vY
			return
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

	// collision detection
	rTransition := (bLeft-pRight)*(bLeft+ball.vX-pRight) <= 0
	lTransition := (bRight-pLeft)*(bRight+ball.vX-pLeft) <= 0
	uTransition := (bDown-pUp)*(bDown+ball.vY-pUp) <= 0
	dTransition := (bUp-pDown)*(bUp+ball.vY-pDown) <= 0
	if rTransition && uTransition {
		isColliding = true
		if ball.vY*(bLeft+ball.vX-pRight) < ball.vX*(bDown+ball.vY-pUp) {
			collisionSide = right
			collisionPointY = pRight + ball.width
			collisionPointX = ball.y + ball.vY
		} else {
			collisionSide = up
			collisionPointY = pUp + ball.height
			collisionPointX = ball.x + ball.vX
		}
		return
	}
	if rTransition && dTransition {
		isColliding = true
		if ball.vY*(bLeft+ball.vX-pRight) < ball.vX*(bUp+ball.vY-pDown) {
			collisionSide = right
			collisionPointY = pRight + ball.width
			collisionPointX = ball.y + ball.vY
		} else {
			collisionSide = down
			collisionPointY = pDown - ball.height
			collisionPointX = ball.x + ball.vX
		}
		return
	}
	if lTransition && uTransition {
		isColliding = true
		if ball.vY*(bRight+ball.vX-pLeft) < ball.vX*(bDown+ball.vY-pUp) {
			collisionSide = left
			collisionPointY = pLeft - ball.width
			collisionPointX = ball.y + ball.vY
		} else {
			collisionSide = up
			collisionPointY = pUp + ball.height
			collisionPointX = ball.x + ball.vX
		}
		return
	}
	if lTransition && dTransition {
		isColliding = true
		if ball.vY*(bRight+ball.vX-pLeft) < ball.vX*(bUp+ball.vY-pDown) {
			collisionSide = left
			collisionPointY = pLeft - ball.width
			collisionPointX = ball.y + ball.vY
		} else {
			collisionSide = down
			collisionPointY = pDown - ball.height
			collisionPointX = ball.x + ball.vX
		}
		return
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
	// player.x += player.vX
	// player.y += player.vY

	// // controls player not to cross bounds
	// // on x && width
	// if player.x-player.width/2 < left {
	// 	player.x = left + player.width/2
	// }
	// if player.x+player.width/2 > right {
	// 	player.x = right - player.width/2
	// }
	// // on y && height
	// if player.y-player.height/2 < down {
	// 	player.y = down + player.height/2
	// }
	// if player.y+player.height/2 > up {
	// 	player.y = up - player.height/2
	// }
}

const (
	noWinner = iota
	p1Win
	p2Win
)

func moveBall(player1 *Movable, player2 *Movable, ball *Movable, height, width float64) int {
	ball.x += ball.vX
	ball.y += ball.vY
	if ball.y-ball.height < 0 {
		ball.y = ball.height
		ball.vY = -ball.vY
	}
	if ball.y+ball.height > height {
		ball.y = height - ball.height
		ball.vY = -ball.vY
	}
	if ball.x-ball.width < 0 {
		return p2Win
	}
	if ball.x+ball.width > width {
		return p1Win
	}
	return noWinner
}

func (pong *Pong) tick() int {
	collide1, collSide1, collPX1, collPY1 := collisionParams(&pong.player1, &pong.ball)
	collide2, collSide2, collPX2, collPY2 := collisionParams(&pong.player1, &pong.ball)

	if collide1 {
		movePlayerWithBall(&pong.player1, &pong.ball, pong.height, 0, 0, pong.width/3, collSide1, collPX1, collPY1)
		movePlayer(&pong.player2, pong.height, 0, pong.width*2/3, pong.width)
		return moveBall(&pong.player1, &pong.player2, &pong.ball, pong.height, pong.width)
	} else if collide2 {
		movePlayer(&pong.player1, pong.height, 0, 0, pong.width/3)
		movePlayerWithBall(&pong.player2, &pong.ball, pong.height, 0, pong.width*2/3, pong.width, collSide2, collPX2, collPY2)
		return moveBall(&pong.player1, &pong.player2, &pong.ball, pong.height, pong.width)
	} else {
		movePlayer(&pong.player1, pong.height, 0, 0, pong.width/3)
		movePlayer(&pong.player2, pong.height, 0, pong.width*2/3, pong.width)
		return moveBall(&pong.player1, &pong.player2, &pong.ball, pong.height, pong.width)
	}

}
