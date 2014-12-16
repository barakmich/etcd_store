package etcd_store

import (
	"bytes"
	"errors"

	"github.com/google/btree"
	"github.com/xiangli-cmu/etcd_store/backend"
)

type store struct {
	horizon uint64
	bt      *btree.BTree

	be backend.Backend
}

func (s *store) Put(key []byte, val []byte) error {
	nexth := s.horizon + 1
	err := s.be.Put(nexth, append(key, val...))
	if err != nil {
		return err
	}

	s.horizon = nexth
	keyi := keyIndex{key: key}

	item := s.bt.Get(keyi)
	if item == nil {
		keyi.val = &content{rev: 1, horizon: s.horizon, next: nil}
		s.bt.ReplaceOrInsert(item)
		return nil
	}
	// add to head
	okeyi := item.(keyIndex)
	oldv := okeyi.val
	okeyi.val = &content{rev: oldv.rev + 1, horizon: s.horizon, next: oldv}
	return nil
}

func (s *store) Get(horizon uint64, key []byte) ([]byte, error) {
	keyi := keyIndex{key: key}
	item := s.bt.Get(keyi)
	if item == nil {
		return nil, errors.New("key not found")
	}

	keyi = item.(keyIndex)
	val := keyi.val
	for val != nil {
		if horizon > val.horizon {
			kv, err := s.be.Get(val.horizon)
			if err != nil {
				return nil, err
			}
			return kv[len(key):], nil
		}
		val = val.next
	}
	return nil, errors.New("key not found")
}

type keyIndex struct {
	key []byte
	val *content
}

func (a keyIndex) Less(b btree.Item) bool {
	return bytes.Compare(a.key, b.(keyIndex).key) == -1
}

type content struct {
	rev     uint64
	horizon uint64
	next    *content
}
