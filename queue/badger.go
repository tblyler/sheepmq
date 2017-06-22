package queue

import (
	"os"
	"sync/atomic"

	"github.com/dgraph-io/badger/badger"
	"github.com/golang/protobuf/proto"
	"github.com/tblyler/sheepmq/lease"
	"github.com/tblyler/sheepmq/shepard"
)

// Badger uses Badger KV store as a queue backend
type Badger struct {
	opts      *Options
	bopts     *badger.Options
	kv        *badger.KV
	currentID uint64
	idconv    *idConverter
	leases    *lease.Manager
}

// NewBadger creates a new instance of a Badger-backed Queue
func NewBadger(opts *Options) (*Badger, error) {
	// use default badger options if none provided
	var bopts badger.Options
	if opts.BadgerOptions == nil {
		bopts = badger.DefaultOptions
	} else {
		bopts = *opts.BadgerOptions
	}

	// make sure the directory exists
	err := os.MkdirAll(opts.Dir, defaultFilePerm)
	if err != nil {
		return nil, err
	}

	// always honor Options' dir setting over badger.Options' dir settings
	bopts.Dir = opts.Dir
	bopts.ValueDir = opts.Dir

	// try to open new badger key value instance with the given options
	kv, err := badger.NewKV(&bopts)
	if err != nil {
		return nil, err
	}

	var currentID uint64

	iter := kv.NewIterator(badger.IteratorOptions{
		PrefetchSize: 5,
		FetchValues:  false,
		Reverse:      true,
	})

	defer iter.Close()

	for iter.Rewind(); iter.Valid(); iter.Next() {
		currentID, err = byteToID(iter.Item().Key())
		if err == nil {
			break
		}

		// try to delete invalid entries
		kv.Delete(iter.Item().Key())
		currentID = 0
	}

	return &Badger{
		opts:      opts,
		bopts:     &bopts,
		kv:        kv,
		currentID: currentID,
		idconv:    newIDConverter(),
		leases:    lease.NewManager(),
	}, nil
}

// Close the internal key value store
func (b *Badger) Close() error {
	return b.kv.Close()
}

// Get the next available ID atomically
func (b *Badger) getID() uint64 {
	return atomic.AddUint64(&b.currentID, 1)
}

// AddItem to the queue
func (b *Badger) AddItem(item *shepard.Item) error {
	item.Id = b.getID()

	data, err := proto.Marshal(item)
	if err != nil {
		return err
	}

	byteID := b.idconv.idToByte(item.Id)
	defer b.idconv.put(byteID)

	return b.kv.Set(byteID, data)
}

// GetItem from the queue
func (b *Badger) GetItem(info *shepard.GetInfo, itemChan chan<- *shepard.Item) error {
	iter := b.kv.NewIterator(badger.IteratorOptions{
		PrefetchSize: 500,
		FetchValues:  true,
		Reverse:      false,
	})

	defer iter.Close()

	var count uint64
	for iter.Rewind(); iter.Valid() && count < info.Count; iter.Next() {
		item := iter.Item()
		id, err := byteToID(item.Key())
		if err != nil {
			// try to delete bad keys (don't care about failures)
			b.kv.Delete(item.Key())
			continue
		}

		err = b.leases.AddLease(id, info)
		if err == nil || err == lease.ErrNoLeaser {
			ret := &shepard.Item{}
			err = proto.Unmarshal(item.Value(), ret)
			if err != nil {
				// delete bad values
				b.kv.Delete(item.Key())
				continue
			}

			count++
			itemChan <- ret
		}
	}

	return nil
}

// DelItem from the queue
func (b *Badger) DelItem(info *shepard.DelInfo) error {
	var err error
	for _, id := range info.GetIds() {
		idByte := b.idconv.idToByte(id)
		err = b.kv.Delete(idByte)
		if err != nil {
			return err
		}

		b.idconv.put(idByte)
	}

	return nil
}
