package lease

import (
	"time"
)

// Timeout expires after a given timeout
type Timeout struct {
	eol time.Time
}

// NewTimeout creates a new timeout instance
func NewTimeout(ttl int64) *Timeout {
	return &Timeout{
		eol: time.Now().Add(time.Nanosecond * time.Duration(ttl)),
	}
}

// Valid Determines whether the timeout has been reached
func (t *Timeout) Valid() bool {
	return t.Check()
}

// Check if the timeout has been reached
func (t *Timeout) Check() bool {
	return time.Now().Before(t.eol)
}
