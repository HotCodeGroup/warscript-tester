package atod

import "math"

func (a *Atod) loadSnapShots(s1, s2 *shot) {
	for i := range s2.Units {
		s2.Units[i].inverse(a.heihgt, a.width)
	}

	// loading movement
	loadSpeed(s1.Units, a.player1Units)
	loadSpeed(s2.Units, a.player2Units)
	// loading projectiles
	a.loadProjectiles(s1.Units, a.player1Units)
	a.loadProjectiles(s2.Units, a.player2Units)

	loadFlags(s1.Units, a.player1Units, a.flags1)
	loadFlags(s1.Units, a.player1Units, a.flags2)
	loadFlags(s2.Units, a.player2Units, a.flags2)
	loadFlags(s2.Units, a.player2Units, a.flags1)

	a.moveProjectiles()
	a.moveUnits(a.player1Units)
	a.moveUnits(a.player2Units)

	a.ticksLeft--
	res := a.checkWinner()
	if res == 1 {
		a.winner = 1
		a.isEnded = true
	}
	if res == 2 {
		a.winner = 2
		a.isEnded = true
	}
	if a.ticksLeft == 0 {
		a.winner = 0
		a.isEnded = true
	}
}

func loadFlags(shot []unitShot, un []*unit, fs []*flag) {
	amount := len(shot)
	if amount > len(un) {
		amount = len(un)
	}
	for i := 0; i < amount; i++ {
		if un[i].health <= 0 {
			continue
		}
		if shot[i].IsCarringFlag && un[i].carriedFlag == nil {
			for _, f := range fs {
				if f.carrier == nil && math.Abs(f.x-un[i].x) < un[i].radius && math.Abs(f.y-un[i].y) < un[i].radius {
					f.carrier = un[i]
					un[i].carriedFlag = f
					f.x = un[i].x
					f.y = un[i].y
				}
			}
		}
		if !shot[i].IsCarringFlag && un[i].carriedFlag != nil {
			un[i].carriedFlag.carrier = nil
			un[i].carriedFlag = nil
		}
	}
}

func loadSpeed(shot []unitShot, un []*unit) {
	amount := len(shot)
	if amount > len(un) {
		amount = len(un)
	}
	for i := 0; i < amount; i++ {
		if un[i].health <= 0 {
			continue
		}
		vX := shot[i].VX
		vY := shot[i].VY
		v := math.Sqrt(vX*vX + vY*vY)
		if v > un[i].maxSpeed {
			vX *= un[i].maxSpeed / v
			vY *= un[i].maxSpeed / v
		}
		un[i].vX = vX
		un[i].vY = vY
	}
}

func (a *Atod) loadProjectiles(shot []unitShot, un []*unit) {
	amount := len(shot)
	if amount > len(un) {
		amount = len(un)
	}
	for i := 0; i < amount; i++ {
		if un[i].health <= 0 || un[i].carriedFlag != nil {
			continue
		}
		if un[i].reloadLeft == 0 && (shot[i].BulletDirX != 0 || shot[i].BulletDirY != 0) {
			un[i].reloadLeft = un[i].reloadTime
			p := un[i].shot(shot[i].BulletDirX, shot[i].BulletDirY)
			if p != nil {
				a.projectiles = append(a.projectiles, p)
			}
		}
	}
}

func (a *Atod) moveProjectiles() {
	ps := make([]projectile, 0, len(a.projectiles))
	for _, p := range a.projectiles {
		if p.move() {
			coll := false
			for _, o := range a.obstacles {
				if p.obstacleIntersect(o) {
					coll = true
					continue
				}
			}
			for _, u := range a.player1Units {
				if p.unitIntersect(u) {
					coll = true
					u.health -= p.getDamage()
					if u.health <= 0 && u.carriedFlag != nil {
						u.carriedFlag.carrier = nil
						u.carriedFlag = nil
					}
				}
			}
			for _, u := range a.player2Units {
				if p.unitIntersect(u) {
					coll = true
					u.health -= p.getDamage()
					if u.health <= 0 && u.carriedFlag != nil {
						u.carriedFlag.carrier = nil
						u.carriedFlag = nil
					}
				}
			}
			if !coll &&
				p.getX() < a.heihgt && 0 < p.getX() &&
				p.getY() < a.width && 0 < p.getY() {
				ps = append(ps, p)
			}
		}
	}

	a.projectiles = ps
}

