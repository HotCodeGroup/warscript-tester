package atod

import (
	"math"
	"math/rand"
	"time"
)

const (
	stdObstAmount  = 10
	stdHeight      = 1000
	stdWidth       = 2000
	stdSafeZone    = 350
	stdDropZoneRad = 50
	stdDropZoneX   = 100
	stdDropZoneY   = 500
	stdUnitRad     = 25
	stdFlag1PosX   = 100
	stdFlag1PosY   = 100
	stdFlag2PosX   = 100
	stdFlag2PosY   = 900
)

func stdField() *Atod {
	return &Atod{
		heihgt:    stdHeight,
		width:     stdWidth,
		obstacles: stdObstacles(),
		dropzone1: dropzone{
			x:      stdDropZoneX,
			y:      stdDropZoneY,
			radius: stdDropZoneRad,
		},
		dropzone2: dropzone{
			x:      stdWidth - stdDropZoneX,
			y:      stdHeight - stdDropZoneY,
			radius: stdDropZoneRad,
		},

		projectiles: make([]*projectile, 0, 0),

		player1Units: stdUnits(false),
		player2Units: stdUnits(true),
		flags1: []*flag{
			&flag{
				x:       stdFlag1PosX,
				y:       stdFlag1PosY,
				carrier: nil,
			},
			&flag{
				x:       stdFlag2PosX,
				y:       stdFlag2PosY,
				carrier: nil,
			},
		},
		flags2: []*flag{
			&flag{
				x:       stdWidth - stdFlag1PosX,
				y:       stdHeight - stdFlag1PosY,
				carrier: nil,
			},
			&flag{
				x:       stdWidth - stdFlag2PosX,
				y:       stdHeight - stdFlag2PosY,
				carrier: nil,
			},
		},
		winner:       0,
		isEnded:      false,
		occuredError: nil,
	}
}

func stdObstacles() []*obstacle {
	obstacles := make([]*obstacle, 0, stdObstAmount*2)
	randS := rand.NewSource(time.Now().UnixNano())
	for i := 0; i < stdObstAmount; i++ {
		height := float64(randS.Int63() % int64(stdHeight/20))
		width := float64(randS.Int63() % int64(stdWidth/20))
		x := float64(randS.Int63() % int64(stdHeight))
		y := float64(randS.Int63() % int64(stdWidth))

		obstacles = append(obstacles, &obstacle{
			width:  width,
			height: height,
			x:      x,
			y:      y,
		})

		obstacles = append(obstacles, &obstacle{
			width:  width,
			height: height,
			x:      stdHeight - x,
			y:      stdWidth - y,
		})
	}

	return obstacles
}

type stdBullet struct {
	prevX    float64
	prevY    float64
	x        float64
	y        float64
	vX       float64
	vY       float64
	v        float64
	distLeft float64
	damage   float64
}

func (b *stdBullet) unitIntersect(u *unit) bool {
	p, _, _, _, _ := circleSectionInter(u.x, u.y, u.radius, b.prevX, b.prevY, b.x, b.y)
	return p > 0
}

func (b *stdBullet) obstacleIntersect(o *obstacle) bool {
	bulletLine := line{
		beg: point{b.prevX, b.prevY},
		end: point{b.x, b.y},
	}

	deltX := o.width / 2
	deltY := o.height / 2

	obstacleLine1 := line{
		beg: point{o.x + deltX, o.y + deltY},
		end: point{o.x - deltX, o.y + deltY},
	}
	if p, _, _ := sectionsInter(bulletLine, obstacleLine1); p {
		return true
	}
	obstacleLine2 := line{
		beg: point{o.x - deltX, o.y + deltY},
		end: point{o.x - deltX, o.y - deltY},
	}
	if p, _, _ := sectionsInter(bulletLine, obstacleLine2); p {
		return true
	}
	obstacleLine3 := line{
		beg: point{o.x - deltX, o.y - deltY},
		end: point{o.x + deltX, o.y - deltY},
	}
	if p, _, _ := sectionsInter(bulletLine, obstacleLine3); p {
		return true
	}
	obstacleLine4 := line{
		beg: point{o.x + deltX, o.y - deltY},
		end: point{o.x + deltX, o.y + deltY},
	}
	if p, _, _ := sectionsInter(bulletLine, obstacleLine4); p {
		return true
	}
	return false
}

