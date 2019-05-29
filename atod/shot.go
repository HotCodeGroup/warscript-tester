package atod

type obstacleShot struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Height float64 `json:"height"`
	Width  float64 `json:"width"`
}

func (s *obstacleShot) inverse(height float64, width float64) {
	s.X = width - s.X
	s.Y = height - s.Y
}

type flagShot struct {
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	IsCarried bool    `json:"is_carried"`
}

func (s *flagShot) inverse(height float64, width float64) {
	s.X = width - s.X
	s.Y = height - s.Y
}

type unitShot struct {
	IsCarringFlag bool    `json:"is_carring_flag"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Radius        float64 `json:"radius"`
	VX            float64 `json:"vX"`
	VY            float64 `json:"vY"`
	Health        float64 `json:"health"`
	ViewRange     float64 `json:"view_range"`
	MaxSpeed      float64 `json:"max_speed"`
	BulletDamage  float64 `json:"bullet_damage"`
	BulletSpeed   float64 `json:"bullet_speed"`
	BulletRange   float64 `json:"bullet_range"`
	ReloadTime    float64 `json:"reload_time"`
	ReloadLeft    float64 `json:"reload_left"`
	SpecialTime   float64 `json:"special_time"`
	SpecialLeft   float64 `json:"special_left"`
	UnitType      string  `json:"unit_type"`
	BulletDirX    float64 `json:"bullet_dir_x"`
	BulletDirY    float64 `json:"bullet_dir_y"`
	SpecialDirX   float64 `json:"special_dir_x"`
	SpecialDirY   float64 `json:"special_dir_y"`
}

func (s *unitShot) inverse(height float64, width float64) {
	s.X = width - s.X
	s.Y = height - s.Y
	s.VX = -s.VX
	s.VY = -s.VY
}

type dropzoneShot struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Radius float64 `json:"radius"`
}

func (s *dropzoneShot) inverse(height float64, width float64) {
	s.X = width - s.X
	s.Y = height - s.Y
}

type projectileShot struct {
	Type string
	X    float64
	Y    float64
	VX   float64
	VY   float64
}

func (s *projectileShot) inverse(height float64, width float64) {
	s.X = width - s.X
	s.Y = height - s.Y
	s.VX = -s.VX
	s.VY = -s.VY
}

type shot struct {
	Dropzone      dropzoneShot     `json:"dropzone"`
	EnemyDropzone dropzoneShot     `json:"enemy_dropzone"`
	Projectiles   []projectileShot `json:"projectiles"`
	Obstacles     []obstacleShot   `json:"obstacles"`
	Units         []unitShot       `json:"units"`
	EnemyUnits    []unitShot       `json:"enemy_units"`
	Flags         []flagShot       `json:"flags"`
	EnemyFlags    []flagShot       `json:"enemy_flags"`
}

func (a *Atod) createShot1() shot {
	return shot{
		Dropzone:      dropzoneToShot(a.dropzone1),
		EnemyDropzone: dropzoneToShot(a.dropzone2),
		Projectiles:   projectilesToShot(a.projectiles),
		Obstacles:     obstaclesToShot(a.obstacles),
		Units:         unitsToShot(a.player1Units),
		EnemyUnits:    unitsToShot(a.player2Units),
		Flags:         flagsToShot(a.flags1),
		EnemyFlags:    flagsToShot(a.flags2),
	}
}

func (a *Atod) createShot2() shot {
	s := shot{
		Dropzone:      dropzoneToShot(a.dropzone2),
		EnemyDropzone: dropzoneToShot(a.dropzone1),
		Projectiles:   projectilesToShot(a.projectiles),
		Obstacles:     obstaclesToShot(a.obstacles),
		Units:         unitsToShot(a.player2Units),
		EnemyUnits:    unitsToShot(a.player1Units),
		Flags:         flagsToShot(a.flags2),
		EnemyFlags:    flagsToShot(a.flags1),
	}
	s.Dropzone.inverse(a.heihgt, a.width)
	s.EnemyDropzone.inverse(a.heihgt, a.width)
	for i := range s.Projectiles {
		s.Projectiles[i].inverse(a.heihgt, a.width)
	}
	for i := range s.Obstacles {
		s.Obstacles[i].inverse(a.heihgt, a.width)
	}
	for i := range s.Units {
		s.Units[i].inverse(a.heihgt, a.width)
	}
	for i := range s.EnemyUnits {
		s.EnemyUnits[i].inverse(a.heihgt, a.width)
	}
	for i := range s.Flags {
		s.Flags[i].inverse(a.heihgt, a.width)
	}
	for i := range s.EnemyFlags {
		s.EnemyFlags[i].inverse(a.heihgt, a.width)
	}

	return s
}
