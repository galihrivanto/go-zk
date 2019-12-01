package push

import (
	"bytes"
	"sync"
)

var (
	bufferPoolSize = 10
	bufferPool     chan *bytes.Buffer
	bufferPoolLock sync.Once
)

// SetBufferPoolSize override default buffer pool size (10)
// it might caused high memory usage when value is too high
// and slow performance if value is too low
func SetBufferPoolSize(size int) {
	bufferPoolSize = size

	// reinit buffer
	initBufferPool()
}

func acquireBuffer() *bytes.Buffer {
	return <-bufferPool
}

func releaseBuffer(buf *bytes.Buffer) {
	// ensure buffer content cleared before put back to pool
	buf.Reset()

	bufferPool <- buf
}

func initBufferPool() {
	// init pool
	bufferPool = make(chan *bytes.Buffer, bufferPoolSize)

	// fill pool with empty buffer
	n := 0
	for n < bufferPoolSize {
		bufferPool <- &bytes.Buffer{}

		n++
	}
}

func init() {
	bufferPoolLock.Do(initBufferPool)
}
