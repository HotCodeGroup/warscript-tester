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

}
