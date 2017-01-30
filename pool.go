package chunked

import (
	"errors"
	"sync"

	"github.com/liferoot/bits"
)

type Pool struct {
	offset            int
	smallest, largest int
	chunk             []sync.Pool
}

func (p *Pool) Smallest() int { return p.smallest }
func (p *Pool) Largest() int  { return p.largest }

func (p *Pool) Nearest(size int) int {
	switch {
	case p.smallest > size:
		size = p.smallest
	case p.largest < size:
		size = p.largest
	case (size & (size - 1)) != 0:
		size = int(1 << (bits.MSB64(uint64(size)) + 1))
	}
	return size
}

func (p *Pool) Get(size int) []byte {
	size = p.Nearest(size)
	if chunk := p.chunk[int(bits.MSB64(uint64(size)))-p.offset].Get(); chunk != nil {
		return chunk.([]byte)
	}
	return make([]byte, 0, size)
}

func (p *Pool) Put(chunk []byte) error {
	if size := cap(chunk); p.smallest <= size && size <= p.largest && (size&(size-1)) == 0 {
		p.chunk[int(bits.MSB64(uint64(size)))-p.offset].Put(chunk[:0])
		return nil
	}
	return ErrWrongSize
}

func NewPool(smallest, largest int) *Pool {
	switch {
	case smallest > largest:
		panic(`chunked.Pool: the smallest chunk size must be less than or equal to the largest chunk size.`)
	case smallest <= 0 || (smallest&(smallest-1)) != 0:
		panic(`chunked.Pool: the smallest chunk size must be greater than zero and a power of 2.`)
	case (largest & (largest - 1)) != 0:
		panic(`chunked.Pool: the largest chunk size must be a power of 2.`)
	}
	pool := &Pool{
		offset:   int(bits.MSB64(uint64(smallest))),
		smallest: smallest,
		largest:  largest,
	}
	pool.chunk = make([]sync.Pool, int(bits.MSB64(uint64(largest)))-pool.offset+1)

	return pool
}

func GetChunk(size int) []byte    { return defaultPool.Get(size) }
func GetSmallestChunk() []byte    { return defaultPool.Get(defaultPool.Smallest()) }
func GetLargestChunk() []byte     { return defaultPool.Get(defaultPool.Largest()) }
func PutChunk(chunk []byte) error { return defaultPool.Put(chunk) }

func init() {
	defaultPool = Pool{smallest: 1 << 8, largest: 1 << 18}
	defaultPool.offset = int(bits.MSB64(uint64(defaultPool.smallest)))
	defaultPool.chunk = make([]sync.Pool, int(bits.MSB64(uint64(defaultPool.largest)))-defaultPool.offset+1)
}

var (
	defaultPool Pool

	ErrWrongSize = errors.New(`wrong chunk size`)
)
