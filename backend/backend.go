package backend

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

var (
	kvBucket = []byte("kv")
)

type backend struct {
	db *bolt.DB

	hbytes []byte
}

func New(path string) (*backend, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &backend{
		db:     db,
		hbytes: make([]byte, 8),
	}, nil
}

func (be *backend) Put(horizon uint64, kv []byte) error {
	// TODO: add transaction coalescer
	binary.BigEndian.PutUint64(be.hbytes, horizon)
	err := be.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(kvBucket)
		if err != nil {
			return err
		}

		if err := b.Put(be.hbytes, kv); err != nil {
			return err
		}
		return nil
	})
	return err
}

// Snapshot creates a snapshot at the given herizon.
// TODO: maybe read out an compacted index.
func (be *backend) Snapshot(horizon uint64) (*Snapshot, error) {
	tx, err := be.db.Begin(false)
	if err != nil {
		return nil, err
	}
	b := tx.Bucket(kvBucket)
	if b == nil {
		return nil, errors.New("backend: empty db")
	}
	snapshot := &Snapshot{bu: b, horizon: horizon}
	return snapshot, nil
}

// Compact remove all given herizons。
// TODO: maybe save an index to the existing herizons for fast
// crash recovery
func (b *backend) Compact(herizons []uint64) error {
	batch := 10000 // do not hold io for a long time
	hbytes := make([]byte, 9)
	for i := 0; i < len(herizons); i = i + batch {
		err := b.db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists(kvBucket)
			if err != nil {
				return err
			}

			for j := i; j < i+batch; j++ {
				binary.BigEndian.PutUint64(hbytes, herizons[j])
				err := b.Delete(hbytes)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			return err
		}
		// sleep for a while before doing another
		// batch removal
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (b *backend) Close() error {
	if b.db == nil {
		return errors.New("backend: empty db")
	}
	return b.db.Close()
}
