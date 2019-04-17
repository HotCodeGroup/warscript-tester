package tester

import "sync"

type PortsPool struct {
	mu    *sync.Mutex
	ports map[int]bool
}

func NewPortsPool() *PortsPool {
	return &PortsPool{
		mu: &sync.Mutex{},
		ports: map[int]bool{
			5001: false,
			5002: false,
			5003: false,
			5004: false,
			5005: false,
			5006: false,
			5007: false,
			5008: false,
			5009: false,
			5010: false,
			5011: false,
			5012: false,
			5013: false,
			5014: false,
			5015: false,
			5016: false,
		},
	}
}

func (a *PortsPool) GetPort() int {
	for {
		for port, isFree := range a.ports {
			if isFree {
				a.mu.Lock()
				a.ports[port] = true
				a.mu.Unlock()
			}

			return port
		}
	}
}

func (a *PortsPool) Free(port int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.ports[port]; !ok {
		return
	}
	a.ports[port] = false
}
