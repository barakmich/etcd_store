package backend

import (
	"encoding/binary"
	"io"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Snapshot struct {
	snapshot *leveldb.Snapshot
	herizon  uint64
}

func (s *Snapshot) NewReader() (io.Reader, error) {
	limit := make([]byte, 8)
	binary.BigEndian.PutUint64(limit, s.herizon+1)
	return &snapshotReader{
		it: s.snapshot.NewIterator(&util.Range{Limit: limit}, nil),
	}, nil
}

func (s *Snapshot) Close() error {
	s.snapshot.Release()
	return nil
}

type snapshotReader struct {
	it        iterator.Iterator
	exhausted bool
	// TODO: resue buffer
	remain []byte
}

func (sr *snapshotReader) Read(b []byte) (int, error) {
	// TODO: read compacted state
	wantN := len(b)
	readN := 0
	for {
		if len(sr.remain) != 0 {
			n := copy(b[readN:], sr.remain)
			sr.remain = sr.remain[n:]
			readN += n
			if readN == wantN {
				return wantN, nil
			}
		}
		if len(sr.remain) != 0 {
			panic("expect len(remain) == 0")
		}
		if sr.exhausted {
			sr.it.Release()
			return readN, io.EOF
		}
		// format?
		sr.remain = append(sr.remain, sr.it.Key()...)
		sr.remain = append(sr.remain, sr.it.Value()...)
		sr.exhausted = sr.it.Next()
	}
}
