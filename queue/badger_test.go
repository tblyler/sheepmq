package queue

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/dgraph-io/badger/badger"
	"github.com/tblyler/sheepmq/shepard"
)

func TestNewBadger(t *testing.T) {
	// test a bad dir setting
	opts := &Options{
		Dir: "",
	}

	_, err := NewBadger(opts)
	if err == nil {
		t.Error("Failed to get error on bad Dir setting")
	}

	opts.Dir = filepath.Join(os.TempDir(), "NewBadgerTest")

	// ensure clean directory
	os.RemoveAll(opts.Dir)
	defer os.RemoveAll(opts.Dir)

	b, err := NewBadger(opts)
	if err != nil {
		t.Fatal("Failed to create badger with default options")
	}

	if b.currentID != 0 {
		t.Error("Current ID of an empty badger db should be 0 not", b.currentID)
	}

	err = b.AddItem(&shepard.Item{})
	if err != nil {
		t.Fatal("Failed to add empty item", err)
	}

	if b.currentID != 1 {
		t.Error("Current ID of a one item db should be 1 not", b.currentID)
	}

	b.Close()

	opts.BadgerOptions = &badger.Options{}

	// make sure custom badger options are honored
	*opts.BadgerOptions = badger.DefaultOptions
	opts.BadgerOptions.SyncWrites = !opts.BadgerOptions.SyncWrites

	b, err = NewBadger(opts)
	if err != nil {
		t.Fatal("Failed to use custom badger optoins", err)
	}

	defer b.Close()

	if b.bopts.SyncWrites != opts.BadgerOptions.SyncWrites {
		t.Error("Failed to use custom badger options")
	}

	if b.currentID != 1 {
		t.Error("current id != 1 got", b.currentID)
	}
}

func TestBadgerAddGetItem(t *testing.T) {
	items := make([]*shepard.Item, 32)

	for i := range items {
		items[i] = &shepard.Item{}

		items[i].Data = make([]byte, 256*(i+1))
		rand.Read(items[i].Data)
		items[i].Ctime = time.Now().Unix()
		items[i].Queue = fmt.Sprint("testing", i)
		items[i].Stats = map[string]int64{
			"cool":    133333337,
			"notcool": 0,
			"#1":      1,
			"datSize": int64(len(items[i].Data)),
		}
	}

	opts := &Options{
		Dir: filepath.Join(os.TempDir(), "TestBadgerAddGetItem"),
	}

	os.RemoveAll(opts.Dir)
	defer os.RemoveAll(opts.Dir)

	b, err := NewBadger(opts)
	if err != nil {
		t.Fatal("Failed to open badger", err)
	}

	defer b.Close()

	for i, item := range items {
		err = b.AddItem(item)
		if err != nil {
			t.Error("Failed to add item", i, err)
		}
	}

	itemChan := make(chan *shepard.Item, len(items))
	err = b.GetItem(&shepard.GetInfo{
		Count: uint64(len(items)),
	}, itemChan)

	close(itemChan)
	if err != nil {
		t.Error("Failed to get items", err)
	}

	i := 0
	for item := range itemChan {
		if item.Ctime != items[i].Ctime {
			t.Error("item ctimes", item.Ctime, items[i].Ctime)
		}

		if item.Queue != items[i].Queue {
			t.Error("item queues", item.Queue, items[i].Queue)
		}

		if !bytes.Equal(item.Data, items[i].Data) {
			t.Error("item data", item.Data, items[i].Data)
		}

		if !reflect.DeepEqual(item.Stats, items[i].Stats) {
			t.Error("item stats", item.Stats, items[i].Stats)
		}

		i++
	}

	if i != len(items) {
		t.Error("got", i, "items expected", len(items))
	}
}

func TestBadgerDelItem(t *testing.T) {
	opts := &Options{
		Dir: filepath.Join(os.TempDir(), "TestBadgerAddGetItem"),
	}

	os.RemoveAll(opts.Dir)
	defer os.RemoveAll(opts.Dir)

	b, err := NewBadger(opts)
	if err != nil {
		t.Fatal("Failed to start badger", err)
	}

	items := make([]*shepard.Item, 32)
	for i := range items {
		items[i] = &shepard.Item{
			Ctime: time.Now().Unix(),
			Data:  make([]byte, 256*(i+1)),
			Queue: "The queue",
			Stats: map[string]int64{
				"lol":       10101010101,
				"index":     int64(i),
				"datasize:": int64(256 * (i + 1)),
			},
		}
		rand.Read(items[i].Data)

		err = b.AddItem(items[i])
		if err != nil {
			t.Error("Failed to add item", i, err)
		}
	}

	delinfo := &shepard.DelInfo{}
	for i := range items {
		if i%2 == 0 {
			continue
		}

		delinfo.Ids = append(delinfo.Ids, uint64(i))
	}

	err = b.DelItem(delinfo)
	if err != nil {
		t.Error("Failed to delete", delinfo.Ids, err)
	}

	getinfo := &shepard.GetInfo{
		Count: uint64(len(items) - len(delinfo.Ids)),
	}
	getChan := make(chan *shepard.Item, getinfo.Count)

	err = b.GetItem(getinfo, getChan)
	if err != nil {
		t.Error("Failed to get items", err)
	}

	close(getChan)

	i := 1
	for item := range getChan {
		if item.Ctime != items[i].Ctime {
			t.Error("item ctimes", item.Ctime, items[i].Ctime)
		}

		if item.Queue != items[i].Queue {
			t.Error("item queues", item.Queue, items[i].Queue)
		}

		if !bytes.Equal(item.Data, items[i].Data) {
			t.Error("item data", item.Data, items[i].Data)
		}

		if !reflect.DeepEqual(item.Stats, items[i].Stats) {
			t.Error("item stats", item.Stats, items[i].Stats)
		}
		i += 2
	}

	if i != len(items)+1 {
		t.Error("only looped to item index", i)
	}
}
