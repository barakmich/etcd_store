package etcd_store

import (
	"errors"
	"io"
)

var (
	ErrHorizonOutOfDate  = errors.New("store: Current horizon is beyond requested horizon")
	ErrRevisionOutOfDate = errors.New("store: Current revision is beyond requested revision")
)

// Horizon/Revision 0 == most recent
type KV struct {
	Key      []byte
	Value    []byte
	Horizon  uint64
	Revision uint64
}

type Store interface {
	Version() int
	Horzon() uint64
	Get(kv KV) (KV, error)
	Set(kv KV) error
	Delete(kv KV) error
	ListKeysWithPrefix(key []byte, horizon uint64) ([][]byte, error)

	Watch(prefix []byte, recursive, stream bool, start, end uint64) (Watcher, error)

	SnapshotTo(w io.Writer) error
	Compact(horizon uint64) error
}

// modified from etcd/store/watcher.go. It seems reasonable enough.
type Watcher interface {
	KVChan() chan *KV
	StartHorizon() uint64 // The horizon at which the Watcher was created
	Close()
}
