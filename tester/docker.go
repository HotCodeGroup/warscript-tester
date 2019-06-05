package tester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	sendCodeEndpoint  = "loadcode"
	sendStateEndpoint = "run"
)

var (
	// ErrTimeount ошибка -- выход по таймауту
	ErrTimeount = errors.New("timeout")
)

// PlayerContainer инкапсуляция коннекта с докером
type PlayerContainer struct {
	PlayerID     int
	Port         docker.Port
	container    *docker.Container
	dockerClient *docker.Client

	httpClient *http.Client
}

// NewPlayerContainer создание и запуск нового докер контейнера
func NewPlayerContainer(playerID, port int, imageName string,
	timeout time.Duration, dockerClient *docker.Client) (*PlayerContainer, error) {
	pContainerUUID := uuid.New()
	pContainer, err := dockerClient.CreateContainer(docker.CreateContainerOptions{
		Name: fmt.Sprintf("p%d_%s_%s", playerID, imageName, pContainerUUID.String()),
		Config: &docker.Config{
			Image: imageName,
			ExposedPorts: map[docker.Port]struct{}{
				docker.Port("5000/tcp"): {},
			},
		},
		HostConfig: &docker.HostConfig{
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port("5000/tcp"): {
					docker.PortBinding{
						HostIP:   "127.0.0.1",
						HostPort: strconv.Itoa(port),
					},
				},
			},
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "can not create p%d container", playerID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = dockerClient.StartContainerWithContext(pContainer.ID, nil, ctx)
	if err != nil {
		if rmErr := dockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID:    pContainer.ID,
			Force: true,
		}); rmErr != nil {
			return nil, errors.Wrapf(rmErr, "can not remove p%d container(intial err %s)", playerID, err)
		}

		if ctx.Err() != nil { // упали по таймауту
			return nil, errors.Wrapf(ErrTimeount, "p%d container creation timeout", playerID)
		}

		return nil, errors.Wrapf(err, "can not start p%d container", playerID)
	}

	pContainer, err = dockerClient.InspectContainer(pContainer.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "can not inspect p%d container", playerID)
	}

	return &PlayerContainer{
		PlayerID:     playerID,
		Port:         docker.Port("5000/tcp"),
		container:    pContainer,
		dockerClient: dockerClient,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// SendCode инициализация докер контейнера кодом решения
func (p *PlayerContainer) SendCode(code string) ([]byte, error) {
	body, err := json.Marshal(struct {
		Code string `json:"code"`
	}{
		Code: code,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can not marshal body")
	}

	return p.SendRequest(body, sendCodeEndpoint)
}

// SendState отправка в контейнер текущего состояния игры
func (p *PlayerContainer) SendState(state []byte) ([]byte, error) {
	return p.SendRequest(state, sendStateEndpoint)
}

// SendRequest отправка запроса в контенер
func (p *PlayerContainer) SendRequest(body []byte, endpoint string) ([]byte, error) {
	portBinding := p.container.HostConfig.PortBindings[p.Port][0]
	resp, err := http.Post(fmt.Sprintf("http://%s:%s/%s", portBinding.HostIP,
		portBinding.HostPort, endpoint), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to send request to p%d docker", p.PlayerID)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read p%d docker response body", p.PlayerID)
	}

	return res, nil
}

// Remove остановка и удаление контейнера
func (p *PlayerContainer) Remove() error {
	err := p.dockerClient.StopContainer(p.container.ID, 1)
	if err != nil { // может быть ситуация при которой контейнер уже стопнут
		log.Printf("can not stop p%d container: %s", p.PlayerID, err)
	}

	err = p.dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:    p.container.ID,
		Force: true,
	})
	if err != nil {
		return errors.Wrapf(err, "can not remove p%d container", p.PlayerID)
	}

	return nil
}
