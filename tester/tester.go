package tester

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/HotCodeGroup/warscript-tester/games"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

// Tester структура инапсулирующая тестер
type Tester struct {
	dockerClient *docker.Client
	ch           *amqp.Channel
	ports        *PortsPool
}

type SendCodeResponce struct {
	Error string `json:"error"`
}

// NewTester создание нового объекта тестировщика
func NewTester(d *docker.Client, ch *amqp.Channel) *Tester {
	return &Tester{
		dockerClient: d,
		ch:           ch,
		ports:        newPortsPool(),
	}
}

// Test - tests bots submitted as RawCode1 and Rawcode2 by game rules
func (t *Tester) Test(rawCode1, rawCode2 string, game games.Game) (info games.Info,
	states []games.State, result games.Result, returnErr error) {
	im1, im2 := game.Images()

	port1 := t.ports.GetPort()
	defer t.ports.Free(port1)

	port2 := t.ports.GetPort()
	defer t.ports.Free(port2)

	port1, _ = strconv.Atoi(os.Getenv("PORT_1"))
	port2, _ = strconv.Atoi(os.Getenv("PORT_2"))

	p1Container, err := NewPlayerContainer(1, port1, im1, 60*time.Second, t.dockerClient)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		err = p1Container.Remove()
		if err != nil {
			returnErr = errors.Wrap(err, "can not remove p1 container")
		}
	}()

	p2Container, err := NewPlayerContainer(2, port2, im2, 60*time.Second, t.dockerClient)
	if err != nil {
		return nil, nil, nil, err
	}
	defer func() {
		err = p2Container.Remove()
		if err != nil {
			returnErr = errors.Wrap(err, "can not remove p2 container")
		}
	}()

	time.Sleep(1 * time.Second)
	// main game loop
	game.Init()

	info = game.GetInfo()
	states = make([]games.State, 0, 0)
	respSendCode1, err := p1Container.SendCode(rawCode1)
	if err != nil {
		returnErr = errors.Wrap(err, "can not init p1 container code")
		return
	}

	respSendCodeStatus1 := &SendCodeResponce{}
	err = json.Unmarshal(respSendCode1, respSendCodeStatus1)
	if err != nil {
		returnErr = errors.Wrap(err, "can not unmarshal p1 container resp")
		return
	}

	if respSendCodeStatus1.Error != "" {
		returnErr = errors.New("player1 load code error: " + respSendCodeStatus1.Error)
		return
	}

	respSendCode2, err := p2Container.SendCode(rawCode2)
	if err != nil {
		returnErr = errors.Wrap(err, "can not init p2 container code")
		return
	}

	respSendCodeStatus2 := &SendCodeResponce{}
	err = json.Unmarshal(respSendCode2, respSendCodeStatus2)
	if err != nil {
		returnErr = errors.Wrap(err, "can not unmarshal p2 container resp")
		return
	}

	if respSendCodeStatus2.Error != "" {
		returnErr = errors.New("player2 load code error: " + respSendCodeStatus2.Error)
		return
	}

	for {
		st1, st2 := game.Snapshots()
		resp1, err1 := p1Container.SendState(st1)
		if err1 != nil {
			returnErr = errors.Wrap(err1, "docker1 error")
			return
		}
		resp2, err2 := p2Container.SendState(st2)
		if err2 != nil {
			returnErr = errors.Wrap(err2, "docker2 error")
			return
		}

		gameErr := game.SaveSnapshots(resp1, resp2)
		if gameErr != nil {
			returnErr = errors.Wrap(gameErr, "save snapshot error")
			return
		}

		state, fin := game.GetState()
		states = append(states, state)
		if fin {
			result = game.GetResult()
			return
		}
	}
}
