package selector

import (
	"sync"

	"github.com/ManoloEsS/load_balancer/internal/backend"
)

type BackendSelector interface {
	Select() *backend.Server
}

type Selector struct {
	backends []backend.Server
	current  int
	mutex    sync.Mutex
	algoFunc func() *backend.Server
}

func NewSelector(algo string, servers []backend.Server) *Selector {
	slctor := &Selector{
		backends: servers,
		current:  0,
		mutex:    sync.Mutex{},
	}

	if algo == "round_robin" {
		slctor.algoFunc = slctor.roundRobin

	}

	return slctor
}

func (sl *Selector) Select() *backend.Server {
	return sl.algoFunc()
}

func (sl *Selector) roundRobin() *backend.Server {
	current := sl.current
	next := (sl.current + 1) % len(sl.backends)
	sl.current = next

	return &sl.backends[current]
}
