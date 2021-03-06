package octrl

import (
	"io"
	"errors"
	"github.com/yulon/go-bin"
)

type Labeler struct{
	ws io.WriteSeeker
	labs map[string]int64
	pits []pit
	base int64
}

type pit struct{
	addr int64
	start string
	end string
	added int64
	wc bin.WordConv
}

func NewLabeler(ws io.WriteSeeker) *Labeler {
	offset, err := ws.Seek(0, 1)
	if err != nil {
		return nil
	}

	laber := &Labeler{
		ws: ws,
		labs: map[string]int64{},
		pits: []pit{},
		base: offset,
	}
	return laber
}

func (laber *Labeler) Label(label string) error {
	offset, err := laber.ws.Seek(0, 1)
	if err != nil {
		return err
	}

	laber.labs[label] = offset
	return nil
}

func (laber *Labeler) Get(label string) (int64, error) {
	off, ok := laber.labs[label]
	if !ok {
		return 0, errors.New(label + " undefined")
	}
	return off, nil
}

func (laber *Labeler) Pit(startLabel string, endLabel string, added int64, wc bin.WordConv) (int, error) {
	addr, err := laber.ws.Seek(0, 1)
	if err != nil {
		return 0, err
	}

	laber.pits = append(laber.pits, pit{
		addr: addr,
		start: startLabel,
		end: endLabel,
		added: added,
		wc: wc,
	})
	return laber.ws.Write(wc(0))
}

func (laber *Labeler) Close() error {
	current, err := laber.ws.Seek(0, 1)
	if err != nil {
		return err
	}

	for i := 0; i < len(laber.pits); i++ {
		var start, end int64
		var err error

		if laber.pits[i].start == "" {
			start = laber.base
		}else{
			start, err = laber.Get(laber.pits[i].start)
			if err != nil {
				return err
			}
		}

		if laber.pits[i].end == "" {
			end = laber.pits[i].addr
		}else{
			end, err = laber.Get(laber.pits[i].end)
			if err != nil {
				return err
			}
		}

		n := end - start + laber.pits[i].added

		_, err = laber.ws.Seek(laber.pits[i].addr, 0)
		if err != nil {
			return err
		}

		laber.ws.Write(laber.pits[i].wc(n))
	}

	_, err = laber.ws.Seek(current, 0)
	return err
}
