package closerset

import (
	"io"
	"sync"

	"emperror.dev/errors"
)

// CloserSet records a set of closers.
// All closers will be closed when the CloserSet closes.
// When a wrapped closer is closed, the record in the CloserSet will be removed. i.e.,
// the closer will only be closed once.
type CloserSet struct {
	closers    map[uint64]io.Closer
	nextID     uint64
	closersMtx sync.Mutex
}

type wrappedCloser struct {
	id  uint64
	set *CloserSet
}

// Close the wrapped closer once.
func (w *wrappedCloser) Close() error {
	w.set.closersMtx.Lock()
	if w.set.closers == nil {
		w.set.closersMtx.Unlock()
		return nil
	}
	c, ok := w.set.closers[w.id]
	if !ok {
		w.set.closersMtx.Unlock()
		return nil
	}
	delete(w.set.closers, w.id)
	w.set.closersMtx.Unlock()
	return c.Close()
}

func (c *CloserSet) ensureSet() {
	if c.closers == nil {
		c.closers = map[uint64]io.Closer{}
	}
}

// WrapAndRecord wrap a closer and return wrapped closer.
// NOTE: the input closer should not be invoked after wrap.
func (c *CloserSet) WrapAndRecord(closer io.Closer) io.Closer {
	c.closersMtx.Lock()
	c.ensureSet()
	id := c.nextID
	c.nextID++
	c.closers[id] = closer
	c.closersMtx.Unlock()
	return &wrappedCloser{
		id:  id,
		set: c,
	}
}

// Close underlying closers.
func (c *CloserSet) Close() error {
	var err error
	c.closersMtx.Lock()
	if c.closers != nil {
		for _, closer := range c.closers {
			err = errors.Append(err, closer.Close())
		}
		c.closers = nil
	}
	c.closersMtx.Unlock()
	return err
}
