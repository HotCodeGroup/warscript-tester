package tester

import (
	"sync"
)

type PortsPool struct {
	mu    *sync.Mutex
	ports map[int]bool
}

func newPortsPool() *PortsPool {
	return &PortsPool{
		mu: &sync.Mutex{},
		ports: map[int]bool{
			5001: true,
			5002: true,
			5003: true,
			5004: true,
			5005: true,
			5006: true,
			5007: true,
			5008: true,
			5009: true,
			5010: true,
			5011: true,
			5012: true,
			5013: true,
			5014: true,
			5015: true,
			5016: true,
		},
	}
}

func (a *PortsPool) GetPort() int {
	for {
		a.mu.Lock()
		for port, isFree := range a.ports {
			if isFree {
				a.ports[port] = false
				a.mu.Unlock()
				return port
			}
		}
		a.mu.Unlock()
	}
}

func (a *PortsPool) Free(port int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.ports[port]; !ok {
		return
	}

	a.ports[port] = true
}
