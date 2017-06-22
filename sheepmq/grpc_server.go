package sheepmq

import (
	"io"

	"github.com/tblyler/sheepmq/shepard"

	context "golang.org/x/net/context"
)

// GServer encapsulates all sheepmq GRPC server activity
type GServer struct {
	sheepmq *SheepMQ
}

// NewGServer creates a new sheepmq GRPC server instance
func NewGServer(sheepmq *SheepMQ) *GServer {
	return &GServer{
		sheepmq: sheepmq,
	}
}

// AddItem to sheepmq's queue
func (l *GServer) AddItem(stream shepard.Sheepmq_AddItemServer) error {
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		resp, _ := l.sheepmq.AddItem(item)
		err = stream.Send(resp)
		if err != nil {
			return err
		}

		if !resp.GetSuccess() {
			return nil
		}
	}
}

// GetItem from sheepmq's queue
func (l *GServer) GetItem(info *shepard.GetInfo, stream shepard.Sheepmq_GetItemServer) error {
	items := make(chan *shepard.Item, 32)

	var err error
	go func() {
		err = l.sheepmq.GetItems(info, items)
		close(items)
	}()

	for item := range items {
		err := stream.Send(item)
		if err != nil {
			return err
		}
	}

	return err
}

// DelItem from sheepmq's queue
func (l *GServer) DelItem(ctx context.Context, info *shepard.DelInfo) (*shepard.Response, error) {
	return l.sheepmq.DelItem(info)
}

// ErrItem from sheepmq's queue
func (l *GServer) ErrItem(ctx context.Context, info *shepard.ErrInfo) (*shepard.Response, error) {
	return l.sheepmq.ErrItem(info)
}
