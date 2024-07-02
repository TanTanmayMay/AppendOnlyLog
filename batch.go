package aol

type BatchEntry struct {
	size int64
}

type Batch struct {
	entries []BatchEntry
	datas   []byte
}

func (b *Batch) Write(data []byte) {
	b.entries = append(b.entries, BatchEntry{int64(len(data))})
	b.datas = append(b.datas, data...)
}

func (b *Batch) Clear() {
	b.entries = b.entries[:0]
	b.datas = b.datas[:0]
}

// todo
func (l *Log) WriteBatch(b *Batch) error {

}

func (l *Log) writeBatch(b *Batch) error {

}
