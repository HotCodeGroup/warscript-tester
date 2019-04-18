package tester

import (
	"time"

	"github.com/HotCodeGroup/warscript-tester/games"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type Tester struct {
	dockerClient *docker.Client
	ch           *amqp.Channel
	ports        *PortsPool
}

func NewTester(d *docker.Client, ch *amqp.Channel) *Tester {
	return &Tester{
		dockerClient: d,
		ch:           ch,
		ports:        newPortsPool(),
	}
}

// Test - tests bots submitted as RawCode1 and Rawcode2 by game rules
func (t *Tester) Test(RawCode1, Rawcode2 string, game games.Game) (states []games.State, result games.Result, returnErr error) {
	game.Init()
	im1, im2 := game.Images()

	port1 := t.ports.GetPort()
	defer t.ports.Free(port1)

	p1Container, err := NewPlayerContainer(1, port1, im1, 5*time.Second, t.dockerClient)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := p1Container.Remove()
		if err != nil {
			returnErr = errors.Wrap(err, "can not remove p1 container")
		}
	}()

	port2 := t.ports.GetPort()
	defer t.ports.Free(port2)

	p2Container, err := NewPlayerContainer(2, port2, im2, 5*time.Second, t.dockerClient)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		err := p2Container.Remove()
		if err != nil {
			returnErr = errors.Wrap(err, "can not remove p2 container")
		}
	}()

	//main game loop
	states = make([]games.State, 0, 0)
	// time.Sleep(2 * time.Second)
	// return states, &pong.Result{Winner: 1}, nil
	for {
		//log.Println("step")
		st1, st2 := game.Snapshots()

		resp1, err1 := p1Container.SendState(st1)
		if err1 != nil {
			returnErr = errors.Wrap(err1, "docker1 error")
		}
		resp2, err2 := p2Container.SendState(st2)
		if err2 != nil {
			returnErr = errors.Wrap(err2, "docker2 error")
		}

		gameErr := game.SaveSnapshots(resp1, resp2)
		if gameErr != nil {
			returnErr = errors.Wrap(gameErr, "save snapshot error")
		}

		state, fin := game.GetState()
		states = append(states, state)
		if fin {
			result = game.GetResult()
			return
		}
	}
}
