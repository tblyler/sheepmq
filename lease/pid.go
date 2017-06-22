package lease

import (
	"os"
	"syscall"
)

// PID contains a lease that is valid until the given PID no longer exists
type PID struct {
	pid int
}

// NewPID creates a new PID leaser instance
func NewPID(pid int) *PID {
	return &PID{
		pid: pid,
	}
}

// Valid checks if the PID still exists
func (p *PID) Valid() bool {
	return p.Check()
}

// Check if the PID still exists
func (p *PID) Check() bool {
	process, err := os.FindProcess(p.pid)
	if err != nil {
		return false
	}

	// if nil, PID exists
	return process.Signal(syscall.Signal(0)) == nil
}
