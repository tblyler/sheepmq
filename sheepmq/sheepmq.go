package sheepmq

import (
	"github.com/tblyler/sheepmq/queue"
	"github.com/tblyler/sheepmq/shepard"
)

// SheepMQ encapsulates the sheepmq queue logic
type SheepMQ struct {
	queues map[string]*queue.Queue
}

// NewSheepMQ creates a new sheepmq instance with the given configuration
func NewSheepMQ() (*SheepMQ, error) {
	return &SheepMQ{
		queues: make(map[string]*queue.Queue),
	}, nil
}

// AddItem to the sheepmq queue
func (l *SheepMQ) AddItem(item *shepard.Item) (*shepard.Response, error) {
	return nil, nil
}

// GetItems from sheepmq's queue
func (l *SheepMQ) GetItems(info *shepard.GetInfo, items chan<- *shepard.Item) error {
	return nil
}

// DelItem from sheepmq's queue
func (l *SheepMQ) DelItem(info *shepard.DelInfo) (*shepard.Response, error) {
	return nil, nil
}

// ErrItem from sheepmq's queue
func (l *SheepMQ) ErrItem(info *shepard.ErrInfo) (*shepard.Response, error) {
	return nil, nil
}
