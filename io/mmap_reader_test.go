package io

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	const filename = "mmap_reader.go"
	r, err := NewMmapReader(filename)
	if err != nil {
		t.Fatalf("NewReader: %v", err)
	}
	defer r.Close()

	got := make([]byte, r.Len())
	if _, err := r.ReadAt(got, 0); err != nil && err != io.EOF {
		t.Fatalf("ReadAt: %v", err)
	}
	want, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("os.ReadFile: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d bytes, want %d", len(got), len(want))
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("\ngot  %q\nwant %q", string(got), string(want))
	}

	r.Reset()
	got, err = r.ReadAll()
	if err != nil && err != io.EOF {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d bytes, want %d", len(got), len(want))
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("\ngot  %q\nwant %q", string(got), string(want))
	}

	got = got[:0]
	for {
		line, err := r.ReadLine()
		if err != nil && err != io.EOF {
			t.Fatalf("ReadAll: %v", err)
		}
		if err != nil && err == io.EOF {
			break
		}
		got = append(got, line...)
		got = append(got, []byte("\n")...)
	}
	if len(got) != len(want) {
		t.Fatalf("got %d bytes, want %d", len(got), len(want))
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("\ngot  %q\nwant %q", string(got), string(want))
	}
}
