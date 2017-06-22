package lease

import "time"

// Heart contains a lease that survives when heartbeats happen within a given ttl interval
type Heart struct {
	lastBeat time.Time
	ttl      time.Duration
}

// NewHeart creates a new heart instance
func NewHeart(ttl int64) *Heart {
	ttlDuration := time.Nanosecond * time.Duration(ttl)
	return &Heart{
		lastBeat: time.Now(),
		ttl:      ttlDuration,
	}
}

// Valid tries to beat the heart, true if still alive, false if dead
func (h *Heart) Valid() bool {
	if !h.Check() {
		return false
	}

	// update the last beat to now
	h.lastBeat = time.Now()
	return true
}

// Check if the heartbeat is valid
func (h *Heart) Check() bool {
	return time.Since(h.lastBeat) <= h.ttl
}
