package atod

import (
	"math"
	"testing"
)

func TestMoveCircle(t *testing.T) {
	corr := float64(0)
	ans := float64(0)
	corr = float64(5)
	ans = moveCircle(0, 0, 5, 10, 5, -5, 25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 1)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}
	corr = float64(-25)
	ans = moveCircle(0, 0, 5, 10, 5, -5, -25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 2)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}
	corr = float64(-5)
	ans = moveCircle(0, 0, 5, -10, 5, -5, -25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 3)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}

	corr = float64(3)
	ans = moveCircle(0, 0, 5, 10, 5, -5, 3)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 4)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}

	corr = float64(5)
	ans = moveCircle(0, 5, 5, 10, 5, -5, 25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 5)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}

	corr = float64(6.46393226874)
	ans = moveCircle(0, 5+3.535, 5, 10, 5, -5, 25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 6)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}
	corr = float64(6.46393226874)
	ans = moveCircle(0, -5-3.535, 5, 10, 5, -5, 25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 7)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}

	corr = float64(-6.46393226874)
	ans = moveCircle(0, 5+3.535, 5, -10, -5, 5, -25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 8)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}
	corr = float64(-6.46393226874)
	ans = moveCircle(0, -5-3.535, 5, -10, 5, -5, -25)
	if math.Abs(corr-ans) > lEPS {
		t.Log("TestIntersection failure", 9)
		t.Errorf("\nexpected: %v\ngot: %v", corr, ans)
	}
}
