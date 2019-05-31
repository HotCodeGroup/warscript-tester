package pong

import "testing"

type intersectionTest struct {
	line1 line
	line2 line
	ans   bool
	x     float64
	y     float64
}

var intersectionTests = []intersectionTest{
	{
		line1: line{
			beg: point{
				x: 1, y: 1,
			},
			end: point{
				x: 0, y: 0,
			},
		},
		line2: line{
			beg: point{
				x: 0, y: 1,
			},
			end: point{
				x: 1, y: 0,
			},
		},
		ans: true,
	},
	{
		line1: line{
			beg: point{
				x: 1, y: 1,
			},
			end: point{
				x: 2, y: 2,
			},
		},
		line2: line{
			beg: point{
				x: 0, y: 1,
			},
			end: point{
				x: 1, y: 0,
			},
		},
		ans: false,
	},
	{
		line1: line{
			beg: point{
				x: 0, y: 0,
			},
			end: point{
				x: 10, y: 0,
			},
		},
		line2: line{
			beg: point{
				x: 5, y: 5,
			},
			end: point{
				x: 5, y: -5,
			},
		},
		ans: true,
	},
	{
		line1: line{
			beg: point{
				x: 0, y: 0,
			},
			end: point{
				x: 10, y: 0,
			},
		},
		line2: line{
			beg: point{
				x: 5, y: 5,
			},
			end: point{
				x: 5, y: 0,
			},
		},
		ans: false,
	},
	{
		line1: line{
			beg: point{
				x: -10, y: 0,
			},
			end: point{
				x: 10, y: 0,
			},
		},
		line2: line{
			beg: point{
				x: 0, y: 10,
			},
			end: point{
				x: 0, y: -10,
			},
		},
		ans: true,
	},
	{
		line1: line{
			beg: point{
				x: 10, y: 0,
			},
			end: point{
				x: 0, y: 0,
			},
		},
		line2: line{
			beg: point{
				x: 10, y: 0,
			},
			end: point{
				x: 20, y: 0,
			},
		},
		ans: false,
	},
}

func TestIntersection(t *testing.T) {
	for i, test := range intersectionTests {
		ans, x, y := intersection(test.line1, test.line2)
		if ans != test.ans {
			t.Log("TestIntersection failure", i)
			t.Errorf("\nexpected: %v, %v, %v\ngot: %v, %v, %v", test.ans, test.x, test.y, ans, x, y)
		}
	}
}