func (a *Atod) moveUnits(un []*unit) {
	for _, u := range un {
		if u.health <= 0 {
			continue
		}

		deltaX := u.vX
		deltaY := u.vY
		movedY := true
		for movedY {
			_, deltaX = a.moveUnitX(u, deltaX)
			movedY, deltaY = a.moveUnitX(u, deltaY)
		}

		if u.carriedFlag != nil {
			u.carriedFlag.x = u.x
			u.carriedFlag.y = u.y
		}
	}
}

func (a *Atod) moveUnitX(u *unit, delta float64) (bool, float64) {
	if math.Abs(delta) < lEPS {
		return false, 0
	}
	movement := delta
	adj := -lEPS
	if delta > 0 {
		for _, obst := range a.obstacles {
			movement = math.Min(
				movement,
				moveCircle(u.x, u.y, u.radius,
					obst.x-obst.width/2, obst.y-obst.height/2, obst.y+obst.height/2, delta),
			)
		}
	} else {
		adj = lEPS
		for _, obst := range a.obstacles {
			movement = math.Min(
				movement,
				moveCircle(u.x, u.y, u.radius,
					obst.x+obst.width/2, obst.y+obst.height/2, obst.y+obst.height/2, delta),
			)
		}
	}

	u.x += movement + adj
	return math.Abs(movement) > lEPS, delta - movement - adj
}

func (a *Atod) moveUnitY(u *unit, delta float64) (bool, float64) {
	if math.Abs(delta) < lEPS {
		return false, 0
	}
	movement := delta
	adj := -lEPS
	if delta > 0 {
		for _, obst := range a.obstacles {
			movement = math.Min(
				movement,
				moveCircle(u.y, u.x, u.radius,
					obst.y-obst.height/2, obst.x-obst.width/2, obst.x+obst.width/2, delta),
			)
		}
	} else {
		adj = lEPS
		for _, obst := range a.obstacles {
			movement = math.Min(
				movement,
				moveCircle(u.y, u.x, u.radius,
					obst.y+obst.height/2, obst.x-obst.width/2, obst.x+obst.width/2, delta),
			)
		}
	}

	u.y += movement + adj
	return math.Abs(movement) > lEPS, delta - movement - adj
}

func (a *Atod) checkWinner() int {
	f1 := make([]*flag, 0, 0)
	for _, f := range a.flags1 {
		if !(f.carrier == nil &&
			math.Abs(f.x-a.dropzone2.x) < a.dropzone2.radius &&
			math.Abs(f.y-a.dropzone2.y) < a.dropzone2.radius) {

			f1 = append(f1, f)
		}
		a.flags1 = f1
	}
	f2 := make([]*flag, 0, 0)
	for _, f := range a.flags2 {
		if !(f.carrier == nil &&
			math.Abs(f.x-a.dropzone2.x) < a.dropzone1.radius &&
			math.Abs(f.y-a.dropzone2.y) < a.dropzone1.radius) {

			f2 = append(f2, f)
		}
		a.flags2 = f2
	}
	p1ct := len(a.player1Units)
	for _, u := range a.player1Units {
		if u.health <= 0 {
			p1ct--
		}
	}
	p2ct := len(a.player2Units)
	for _, u := range a.player2Units {
		if u.health <= 0 {
			p2ct--
		}
	}

	if len(f1) == 0 || len(f2) == 0 || p1ct == 0 || p2ct == 0 {
		if len(f1) == len(f2) {
			if p1ct > p2ct {
				return 1
			} else if p1ct < p2ct {
				return 2
			} else {
				return 0
			}
		}
		if p1ct == 0 && p2ct == 0 {
			if len(f1) > len(f2) {
				return 1
			} else if len(f1) < len(f2) {
				return 2
			} else {
				return 0
			}
		}
		if len(f2) == 0 || p2ct == 0 {
			return 1
		}
		if len(f1) == 0 || p1ct == 0 {
			return 2
		}
	}
	return 0
}
