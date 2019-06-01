package atod

type obstacle struct {
	width  float64
	height float64
	x      float64
	y      float64
}

type unit struct {
	carriedFlag  *flag
	x            float64
	y            float64
	radius       float64
	vX           float64
	vY           float64
	health       float64
	viewRange    float64
	maxSpeed     float64
	bulletDamage float64
	bulletSpeed  float64
	bulletRange  float64
	reloadTime   float64
	reloadLeft   float64
	specialTime  float64
	specialLeft  float64
	unitType     string
	shot         func(float64, float64) projectile
	special      func(float64, float64) projectile
}

type projectile interface {
	unitIntersect(*unit) bool
	obstacleIntersect(*obstacle) bool
	move() bool

	getType() string
	getX() float64
	getY() float64
	getVX() float64
	getVY() float64
	getDamage() float64
}

type dropzone struct {
	x      float64
	y      float64
	radius float64
}

type flag struct {
	x       float64
	y       float64
	carrier *unit
}

// Atod - game of 2AtoD
type Atod struct {
	ticksLeft int
	width     float64
	heihgt    float64
	obstacles []*obstacle
	dropzone1 dropzone
	dropzone2 dropzone

	projectiles []projectile

	player1Units []*unit
	player2Units []*unit
	flags1       []*flag
	flags2       []*flag

	logs1 []string
	logs2 []string

	winner  int
	isEnded bool
	Error1  string
	Error2  string
}
