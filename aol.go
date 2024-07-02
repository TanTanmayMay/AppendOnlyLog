package aol

import (
	"errors"
	"os"
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
	NoSync: false,	// 'fsync' after every write
	SegmentSize: 20 * 1024 * 1024, // 20MB Segments
	DirectoryPermissions: 0750, // Owner Group Other (1-R 2-W 4-E) ->  || 750 => Owner(RWE), Group(RE), Others(NoPermission)
	FilePermissions: 0640,
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
