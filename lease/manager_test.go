package lease

import (
	"testing"
	"time"

	"github.com/tblyler/sheepmq/shepard"
)

func TestManagerAddCheckLease(t *testing.T) {
	manager := NewManager()

	info := &shepard.GetInfo{}

	err := manager.AddLease(123, info)
	if err != ErrNoLeaser {
		t.Error("Failed to get ErrNoLeaser on add lease got", err)
	}

	info.PidLease = &shepard.PidLease{
		// use a bad PID
		Pid: 0,
	}

	err = manager.AddLease(123, info)
	if err != nil {
		t.Error("Failed to add a crappy PID leaser", err)
	}

	err = manager.AddLease(123, info)
	if err != nil {
		t.Error("Failed to add a crappy PID leaser ontop of another", err)
	}

	info.PidLease = nil

	info.HeartLease = &shepard.HeartbeatLease{
		Ttl: int64(time.Second / 2),
	}

	err = manager.AddLease(2, info)
	if err != nil {
		t.Error("Failed to add a valid heartbeat lease", err)
	}

	if !manager.CheckLease(2) {
		t.Error("Failed to check for valid heartbeat lease")
	}

	info.HeartLease = nil
	info.TimeoutLease = &shepard.TimeLease{
		Ttl: int64(time.Second),
	}

	err = manager.AddLease(123, info)
	if err != nil {
		t.Error("Failed to add a good timeout lease of 1 second")
	}

	err = manager.AddLease(123, info)
	if err == nil {
		t.Error("Should not be able to add a lease ontop of another valid one")
	}

	err = manager.AddLease(124, info)
	if err != nil {
		t.Error("Failed to add a valid lease against a different id", err)
	}

	if !manager.CheckLease(123) {
		t.Error("Failed to make sure valid lease was valid")
	}

	if !manager.CheckLease(124) {
		t.Error("Failed to make sure valid lease was valid")
	}

	time.Sleep(time.Second)

	if manager.CheckLease(123) {
		t.Error("failed to make sure invalid lease was invalid")
	}

	if manager.CheckLease(124) {
		t.Error("failed to make sure invalid lease was invalid")
	}
}

func TestManagerPruneLeases(t *testing.T) {
	manager := NewManager()

	info := &shepard.GetInfo{}
	info.PidLease = &shepard.PidLease{
		Pid: 0,
	}

	err := manager.AddLease(5, info)
	if err != nil {
		t.Error("Failed to add crappy PID lease", err)
	}

	info.PidLease = nil
	info.TimeoutLease = &shepard.TimeLease{
		Ttl: int64(time.Second),
	}

	err = manager.AddLease(2, info)
	if err != nil {
		t.Error("Failed to add valid timeout lease", err)
	}

	info.TimeoutLease = nil
	info.HeartLease = &shepard.HeartbeatLease{
		Ttl: int64(time.Second),
	}

	err = manager.AddLease(1, info)
	if err != nil {
		t.Error("Failed to add valid heartbeat lease", err)
	}

	manager.PruneLeases()

	if len(manager.leases) != 2 {
		t.Errorf("Should have 2 leases left, have %d", len(manager.leases))
	}

	time.Sleep(time.Second)

	manager.PruneLeases()
	if len(manager.leases) != 0 {
		t.Errorf("Should have 0 leases left, have %d", len(manager.leases))
	}
}
