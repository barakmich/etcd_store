package backend

import (
	"encoding/binary"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type backend struct {
	db *leveldb.DB

	hbytes []byte
}

func New() (*backend, error) {
	// TODO: disable default compaction.
	// we never overwrite/delete a key.
	// defuault compaction is simply a waste of time.
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		return nil, err
	}
	return &backend{
		db:     db,
		hbytes: make([]byte, 8),
	}, nil
}

func (b *backend) Put(herizon uint64, kv []byte) error {
	// TODO: add prefix
	binary.BigEndian.PutUint64(b.hbytes, herizon)
	return b.db.Put(b.hbytes, kv, nil)
}

// Snapshot creates a snapshot at the given herizon.
func (b *backend) Snapshot(herizon uint64) (*Snapshot, error) {
	// force sync the db
	binary.BigEndian.PutUint64(b.hbytes, herizon)
	err := b.db.Put([]byte("snapshot"), b.hbytes, &opt.WriteOptions{true})
	if err != nil {
		return nil, err
	}

	// create the snapshot
	dbsnapshot, err := b.db.GetSnapshot()
	if err != nil {
		return nil, err
	}
	snapshot := &Snapshot{
		snapshot: dbsnapshot,
		herizon:  herizon,
	}
	return snapshot, nil
}

// PutCompact puts a compacted state at the given herizon into the db.
// db can release the herizon before the given herizon(including).
func (b *backend) PutCompact(state []byte, herizon uint64) error {
	// Save compacted state
	// background goroutine
	//     1. Batch Remove
	//     2. Trigger db Compaction
	return nil
}

func (b *backend) Close() error {
	if b.db == nil {
		return errors.New("backend: empty db")
	}
	return b.db.Close()
}
