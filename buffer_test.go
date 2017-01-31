package chunked

import (
	"bytes"
	"io"
	"testing"
)

func TestBuffer_NewBufferWithoutPool(t *testing.T) {
	buf := NewBuffer(nil)
	if buf.pool != &defaultPool {
		t.Error(`unexpected pool`)
	}
}

func TestBuffer_Pool(t *testing.T) {
	buf := new(Buffer)
	if buf.Pool() != &defaultPool {
		t.Error(`unexpected pool`)
	}
}

func TestBuffer_Reset(t *testing.T) {
	buf := NewBuffer(NewPool(4, 4))
	buf.Write(bb)
	if buf.Len() != len(bb) {
		t.Errorf("expected length %d, got %d", len(bb), buf.Len())
	}
	buf.Reset()
	if buf.Len() > 0 {
		t.Error(`buffer must be empty`)
	}
}

func TestBuffer_Write(t *testing.T) {
	cases := [...]struct {
		in  []byte
		exp []byte
		len int
		cap int
	}{
		{bb[:0], nil, 0, 0},
		{bb[:4], bb[:4], 4, 8},
		{bb[4:7], bb[:7], 7, 8},
		{bb[:26], bb[17:26], 33, 40},
	}
	buf := NewBuffer(NewPool(8, 16))

	for i, c := range cases {
		n, err := buf.Write(c.in)
		out := last(buf)
		if err != nil || len(c.in) != n || c.len != buf.Len() || c.cap != buf.Cap() || !bytes.Equal(c.exp, out) {
			t.Errorf("case %d:\n\tfor: %q\n"+
				"\texp: %q, writ/len/cap: %d/%d/%d, err: <nil>\n"+
				"\tgot: %q, writ/len/cap: %d/%d/%d, err: %v\n",
				i, c.in,
				c.exp, len(c.in), c.len, c.cap,
				out, n, buf.Len(), buf.Cap(), err)
		}
	}
}

func TestBuffer_WriteByte(t *testing.T) {
	cases := [...]struct {
		in  byte
		exp []byte
		len int
		cap int
	}{
		{bb[0], bb[:1], 1, 4},
		{bb[1], bb[:2], 2, 4},
		{bb[2], bb[:3], 3, 4},
		{bb[3], bb[:4], 4, 4},
		{bb[4], bb[4:5], 5, 8},
	}
	buf := NewBuffer(NewPool(4, 4))

	for i, c := range cases {
		err := buf.WriteByte(c.in)
		out := last(buf)
		if err != nil || c.len != buf.Len() || c.cap != buf.Cap() || !bytes.Equal(c.exp, out) {
			t.Errorf("case %d:\n\tfor: %q\n"+
				"\texp: %q, len/cap: %d/%d, err: <nil>\n"+
				"\tgot: %q, len/cap: %d/%d, err: %v\n",
				i, c.in,
				c.exp, c.len, c.cap,
				out, buf.Len(), buf.Cap(), err)
		}
	}
}

func TestBuffer_WriteByteWithoutPool(t *testing.T) {
	buf := new(Buffer)
	buf.WriteByte('a')
	if buf.Len() != 1 {
		t.Errorf("expected length 1, got %d", buf.Len())
	}
}

func TestBuffer_Bytes(t *testing.T) {
	buf := NewBuffer(NewPool(8, 16))
	out := buf.Bytes()
	if out != nil {
		t.Errorf("expected <nil>, got %q", out)
	}
	buf.Write(bb)
	out = buf.Bytes()
	if !bytes.Equal(bb, out) {
		t.Errorf("\n\tfor: %q\n\texp: %q\n\tgot: %q\n", bb, bb, out)
	}
}

func TestBuffer_ReadFrom(t *testing.T) {
	buf := NewBuffer(NewPool(8, 16))
	n, err := buf.ReadFrom(bytes.NewBuffer(bb))
	exp, out := bb[32:], last(buf)
	if err != nil || len(bb) != int(n) || !bytes.Equal(exp, out) {
		t.Errorf("\n\tfor: %q\n\texp: %q, read: %d, err: <nil>\n\tgot: %q, read: %d, err: %v\n", bb, exp, len(bb), out, n, err)
	}
}

