package backend

import (
	"encoding/binary"
	"io"

	"github.com/boltdb/bolt"
)

type Snapshot struct {
	bu      *bolt.Bucket
	horizon uint64
}

func (s *Snapshot) NewReader() (io.Reader, error) {
	limit := make([]byte, 8)
	binary.BigEndian.PutUint64(limit, s.horizon+1)
	return &snapshotReader{
		limit: limit,
		c:     s.bu.Cursor(),
	}, nil
}

func (s *Snapshot) Close() error {
	return s.bu.Tx().Commit()
}

type snapshotReader struct {
	limit []byte
	c     *bolt.Cursor
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
		// format?
		k, v := sr.c.Next()
		if k == nil {
			return readN, io.EOF
		}
		sr.remain = append(sr.remain, k...)
		sr.remain = append(sr.remain, v...)
	}
}
