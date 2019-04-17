package tester

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	uuid "github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	sendCodeExt  = "/loadcode"
	sendStateExt = "/run"
)

type docker struct {
	port            int
	uuid            uuid.UUID
	adress          string
	sendCodeAdress  string
	sendStateAdress string
}

func InitDocker(imgName string, port int, code string) (*docker, error) {
	baseAdress := "http://127.0.0.1:" + strconv.Itoa(port)
	newDocker := &docker{
		port:            port,
		adress:          baseAdress,
		sendCodeAdress:  baseAdress + sendCodeExt,
		sendStateAdress: baseAdress + sendStateExt,
	}

	var err error
	newDocker.uuid, err = uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create uuid")
	}

	fmt.Println("initing")
	startDocker := exec.Command("docker", "run", "-p", strconv.Itoa(port)+":5000", "--name", newDocker.uuid.String(), "-t", imgName)
	err = startDocker.Start()
	if err != nil {
		return nil, errors.New(`failed to run "docker run"`)
	}

	ok := false
	defer func(cmd *exec.Cmd, d *docker) {
		//cmd.Process.Kill()
		if !ok {
			fmt.Println("defer closing")
			newDocker.Kill()
			startDocker.Process.Kill()
		}
	}(startDocker, newDocker)
	timer := time.NewTimer(time.Second * 5)
	ch := make(chan (error), 0)
	go func(cmd *exec.Cmd, ch chan<- (error)) {
		err := cmd.Wait()
		ch <- err
	}(startDocker, ch)

	select {
	case <-timer.C:
		// checking if docker container alive
		if !newDocker.IsUp() {
			if err := startDocker.Process.Kill(); err != nil {
				return nil, errors.Wrap(err, "time for docker launch is up: can not kill process")
			}
			if err := newDocker.Kill(); err != nil {
				return nil, errors.Wrap(err, "time for docker launch is up: can not stop newDocker")
			}

			return nil, errors.New("time for docker launch is up")
		}
	case err, _ := <-ch:
		// docker cmd down. Usually error in start
		if err != nil {
			return nil, errors.Wrap(err, "failed to start docker")
		}

		if err := startDocker.Process.Kill(); err != nil {
			return nil, errors.Wrap(err, "docker suddenly stopped: can not kill process")
		}
		if err := newDocker.Kill(); err != nil {
			return nil, errors.Wrap(err, "docker suddenly stopped: can not stop newDocker")
		}

		return nil, errors.New("docker suddenly stopped")
	}
	fmt.Println("initited")

	resp, err := newDocker.SendCode(code)
	fmt.Println("docker response:", string(resp))
	if err != nil {
		return nil, errors.Wrap(err, "failed to send code")
	}

	ok = true
	return newDocker, nil
}

func (d *docker) SendCode(code string) ([]byte, error) {
	return d.SendRequest([]byte(`{"code":"`+code+`"}`), d.sendCodeAdress)
}

func (d *docker) SendState(state []byte) ([]byte, error) {
	return d.SendRequest(state, d.sendStateAdress)
}

func (d *docker) SendRequest(request []byte, url string) ([]byte, error) {
	if !d.IsUp() {
		return nil, errors.New("docker is down")
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(request))
	if err != nil {
		return nil, errors.Wrap(err, "failed to send request")
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return res, nil
}

func (d *docker) Kill() error {
	return exec.Command("docker", "kill", d.uuid.String()).Run()
}

func KillDockers(ds ...*docker) error {
	for i, d := range ds {
		if err := d.Kill(); err != nil {
			return errors.Wrapf(err, "[%d] ", i)
		}
	}

	return nil
}

func (d *docker) IsUp() bool {
	resp, err := exec.Command("docker", "ps", "-f", "name="+d.uuid.String()).Output()
	if err != nil {
		fmt.Println("isUp err:", err)
		return false
	}

	//fmt.Println(string(resp), len(resp))
	// len("CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                              NAMES") == 145
	return len(resp) >= 145
}
