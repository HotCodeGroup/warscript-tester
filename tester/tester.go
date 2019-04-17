package tester

import (
	"github.com/HotCodeGroup/warscript-tester/games"
	"github.com/pkg/errors"
)

var ports *PortsPool

func init() {
	ports = NewPortsPool()
}

// Test - tests bots submitted as RawCode1 and Rawcode2 by game rules
func Test(RawCode1, Rawcode2 string, game games.Game) (states []games.State, result games.Result, returnErr error) {
	game.Init()
	im1, im2 := game.Images()

	port1 := ports.GetPort()
	defer ports.Free(port1)

	port2 := ports.GetPort()
	defer ports.Free(port2)

	docker1, err := InitDocker(im1, port1, RawCode1)
	if err != nil {
		returnErr = errors.Wrap(err, "failed to init docker1")
		return
	}
	defer docker1.Kill()
	docker2, err := InitDocker(im2, port2, RawCode1)
	if err != nil {
		returnErr = errors.Wrap(err, "failed to init docker2")
		return
	}
	defer docker2.Kill()

	// main game loop
	states = make([]games.State, 0, 0)
	// time.Sleep(2 * time.Second)
	// return states, &pong.Result{Winner: 1}, nil
	for {
		st1, st2 := game.Snapshots()

		resp1, err1 := docker1.SendState(st1)
		if err1 != nil {
			returnErr = errors.Wrap(err1, "docker1 error")
		}
		resp2, err2 := docker2.SendState(st2)
		if err2 != nil {
			returnErr = errors.Wrap(err1, "docker2 error")
		}

		gameErr := game.SaveSnapshots(resp1, resp2)
		if gameErr != nil {
			returnErr = errors.Wrap(err1, "docker1 error")
		}

		state, fin := game.GetState()
		states = append(states, state)
		if fin {
			result = game.GetResult()
			return
		}
	}
}
