# CloserSet

CloserSet is a set of `io.Closer`. It modifies and records `io.Closer`.

When `CloserSet.Close()` is called, it closes all underlying closers which have not been closed.

