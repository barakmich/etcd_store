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
	// ListKeysWithPrefix lists the keys with the given prefix.
	// The returned keys are sorted in (level, bytes) order. Each "/" after prefix increments the level of a returned key (to simulate directory depth).
	// level is the max level of the listed keys. Zero is unlimited.
	ListKeysWithPrefix(prefix []byte, horizon uint64, level int) ([][]byte, error)

	// ListKeysInRange lists the keys between [start, end).
	// The returned keys are sorted in (level, bytes) order. Each "/" after prefix increments the level of a returned key (to simulate directory depth).
	// level is the max level of the listed keys. Zero is unlimited.
	ListKeysInRange(start, end []byte, horizon uint64, level int) ([][]byte, error)

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
