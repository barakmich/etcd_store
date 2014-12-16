package backend

import (
	"fmt"
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

func BenchmarkPut(b *testing.B) {
	toPut := make([][]byte, 1000000)
	for i := range toPut {
		toPut[i] = []byte(fmt.Sprintf("key_%d", i))
	}
	be, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer be.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := be.Put(uint64(i), toPut[i])
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	os.RemoveAll("db")
}

func BenchmarkBatchDelete(b *testing.B) {
	toPut := make([][]byte, 1000000)
	for i := range toPut {
		toPut[i] = []byte(fmt.Sprintf("key_%d", i))
	}
	be, err := New()
	if err != nil {
		b.Fatal(err)
	}
	defer be.Close()
	for i := 0; i < b.N; i++ {
		err := be.db.Put(toPut[i], toPut[i], nil)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.ResetTimer()
	batch := new(leveldb.Batch)
	for i := 0; i < b.N; i++ {
		batch.Delete(toPut[i])
	}
	err = be.db.Write(batch, nil)
	if err != nil {
		b.Fatal(err)
	}
	os.RemoveAll("db")
}
