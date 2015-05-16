package octrl

import (
	"io"
	"github.com/yulon/go-bin"
)

func Align(ws io.WriteSeeker, alignment int64) {
	off, err := ws.Seek(0, 1)
	if err != nil {
		return
	}
	m := off % alignment
	if m > 0 {
		ws.Write(bin.Zeros(alignment - m))
	}
}
