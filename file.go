package octrl

import (
	"os"
	"github.com/yulon/go-bin"
)

func FileAlign(f *os.File, size int64) {
	fi, err := f.Stat()
	if err != nil {
		return
	}
	m := fi.Size() % size
	if m > 0 {
		f.Write(bin.Zeros(size - m))
	}
}
