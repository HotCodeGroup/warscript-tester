package atod

func obstaclesToResp(os []*obstacle, height float64, width float64) (r []obstacleResp) {
	r = make([]obstacleResp, len(os), len(os))
	for i, o := range os {
		r[i] = obstacleResp{
			X:      o.x / width,
			Y:      o.y / height,
			Height: o.height / height,
			Width:  o.width / width,
		}
	}
	return
}

func unitsToResp(us []*unit, height float64, width float64) (r []unitResp) {
	r = make([]unitResp, len(us), len(us))
	for i, u := range us {
		r[i] = unitResp{
			X:           u.x / width,
			Y:           u.y / height,
			Radius:      u.radius / height,
			Health:      u.health,
			ViewRange:   u.viewRange / height,
			ReloadLeft:  u.reloadLeft,
			ReloadTime:  u.reloadTime,
			SpecialLeft: u.specialLeft,
			SpecialTime: u.specialTime,
			UnitType:    u.unitType,
		}
	}
	return
}

func flagsToResp(fs []*flag, height float64, width float64) (r []flagResp) {
	r = make([]flagResp, len(fs), len(fs))
	for i, f := range fs {
		r[i] = flagResp{
			X: f.x / width,
			Y: f.y / height,
		}
	}
	return
}

func dropzoneToResp(d dropzone, height float64, width float64) (r dropzoneResp) {
	r = dropzoneResp{
		X:      d.x / width,
		Y:      d.y / height,
		Radius: d.radius / height,
	}
	return
}

func projectilesToResp(ps []projectile, height float64, width float64) (r []projectileResp) {
	r = make([]projectileResp, len(ps), len(ps))
	for i, p := range ps {
		r[i] = projectileResp{
			X:    p.getX() / width,
			Y:    p.getY() / height,
			VX:   p.getVX() / width,
			VY:   p.getVY() / height,
			Type: p.getType(),
		}
	}
	return
}

func obstaclesToShot(os []*obstacle) (r []obstacleShot) {
	r = make([]obstacleShot, len(os), len(os))
	for i, o := range os {
		r[i] = obstacleShot{
			X:      o.x,
			Y:      o.y,
			Height: o.height,
			Width:  o.width,
		}
	}
	return
}

func unitsToShot(us []*unit) (r []unitShot) {
	r = make([]unitShot, len(us), len(us))
	for i, u := range us {
		r[i] = unitShot{
			IsCarringFlag: u.carriedFlag != nil,
			X:             u.x,
			Y:             u.y,
			Radius:        u.radius,
			VX:            u.vX,
			VY:            u.vY,
			Health:        u.health,
			ViewRange:     u.viewRange,
			MaxSpeed:      u.maxSpeed,
			BulletDamage:  u.bulletSpeed,
			BulletSpeed:   u.bulletSpeed,
			BulletRange:   u.bulletRange,
			ReloadLeft:    u.reloadLeft,
			ReloadTime:    u.reloadTime,
			SpecialLeft:   u.specialLeft,
			SpecialTime:   u.specialTime,
			UnitType:      u.unitType,
		}
	}
	return
}

func flagsToShot(fs []*flag) (r []flagShot) {
	r = make([]flagShot, len(fs), len(fs))
	for i, f := range fs {
		r[i] = flagShot{
			X:         f.x,
			Y:         f.y,
			IsCarried: f.carrier != nil,
		}
	}
	return
}

func dropzoneToShot(d dropzone) (r dropzoneShot) {
	r = dropzoneShot{
		X:      d.x,
		Y:      d.y,
		Radius: d.radius,
	}
	return
}

func projectilesToShot(ps []projectile) (r []projectileShot) {
	r = make([]projectileShot, len(ps), len(ps))
	for i, p := range ps {
		r[i] = projectileShot{
			X:    p.getX(),
			Y:    p.getY(),
			VX:   p.getVX(),
			VY:   p.getVY(),
			Type: p.getType(),
		}
	}
	return
}
