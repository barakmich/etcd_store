package etcd_store

// call in a go-routine
func compact(be BackEnd, herizon uint64) {
	snap := be.GetSnapshot(herizon)
	reader := snap.NewReader()
	compacted := compact(reader)
	be.Compact(compacted, herizon)
}
