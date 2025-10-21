package store

import (
	"fmt"
	"slices"
	"sync"

	"github.com/Unfield/Valdock/namespaces"
)

type PortAllocator struct {
	mu            sync.Mutex
	store         Store
	minPort       int
	maxPort       int
	releasedPorts []int
}

func NewPortAllocator(store Store, minPort, maxPort int) (*PortAllocator, error) {
	var releasedPorts = []int{}
	_ = store.GetJSON(namespaces.INTERNAL_RELEASED_PORTS, &releasedPorts)

	pa := PortAllocator{
		store:         store,
		minPort:       minPort,
		maxPort:       maxPort,
		releasedPorts: releasedPorts,
	}

	return &pa, nil
}

func (pa *PortAllocator) GetPort() (int, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	lastPort, err := getLastPort(pa.store, pa.minPort, pa.maxPort)
	if err != nil {
		return 0, err
	}

	if len(pa.releasedPorts) > 0 {
		n := len(pa.releasedPorts)
		selectedPort := pa.releasedPorts[n-1]
		pa.releasedPorts = pa.releasedPorts[:n-1]
		_ = pa.store.SetJSON(namespaces.INTERNAL_RELEASED_PORTS, pa.releasedPorts)
		return selectedPort, nil
	}

	selectedPort := lastPort + 1
	if selectedPort > pa.maxPort {
		return 0, fmt.Errorf("maxPort reached")
	}

	if err := pa.store.SetInt(namespaces.INTERNAL_LAST_PORT, selectedPort); err != nil {
		return 0, fmt.Errorf("failed to save last used port")
	}

	return selectedPort, nil
}

func (pa *PortAllocator) FreePort(port int) error {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	lastPort, err := getLastPort(pa.store, pa.minPort, pa.maxPort)
	if err != nil {
		return err
	}

	if port > pa.maxPort {
		return fmt.Errorf("port too large")
	}
	if port < pa.minPort {
		return fmt.Errorf("port too small")
	}

	if slices.Contains(pa.releasedPorts, port) {
		return fmt.Errorf("port already released")
	}

	if port > lastPort {
		return fmt.Errorf("port is bigger than current port")
	}

	pa.releasedPorts = append(pa.releasedPorts, port)

	if err := pa.store.SetJSON(namespaces.INTERNAL_RELEASED_PORTS, pa.releasedPorts); err != nil {
		return fmt.Errorf("failed to persist released ports: %w", err)
	}

	return nil
}

func getLastPort(store Store, minPort, maxPort int) (int, error) {
	lastPort, err := store.GetInt(namespaces.INTERNAL_LAST_PORT)
	if err != nil {
		if err.Error() == "not found" {
			return minPort - 1, nil
		}
		return 0, fmt.Errorf("failed to fetch last port: %w", err)
	}

	if lastPort < minPort {
		return 0, fmt.Errorf("lastPort smaller than minPort (%d < %d)", lastPort, minPort)
	}

	if lastPort >= maxPort {
		return 0, fmt.Errorf("maxPort reached")
	}

	return lastPort, nil
}
