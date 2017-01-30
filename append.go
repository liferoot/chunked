package chunked

import "github.com/liferoot/bits"

func Append(b []byte, p ...[]byte) (r []byte) {
	z := 0

	for _, r = range p {
		z += len(r)
	}
	if z == 0 {
		return b
	}
	n := len(b)
	z += n

	if z <= cap(b) {
		r = b[:z]
	} else {
		if z > defaultPool.largest {
			r = make([]byte, z, int(1<<(bits.MSB64(uint64(z))+1)))
		} else {
			r = GetChunk(z)[:z]
		}
		copy(r, b)
	}
	for _, b = range p {
		n += copy(r[n:], b)
	}
	return
}

func AppendByte(b []byte, p ...byte) (r []byte) {
	if len(b) == 0 {
		return p
	}
	if len(p) == 0 {
		return b
	}
	n := len(b)
	z := n + len(p)

	if z <= cap(b) {
		r = b[:z]
	} else {
		if z > defaultPool.largest {
			r = make([]byte, z, int(1<<(bits.MSB64(uint64(z))+1)))
		} else {
			r = GetChunk(z)[:z]
		}
		copy(r, b)
	}
	for i, z := 0, len(p); i < z; i++ {
		r[n] = p[i]
		n++
	}
	return
}
