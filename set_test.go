package closerset

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type closerCount struct {
	counter int
}

func (c *closerCount) Close() error {
	c.counter++
	return nil
}

func TestCloserSet(t *testing.T) {
	set := &CloserSet{}
	cnt := &closerCount{}
	closer := set.WrapAndRecord(cnt)
	require.NoError(t, closer.Close())
	require.NoError(t, closer.Close())
	require.NoError(t, set.Close())
	require.NoError(t, closer.Close())
	require.Equal(t, 1, cnt.counter)
}
