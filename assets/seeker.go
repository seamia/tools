// Copyright 2017 Seamia Corporation. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package assets

import (
	"fmt"
	"io"
)

func CreateReader(from []byte) (io.ReadSeeker, error) {
	return &readerSeekerCloser{
		data:    from,
		total:   int64(len(from)),
		open:    true,
		current: 0,
	}, nil
}

type readerSeekerCloser struct {
	io.Reader
	io.Seeker
	io.Closer

	data    []byte
	total   int64
	current int64
	open    bool
}

// Seek implements the io.Seeker interface.
func (r *readerSeekerCloser) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case io.SeekStart:
		abs = offset
	case io.SeekCurrent:
		abs = r.current + offset
	case io.SeekEnd:
		abs = int64(r.total) + offset
	default:
		return 0, fmt.Errorf("strings.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, fmt.Errorf("strings.Reader.Seek: negative position")
	} else if abs > r.total {
		return 0, fmt.Errorf("strings.Reader.Seek: post end.of.data position")
	}
	r.current = abs
	return abs, nil
}

func (r *readerSeekerCloser) Read(p []byte) (n int, err error) {
	if !r.open {
		return 0, fmt.Errorf("the readed already closed")
	}
	if p == nil || len(p) == 0 {
		return 0, nil
	}
	looking4 := int64(len(p))
	if r.current+looking4 > r.total {
		looking4 = r.total - r.current
	}

	if looking4 == 0 {
		return 0, io.EOF
	}

	copy(p, r.data[r.current:])
	r.current += looking4
	return int(looking4), nil
}

func (r *readerSeekerCloser) Close() error {
	r.open = false
	return nil
}