func (b *stdBullet) move() bool {
	if b.distLeft <= 0 {
		return false
	}
	b.distLeft -= b.v
	b.prevX = b.vX
	b.prevY = b.vY
	if b.distLeft <= 0 {
		newV := b.distLeft + b.v
		b.vX *= newV / b.v
		b.vY *= newV / b.v
	}
	b.x += b.vX
	b.y += b.vY

	return true
}

func bulletProducer(u *unit) func(float64, float64) projectile {
	return func(x float64, y float64) projectile {
		if x == 0 && y == 0 {
			return nil
		}

		mod := math.Sqrt(x*x+y*y) + lEPS
		dirX := mod / x
		dirY := mod / y

		return &stdBullet{
			v:        u.bulletSpeed,
			x:        u.x + dirX*u.radius,
			y:        u.y + dirY*u.radius,
			vX:       dirX * u.bulletSpeed,
			vY:       dirY * u.bulletSpeed,
			damage:   u.bulletDamage,
			distLeft: u.bulletRange,
		}
	}
}

func stdUnits(reversed bool) []*unit {

	sniper := &unit{
		x:            125,
		y:            625,
		radius:       stdUnitRad,
		maxSpeed:     15,
		vX:           0,
		vY:           0,
		health:       500,
		viewRange:    500,
		bulletDamage: 800,
		bulletSpeed:  200,
		bulletRange:  2000,
		reloadTime:   50,
		reloadLeft:   0,
		specialTime:  0,
		specialLeft:  0,
		unitType:     "sniper",
	}
	sniper.shot = bulletProducer(sniper)
	healer := &unit{
		x:            125,
		y:            375,
		radius:       stdUnitRad,
		maxSpeed:     15,
		vX:           0,
		vY:           0,
		health:       600,
		viewRange:    200,
		bulletDamage: 50,
		bulletSpeed:  50,
		bulletRange:  200,
		reloadTime:   25,
		reloadLeft:   0,
		specialTime:  0,
		specialLeft:  0,
		unitType:     "medic",
	}
	healer.shot = bulletProducer(healer)
	tank := &unit{
		x:            225,
		y:            500,
		radius:       stdUnitRad * 2,
		maxSpeed:     15,
		vX:           0,
		vY:           0,
		health:       2000,
		viewRange:    200,
		bulletDamage: 200,
		bulletSpeed:  50,
		bulletRange:  200,
		reloadTime:   25,
		reloadLeft:   0,
		specialTime:  0,
		specialLeft:  0,
		unitType:     "tank",
	}
	tank.shot = bulletProducer(tank)
	sld1 := &unit{
		x:            200,
		y:            575,
		radius:       stdUnitRad,
		maxSpeed:     40,
		vX:           0,
		vY:           0,
		health:       1000,
		viewRange:    200,
		bulletDamage: 400,
		bulletSpeed:  50,
		bulletRange:  300,
		reloadTime:   25,
		reloadLeft:   0,
		specialTime:  0,
		specialLeft:  0,
		unitType:     "soldier",
	}
	sld1.shot = bulletProducer(sld1)
	sld2 := &unit{
		x:            200,
		y:            425,
		radius:       stdUnitRad,
		maxSpeed:     40,
		vX:           0,
		vY:           0,
		health:       1000,
		viewRange:    200,
		bulletDamage: 400,
		bulletSpeed:  50,
		bulletRange:  300,
		reloadTime:   25,
		reloadLeft:   0,
		specialTime:  0,
		specialLeft:  0,
		unitType:     "soldier",
	}
	sld2.shot = bulletProducer(sld2)

	units := make([]*unit, 0, 5)
	units = append(units, sniper, healer, tank, sld1, sld2)
	if reversed {
		for _, u := range units {
			u.x = stdWidth - u.x
			u.y = stdWidth - u.y
		}
	}
	return units
}
