package atod

import (
	"math"
)

const (
	stdObstAmount     = 0
	stdHeight         = 1000
	stdWidth          = 2000
	stdNoObstacleZone = 350
	stdSafeZone       = 350
	stdDropZoneRad    = 50
	stdDropZoneX      = 100
	stdDropZoneY      = 500
	stdUnitRad        = 25
	stdFlag1PosX      = 100
	stdFlag1PosY      = 100
	stdFlag2PosX      = 100
	stdFlag2PosY      = 900

	stdObstHeight = 200
	stdObstWidth  = 200
	stdObst1PosX  = 500
	stdObst1PosY  = 300
	stdObst2PosX  = 500
	stdObst2PosY  = 700
)

func stdField() *Atod {
	return &Atod{
		ticksLeft: 10000,
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

		projectiles: make([]projectile, 0, 0),

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
		winner:  0,
		isEnded: false,
	}
}

func stdObstacles() []*obstacle {
	return []*obstacle{
		{
			width:  stdObstWidth,
			height: stdObstHeight,
			x:      stdObst1PosX,
			y:      stdObst1PosY,
		},
		{
			width:  stdObstWidth,
			height: stdObstHeight,
			x:      stdObst2PosX,
			y:      stdObst2PosY,
		},
		{
			width:  stdObstWidth,
			height: stdObstHeight,
			x:      stdWidth - stdObst1PosX,
			y:      stdHeight - stdObst1PosY,
		},
		{
			width:  stdObstWidth,
			height: stdObstHeight,
			x:      stdWidth - stdObst2PosX,
			y:      stdHeight - stdObst2PosY,
		},
	}
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
	b.prevX = b.x
	b.prevY = b.y
	if b.distLeft <= 0 {
		newV := b.distLeft + b.v
		b.vX *= newV / b.v
		b.vY *= newV / b.v
	}
	b.x += b.vX
	b.y += b.vY

	return true
}

func (b *stdBullet) getType() string {
	return "bullet"
}

func (b *stdBullet) getX() float64 {
	return b.x
}

func (b *stdBullet) getY() float64 {
	return b.y
}

func (b *stdBullet) getVX() float64 {
	return b.vX
}

func (b *stdBullet) getVY() float64 {
	return b.vY
}

func (b *stdBullet) getDamage() float64 {
	return b.damage
}

func bulletProducer(u *unit) func(float64, float64) projectile {
	return func(x float64, y float64) projectile {
		if x == 0 && y == 0 {
			return nil
		}

		mod := math.Sqrt(x*x+y*y) * (1 - lEPS)
		dirX := x / mod
		dirY := y / mod

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
		unitType:     "soldier1",
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
		unitType:     "soldier2",
	}
	sld2.shot = bulletProducer(sld2)

	units := make([]*unit, 0, 5)
	units = append(units, sniper, healer, tank, sld1, sld2)
	if reversed {
		for _, u := range units {
			u.x = stdWidth - u.x
			u.y = stdHeight - u.y
		}
	}
	return units
}
