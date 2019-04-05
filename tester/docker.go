package tester

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"

	uuid "github.com/google/uuid"
	"github.com/pkg/errors"
)

type docker struct {
	port   int
	uuid   uuid.UUID
	adress string
}

func InitDocker(imgName string, port int, code string) (*docker, error) {
	newDocker := &docker{
		port:   port,
		adress: "127.0.0.1:" + strconv.Itoa(port),
	}
	var err error
	newDocker.uuid, err = uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create uuid")
	}

	fmt.Println("initing")
	startDocker := exec.Command("docker", "run", "-p", "5000:"+strconv.Itoa(port), "--name", newDocker.uuid.String(), "-t", imgName)
	err = startDocker.Start()
	if err != nil {
		return nil, errors.New(`failed to run "docker run"`)
	}

	defer func(cmd *exec.Cmd, d *docker) {
		//cmd.Process.Kill()
		if !newDocker.IsUp() {
			fmt.Println("defer closing")
			newDocker.Kill()
		}
	}(startDocker, newDocker)
	timer := time.NewTimer(time.Second * 2)
	ch := make(chan (error), 0)
	go func(cmd *exec.Cmd, ch chan<- (error)) {
		err := cmd.Wait()
		ch <- err
	}(startDocker, ch)

	select {
	case <-timer.C:
		// checking if docker container alive
		if !newDocker.IsUp() {
			startDocker.Process.Kill()
			newDocker.Kill()
			return nil, errors.New("time for docker launch is up")
		}
	case err, _ := <-ch:
		// docker cmd down. Usually error in start
		if err != nil {
			return nil, errors.Wrap(err, "failed to start docker")
		}
		startDocker.Process.Kill()
		newDocker.Kill()
		return nil, errors.New("docker suddenly stopped")

	}

	fmt.Println("initited")
	if err != nil {
		return nil, errors.Wrap(err, "failed to start docker")
	}

	return newDocker, nil
}

func (d *docker) MakeRequest(request []byte) ([]byte, error) {
	return request, nil
}

func (d *docker) Kill() {
	exec.Command("docker", "kill", d.uuid.String()).Run()
}

func KillDockers(ds ...*docker) {
	for _, d := range ds {
		d.Kill()
	}
}

func (d *docker) IsUp() bool {
	resp, err := exec.Command("docker", "ps", "-f", "name="+d.uuid.String()).Output()
	if err != nil {
		fmt.Println("isUp err:", err)
		return false
	}

	fmt.Println(string(resp), len(resp))
	// len("CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                              NAMES") == 145
	return len(resp) >= 145
}
