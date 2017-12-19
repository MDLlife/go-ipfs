// package chunk implements streaming block splitters
package chunk

import (
	"io"

	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
)

var log = logging.Logger("chunk")

var DefaultBlockSize int64 = 1024 * 256

type Splitter interface {
	Reader() io.Reader
	NextBytes() ([]byte, error)
}

type SplitterGen func(r io.Reader) Splitter

func DefaultSplitter(r io.Reader) Splitter {
	return NewSizeSplitter(r, DefaultBlockSize)
}

func SizeSplitterGen(size int64) SplitterGen {
	return func(r io.Reader) Splitter {
		return NewSizeSplitter(r, size)
	}
}

func Chan(s Splitter) (<-chan []byte, <-chan error) {
	out := make(chan []byte)
	errs := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errs)

		// all-chunks loop (keep creating chunks)
		for {
			b, err := s.NextBytes()
			if err != nil {
				errs <- err
				return
			}

			out <- b
		}
	}()
	return out, errs
}

type sizeSplitterv2 struct {
	r   io.Reader
	buf []byte
	err error
}

func NewSizeSplitter(r io.Reader, size int64) Splitter {
	return &sizeSplitterv2{
		r:   r,
		buf: make([]byte, size),
	}
}

func (ss *sizeSplitterv2) NextBytes() ([]byte, error) {
	if ss.err != nil {
		return nil, ss.err
	}

	size := len(ss.buf)
	n, err := io.ReadFull(ss.r, ss.buf)
	switch err {
	case io.ErrUnexpectedEOF:
		ss.err = io.EOF
		ret := make([]byte, n)
		copy(ret, ss.buf)
		ss.buf = nil
		return ret, nil
	case nil:
		var ret []byte
		ss.buf, ret = make([]byte, size), ss.buf
		return ret, nil
	default:
		return nil, err
	}
}

func (ss *sizeSplitterv2) Reader() io.Reader {
	return ss.r
}
