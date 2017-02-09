package chunked

import (
	"io"

	"github.com/liferoot/bits"
)

type Buffer struct {
	pool     *Pool
	chunk    [][]byte
	offset   int
	length   int
	capacity int
}

// Len returns the number of bytes of the unread portion of the buffer.
func (b *Buffer) Len() int { return b.length }

// Cap returns the capacity of the buffer.
func (b *Buffer) Cap() int { return b.capacity }

// Pool
func (b *Buffer) Pool() *Pool {
	if b.pool == nil {
		b.pool = &defaultPool
	}
	return b.pool
}

// Bytes returns a slice of the contents of the unread portion of the buffer.
func (b *Buffer) Bytes() (p []byte) {
	if len(b.chunk) == 0 {
		return nil
	}
	if z := b.pool.Nearest(b.length); z < b.length {
		p = make([]byte, int(1<<(bits.MSB64(uint64(b.length))+1)))
	} else {
		p = (b.pool.Get(z))[:z]
	}
	n := copy(p[0:], b.chunk[0][b.offset:])

	for i, l := 1, len(b.chunk); i < l; i++ {
		n += copy(p[n:], b.chunk[i])
	}
	return p[:n]
}

// Reset resets the buffer so it has no content.
func (b *Buffer) Reset() {
	for i, l := 0, len(b.chunk); i < l; i++ {
		b.pool.Put(b.chunk[i])
		b.chunk[i] = nil
	}
	b.chunk = b.chunk[:0]
	b.offset = 0
	b.length = 0
	b.capacity = 0
}

// PeekByte returns the next byte without advancing the buffer.
// This byte is valid until the next read call.
// If the buffer is empty, it returns byte(0).
func (b *Buffer) PeekByte() byte {
	if len(b.chunk) == 0 {
		return 0
	}
	return b.chunk[0][b.offset]
}

// LastByte returns the last byte without advancing the buffer.
// This byte is valid until the next write call.
// If the buffer is empty, it returns byte(0).
func (b *Buffer) LastByte() byte {
	i := len(b.chunk) - 1
	if i < 0 {
		return 0
	}
	return b.chunk[i][len(b.chunk[i])-1]
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	for a, m := 0, len(p); m > n; {
		if len(b.chunk) == 0 {
			err = io.EOF
			break
		}
		c := b.chunk[0]
		a = copy(p[n:], c[b.offset:])
		n += a

		if b.offset += a; b.offset == len(c) {
			b.pop()
		}
	}
	b.length -= n
	return
}

func (b *Buffer) ReadByte() (byte, error) {
	if len(b.chunk) == 0 {
		return 0, io.EOF
	}
	c := b.chunk[0][b.offset]
	if b.offset++; b.offset == len(b.chunk[0]) {
		b.pop()
	}
	b.length--
	return c, nil
}

func (b *Buffer) ReadFrom(r io.Reader) (_ int64, err error) {
	var i, n int

	if i = len(b.chunk) - 1; i < 0 {
		if b.pool == nil {
			b.pool = &defaultPool
		}
		b.push(b.pool.Largest())
		i++
	}
	for a, c := 0, &b.chunk[i]; err == nil; {
		if len(*c) == cap(*c) {
			b.push(b.pool.Largest())
			i++
			c = &b.chunk[i]
		}
		a, err = r.Read((*c)[len(*c):cap(*c)])
		*c = (*c)[:len(*c)+a]
		n += a
	}
	if err == io.EOF {
		err = nil
	}
	b.length += n
	return int64(n), err
}

// Write appends the contents of the slice to the buffer.
func (b *Buffer) Write(p []byte) (n int, err error) {
	m := len(p)
	if m == 0 {
		return
	}
	i := len(b.chunk) - 1
	if i < 0 {
		if b.pool == nil {
			b.pool = &defaultPool
		}
		b.push(m)
		i++
	}
	for a, c := 0, &b.chunk[i]; m > n; {
		if len(*c) == cap(*c) {
			b.push(m - n) // mb simply m?
			i++
			c = &b.chunk[i]
		}
		a = copy((*c)[len(*c):cap(*c)], p[n:])
		*c = (*c)[:len(*c)+a]
		n += a
	}
	b.length += n
	return
}

// WriteByte appends the byte to the buffer.
func (b *Buffer) WriteByte(p byte) error {
	i := len(b.chunk) - 1
	if i < 0 || len(b.chunk[i]) == cap(b.chunk[i]) {
		if b.pool == nil {
			b.pool = &defaultPool
		}
		b.push(1)
		i++
	}
	c := &b.chunk[i]
	j := len(*c)

	*c = (*c)[:j+1]
	(*c)[j] = p

	b.length++
	return nil
}

func (b *Buffer) WriteTo(w io.Writer) (_ int64, err error) {
	var d, i, n int

	for a, m := 0, len(b.chunk); i < m && err == nil; i++ {
		d = len(b.chunk[i]) - b.offset
		a, err = w.Write(b.chunk[i][b.offset:])
		n += a

		if b.offset += a; b.offset == len(b.chunk[i]) {
			b.pool.Put(b.chunk[i])
			b.capacity -= cap(b.chunk[i])
			b.chunk[i] = nil
			b.offset = 0
		}
		if a != d {
			err = io.ErrShortWrite
		}
	}
	b.chunk = b.chunk[i:]
	b.length -= n
	return int64(n), err
}

func (b *Buffer) pop() {
	b.pool.Put(b.chunk[0])
	b.capacity -= cap(b.chunk[0])
	b.chunk[0], b.chunk = nil, b.chunk[1:]
	b.offset = 0
}

func (b *Buffer) push(n int) {
	c := b.pool.Get(n)
	b.chunk = append(b.chunk, c)
	b.capacity += cap(c)
}

func NewBuffer(pool *Pool) (b *Buffer) {
	if pool == nil {
		pool = &defaultPool
	}
	return &Buffer{pool: pool, chunk: make([][]byte, 0, 64)}
}
