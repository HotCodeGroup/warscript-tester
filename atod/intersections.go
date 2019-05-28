package atod

import "math"

func circleSectionInter(cX float64, cY float64, cR float64,
	x1 float64, y1 float64, x2 float64, y2 float64) (int, float64, float64, float64, float64) {
	x1 -= cX
	x2 -= cX
	y1 -= cY
	y2 -= cY

	a := y2 - y1
	b := x1 - x2
	c := x1*(y1-y2) + y1*(x2-x1)
	x0 := -a * c / (a*a + b*b)
	y0 := -b * c / (a*a + b*b)
	if c*c > cR*cR*(a*a+b*b)+lEPS {
		return 0, 0, 0, 0, 0
	} else if math.Abs(c*c-cR*cR*(a*a+b*b)) < lEPS {
		return 1, x0, y0, 0, 0
	} else {
		d := cR*cR - c*c/(a*a+b*b)
		mult := math.Sqrt(d / (a*a + b*b))
		ax := x0 + b*mult
		bx := x0 - b*mult
		ay := y0 - a*mult
		by := y0 + a*mult
		return 2, ax, ay, bx, by
	}
}

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

func sectionsInter(l1 line, l2 line) (intersect bool, x float64, y float64) {
	v1 := vMult(l2.end.x-l2.beg.x, l2.end.y-l2.beg.y, l1.beg.x-l2.beg.x, l1.beg.y-l2.beg.y)
	v2 := vMult(l2.end.x-l2.beg.x, l2.end.y-l2.beg.y, l1.end.x-l2.beg.x, l1.end.y-l2.beg.y)
	v3 := vMult(l1.end.x-l1.beg.x, l1.end.y-l1.beg.y, l2.beg.x-l1.beg.x, l2.beg.y-l1.beg.y)
	v4 := vMult(l1.end.x-l1.beg.x, l1.end.y-l1.beg.y, l2.end.x-l1.beg.x, l2.end.y-l1.beg.y)

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
