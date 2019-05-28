package atod

func obstaclesToResp(os []*obstacle) (r []obstacleResp) {
	r = make([]obstacleResp, len(os), len(os))
	for i, o := range os {
		r[i] = obstacleResp{
			X:      o.x,
			Y:      o.y,
			Height: o.height,
			Width:  o.width,
		}
	}
	return
}

func unitsToResp(us []*unit) (r []unitResp) {
	r = make([]unitResp, len(us), len(us))
	for i, u := range us {
		r[i] = unitResp{
			X:           u.x,
			Y:           u.y,
			Radius:      u.radius,
			Health:      u.health,
			ViewRange:   u.viewRange,
			ReloadLeft:  u.reloadLeft,
			ReloadTime:  u.reloadTime,
			SpecialLeft: u.specialLeft,
			SpecialTime: u.specialTime,
			UnitType:    u.unitType,
		}
	}
	return
}

func flagsToResp(fs []*flag) (r []flagResp) {
	r = make([]flagResp, len(fs), len(fs))
	for i, f := range fs {
		r[i] = flagResp{
			X: f.x,
			Y: f.y,
		}
	}
	return
}

func dropzoneToResp(d dropzone) (r dropzoneResp) {
	r = dropzoneResp{
		X:      d.x,
		Y:      d.y,
		Radius: d.radius,
	}
	return
}
