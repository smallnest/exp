package io

import (
	"errors"
	"io"
	"unsafe"

	"golang.org/x/exp/mmap"
)

var _ io.Reader = (*MmapReader)(nil)
var _ io.Closer = (*MmapReader)(nil)
var _ io.ReaderAt = (*MmapReader)(nil)

type readerAt struct {
	data []byte
}

// MmapReader is a mmap.ReaderAt wrapper that implements io.Reader.
type MmapReader struct {
	*mmap.ReaderAt
	off int64
}

// NewMmapReader returns a new Reader that reads from the file named filename.
func NewMmapReader(filename string) (*MmapReader, error) {
	readerAt, err := mmap.Open(filename)
	if err != nil {
		return nil, err
	}

	r := &MmapReader{
		ReaderAt: readerAt,
	}

	return r, nil
}

// Close closes the Reader.
func (r *MmapReader) Read(p []byte) (n int, err error) {
	n, err = r.ReaderAt.ReadAt(p, r.off)
	r.off += int64(n)
	return
}

// ReadAll reads all data from the Reader.
func (r *MmapReader) ReadAll() ([]byte, error) {
	l := r.ReaderAt.Len()
	data := make([]byte, l)

	n, err := r.ReaderAt.ReadAt(data, r.off)
	if err != nil {
		if errors.Is(err, io.EOF) {
			err = nil
		}
	}

	return data[:n], err
}

// ReadLine reads a single line from the Reader.
func (r *MmapReader) ReadLine() (b []byte, err error) {
	rat := (*readerAt)(unsafe.Pointer(r.ReaderAt))

	i := r.off
	data := rat.data

	l := int64(len(data))
	if len(data) == 0 || i >= l {
		err = io.EOF
		return
	}

	// find the newline
	for i < l && data[i] != '\n' {
		i++
	}

	b = data[r.off:i]

	// move to next line
	if i < l {
		r.off = i + 1
	}

	return
}

// Reset resets the Reader to the beginning of the file.
func (r *MmapReader) Reset() {
	r.off = 0
}

// Close closes the Reader.
func (r *MmapReader) Close() error {
	return r.ReaderAt.Close()
}
