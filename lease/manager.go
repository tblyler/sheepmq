package lease

import (
	"errors"
	"sync"

	"github.com/tblyler/sheepmq/shepard"
)

var (
	// ErrLeased denotes the item is already leased
	ErrLeased = errors.New("Item is already leased")

	// ErrNoLeaser denotes no leaser was provided to add a lease
	ErrNoLeaser = errors.New("No leaser provided")
)

// Manager contains many leases and their validity
type Manager struct {
	leases map[uint64]Leaser
	locker sync.RWMutex
}

// NewManager creates a new Manager instance
func NewManager() *Manager {
	return &Manager{
		leases: make(map[uint64]Leaser),
	}
}

// AddLease to the manager for the given item
func (m *Manager) AddLease(id uint64, info *shepard.GetInfo) error {
	if m.CheckLease(id) {
		return ErrLeased
	}

	var leaser Leaser
	if info.TimeoutLease != nil {
		leaser = NewTimeout(info.TimeoutLease.Ttl)
	} else if info.PidLease != nil {
		leaser = NewPID(int(info.PidLease.Pid))
	} else if info.HeartLease != nil {
		leaser = NewHeart(info.HeartLease.Ttl)
	} else {
		return ErrNoLeaser
	}

	m.locker.Lock()
	m.leases[id] = leaser
	m.locker.Unlock()

	return nil
}

// CheckLease for validity
func (m *Manager) CheckLease(id uint64) bool {
	m.locker.RLock()
	lease, exists := m.leases[id]
	if !exists {
		m.locker.RUnlock()
		return false
	}

	ret := lease.Valid()
	m.locker.RUnlock()
	if !ret {
		// delete the lease since it is no longer valid
		m.locker.Lock()
		delete(m.leases, id)
		m.locker.Unlock()
	}

	return ret
}

// PruneLeases that fail their checks
func (m *Manager) PruneLeases() {
	deleteKeys := []uint64{}
	m.locker.RLock()

	for key, lease := range m.leases {
		if !lease.Check() {
			deleteKeys = append(deleteKeys, key)
		}
	}

	m.locker.RUnlock()
	m.locker.Lock()
	for _, key := range deleteKeys {
		delete(m.leases, key)
	}

	m.locker.Unlock()
}