func TestBuffer_WriteTo(t *testing.T) {
	buf, w := NewBuffer(NewPool(8, 16)), new(bytes.Buffer)
	exp := bb
	buf.Write(exp)
	n, err := buf.WriteTo(w)
	out := w.Bytes()
	if err != nil || len(exp) != int(n) || !bytes.Equal(exp, out) {
		t.Errorf("\n\tfor: %q\n\texp: %q, read: %d, err: <nil>\n\tgot: %q, read: %d, err: %v\n", exp, exp, len(bb), out, n, err)
	}
	if buf.Len() > 0 {
		t.Error(`buffer should be empty`)
	}
}

func TestBuffer_Read(t *testing.T) {
	cases := [...]struct {
		in  int
		exp []byte
		off int
		len int
		cap int
		err error
	}{
		{4, bb[:4], 4, 31, 40, nil},
		{5, bb[4:9], 1, 26, 32, nil},
		{10, bb[9:19], 3, 16, 24, nil},
		{16, bb[19:], 0, 0, 0, nil},
		{4, []byte{0, 0, 0, 0}, 0, 0, 0, io.EOF},
	}
	buf := NewBuffer(NewPool(8, 8))
	buf.Write(bb)

	for i, c := range cases {
		out := make([]byte, c.in)
		n, err := buf.Read(out)
		m := 0
		if err != io.EOF {
			m = len(out)
		}
		if c.err != err || m != n || c.len != buf.Len() || c.cap != buf.Cap() || !bytes.Equal(c.exp, out) {
			t.Errorf("case %d:\n\tfor: [0:%d]\n"+
				"\texp: %q, read/len/cap: %d/%d/%d, err: %v\n"+
				"\tgot: %q, read/len/cap: %d/%d/%d, err: %v\n",
				i, c.in,
				c.exp, m, c.len, c.cap, c.err,
				out, n, buf.Len(), buf.Cap(), err)
		}
	}
}

func TestBuffer_ReadByte(t *testing.T) {
	m := 15
	buf := NewBuffer(NewPool(4, 4))
	buf.Write(bb[:m])

	for i := 0; i <= m; i++ {
		b, err := buf.ReadByte()
		if err != nil {
			if err != io.EOF {
				t.Errorf("case %d: err should always be io.EOF, got %v", i, err)
			}
			break
		}
		if b != bb[i] {
			t.Errorf("case %d: expected %q, got %q", i, b, bb[i])
		}
	}
	if buf.Len() > 0 || buf.Cap() > 0 {
		t.Errorf("buffer should be empty in the final, but has a len/cap == %d/%d", buf.Len(), buf.Cap())
	}
}

func TestBuffer_PeekByte(t *testing.T) {
	buf := NewBuffer(NewPool(4, 4))
	exp, out := byte(0), buf.PeekByte()
	if exp != out {
		t.Errorf("buffer is empty: expected %q, got %q", exp, out)
	}
	buf.Write(bb)

	for i, m := 0, len(bb); i < m; i++ {
		exp, out = bb[i], buf.PeekByte()
		buf.ReadByte()
		if exp != out {
			t.Errorf("case %d: expected %q, got %q", i, exp, out)
		}
	}
	exp, out = byte(0), buf.PeekByte()
	if exp != out {
		t.Errorf("buffer is completely read: expected %q, got %q", exp, out)
	}
}

func TestBuffer_LastByte(t *testing.T) {
	buf := NewBuffer(NewPool(4, 4))
	exp, out := byte(0), buf.LastByte()
	if exp != out {
		t.Errorf("buffer is empty: expected %q, got %q", exp, out)
	}
	for i, m := 0, len(bb); i < m; i++ {
		exp = bb[i]
		buf.WriteByte(exp)
		out = buf.LastByte()
		if exp != out {
			t.Errorf("case %d: expected %q, got %q", i, exp, out)
		}
	}
}

var (
	bb = []byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		194, 181, // µ
		225, 180, 131, // ᴃ
		240, 166, 164, 128, // 𦤀
	}
)

func first(b *Buffer) []byte {
	if len(b.chunk) > 0 {
		return b.chunk[0]
	}
	return nil
}

func last(b *Buffer) []byte {
	if i := len(b.chunk) - 1; i >= 0 {
		return b.chunk[i]
	}
	return nil
}
