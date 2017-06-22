package queue

import (
	"encoding/binary"
	"fmt"
	"os"
	"sync"

	"github.com/dgraph-io/badger/badger"
	"github.com/tblyler/sheepmq/shepard"
)

const (
	idByteSize      = 8
	defaultFilePerm = os.FileMode(0700)
)

// Queue defines a resource to store queue items
type Queue interface {
	AddItem(*shepard.Item) error
	GetItem(*shepard.GetInfo) (*shepard.Item, error)
	DelItem(*shepard.DelInfo) error
}

// Options to be used when creating a new Queue
type Options struct {
	// Directory to store queue data
	Dir string

	// Badger queue specific options
	BadgerOptions *badger.Options
}

type idConverter struct {
	pool sync.Pool
}

func newIDConverter() *idConverter {
	return &idConverter{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]byte, idByteSize)
			},
		},
	}
}

func (i *idConverter) idToByte(id uint64) []byte {
	buf := i.pool.Get().([]byte)

	binary.LittleEndian.PutUint64(buf, id)

	return buf
}

func (i *idConverter) put(data []byte) {
	i.pool.Put(data)
}

func byteToID(data []byte) (uint64, error) {
	if len(data) < idByteSize {
		return 0, fmt.Errorf(
			"unable to convert byte slice length of %d, need at least %d",
			len(data),
			idByteSize,
		)
	}

	return binary.LittleEndian.Uint64(data), nil
}
