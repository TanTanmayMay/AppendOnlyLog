package aol

import (
	"errors"
	"os"
	"sync"
)

var (
	ErrCorrupt  = errors.New("log corrupt")
	ErrClosed   = errors.New("log closed")
	ErrNotFound = errors.New("log not found")
	ErrEOF      = errors.New("end of file reached")
)

type Options struct {
	// NoSync field type is used to give user ability to call "fsync" system call that flushes data to disk, ensuring that the data is written to permanent storage.
	// Disabling fsync can improve performance but increases the risk of data loss in case of a server crash because data might not be persisted to disk immediately.
	NoSync               bool
	SegmentSize          int
	DirectoryPermissions os.FileMode
	FilePermissions      os.FileMode
}

var DefaultOptions = &Options{
	NoSync:               false,            // 'fsync' after every write
	SegmentSize:          20 * 1024 * 1024, // 20MB Segments
	DirectoryPermissions: 0750,             // Owner Group Other (1-R 2-W 4-E) ->  || 750 => Owner(RWE), Group(RE), Others(NoPermission)
	FilePermissions:      0640,
}

func (o *Options) validate() {
	if o.SegmentSize <= 0 {
		o.SegmentSize = DefaultOptions.SegmentSize
	}

	if o.DirectoryPermissions == 0 {
		o.DirectoryPermissions = DefaultOptions.DirectoryPermissions
	}

	if o.FilePermissions == 0 {
		o.FilePermissions = DefaultOptions.FilePermissions
	}
}

// Log struct that represents an Append Only Log
type Log struct {
	mtx              sync.RWMutex // For Multiple readers and single writer problem
	path             string       // Absolute path of the local directory where log resides
	segments         []*segment   // A slice of all of the known Log Segments
	latestFileHandle *os.File     // The file handle for the most recent (tail) segment, used for writing new log entries.
	rytBatch         Batch        // Reusable Write Batch
	opts             Options
	closed           bool
	corrupt          bool
}

// Each segment is a part of the log, allowing the log to manage data in smaller, manageable files.
type segment struct {
	path               string // Absolute path to Segment file	|| 	This is FUN...I'm getting Tanmayed creating this project
	index              uint64 // Field to store first index of segment -> helps in identifying and ordering segments within the log.
	cacheEntryBuffer   []byte // A buffer that caches entries for the segment
	cacheEntryPosition []bpos // A slice of bpos structs that store the positions of entries within the cbuf. This allows for efficient lookup and retrieval of entries from the buffer.
}

type bpos struct {
	pos int // The byte position where an entry starts in the buffer.
	end int // The byte position just past the end of the entry in the buffer.
}
